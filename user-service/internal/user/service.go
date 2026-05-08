package user

import (
	"errors"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"user-service/config"
	"user-service/pkg/utils"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"time"
)

type Service struct {
	DB *gorm.DB

	RDB *redis.Client
}

func NewService(db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{
		DB:  db,
		RDB: rdb,
	}
}
func generateOTP() string {
	return fmt.Sprintf("%04d", rand.Intn(10000))
}

func generateVerificationToken() string {
	return fmt.Sprintf(
		"%d%d",
		time.Now().UnixNano(),
		rand.Intn(100000),
	)
}

func hashToken(token string) string {

	hash := sha256.Sum256([]byte(token))

	return hex.EncodeToString(hash[:])
}
func (s *Service) SendOTP(
	user *User,
	identifier string,
	channel string,
) error {

	otp := generateOTP()

	fmt.Println("SAVING OTP FOR:", identifier)
	fmt.Println("OTP:", otp)

	err := s.RDB.Set(
		config.Ctx,
		"otp:"+identifier,
		otp,
		5*time.Minute,
	).Err()

	if err != nil {
		return err
	}

	// -------- EMAIL OTP --------
	if channel == "email" {

		return utils.SendOTPEmail(
			user.Name,
			user.Email,
			otp,
		)
	}

	// -------- SMS OTP --------
	if channel == "phone" {

		return utils.SendOTPSMS(
			user.Phone,
			otp,
		)
	}

	return fmt.Errorf("invalid otp channel")
}
func (s *Service) VerifyOTP(identifier, otp string) error {
	key := "otp:" + identifier

	storedOTP, err := s.RDB.Get(config.Ctx, key).Result()
	if err != nil {
		return fmt.Errorf("otp expired or not found")
	}

	if storedOTP != otp {
		return fmt.Errorf("invalid otp")
	}
	fmt.Println("VERIFYING OTP FOR:", identifier)
fmt.Println("INPUT OTP:", otp)
	// delete after use
	s.RDB.Del(config.Ctx, key)

	return nil
}
func (s *Service) FindByIdentifier(email, phone string) (*User, error) {
	var user User
	var err error

	if email != "" {
		err = s.DB.Where("email = ?", email).First(&user).Error
	} else {
		err = s.DB.Where("phone = ?", phone).First(&user).Error
	}

	if err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (s *Service) Register(req RegisterRequest) (*User, error) {

	// 1. normalize input
	req.Email = utils.NormalizeEmail(req.Email)
	req.Phone = utils.NormalizePhone(req.Phone)

	// 2. REQUIRED FIELD CHECK (IMPORTANT ORDER FIX)

	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	if req.Email == "" {
		return nil, errors.New("email is required")
	}

	if req.Phone == "" {
		return nil, errors.New("phone is required")
	}

	if req.Password == "" {
		return nil, errors.New("password is required")
	}

	if req.ConfirmPassword == "" {
		return nil, errors.New("confirm password is required")
	}

	// 3. FORMAT VALIDATION
	if err := utils.ValidateEmail(req.Email); err != nil {
		return nil, err
	}

	if err := utils.ValidatePhone(req.Phone); err != nil {
		return nil, err
	}

	if err := utils.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	// 4. CHECK DUPLICATES
	var existing User
	if err := s.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		return nil, errors.New("email already exists")
	}

	if err := s.DB.Where("phone = ?", req.Phone).First(&existing).Error; err == nil {
		return nil, errors.New("phone already exists")
	}
	// 5. PASSWORD MATCH
	if req.Password != req.ConfirmPassword {
		return nil, errors.New("passwords do not match")
	}

	// 6. HASH PASSWORD
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// 7. CREATE USER
	user := User{
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
		Password:   string(hashed),
		IsVerified: false,
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	go utils.NotifyUserRegistered(
	user.Name,
	user.Email,
	user.Phone,
)

go func() {

	time.Sleep(1 * time.Minute)

	rawToken := generateVerificationToken()

	hashedToken := hashToken(rawToken)

	verification := EmailVerification{
		UserID:    user.ID,
		TokenHash: hashedToken,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	s.DB.Create(&verification)
verifyLink := fmt.Sprintf(
	"https://a15b-115-245-207-91.ngrok-free.app/api/v1/user/verify-email?token=%s",
	rawToken,
)
	utils.SendVerificationEmail(
		user.Name,
		user.Email,
		verifyLink,
	)

}()
	return &user, nil
}

func (s *Service) Login(req LoginRequest) (*User, error) {

	// 1. normalize input
	req.Email = utils.NormalizeEmail(req.Email)
	req.Phone = utils.NormalizePhone(req.Phone)

	// 2. validation
	if req.Email == "" && req.Phone == "" {
		return nil, errors.New("email or phone is required")
	}

	if req.Password == "" {
		return nil, errors.New("password is required")
	}

	if req.Email != "" {
		if err := utils.ValidateEmail(req.Email); err != nil {
			return nil, err
		}
	}

	if req.Phone != "" {
		if err := utils.ValidatePhone(req.Phone); err != nil {
			return nil, err
		}
	}

	// 3. find user (clean query)
	var user User
	query := s.DB

	if req.Email != "" {
		query = query.Where("email = ?", req.Email)
	} else {
		query = query.Where("phone = ?", req.Phone)
	}

	if err := query.First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}
	if !user.IsVerified {
	return nil, errors.New("please verify your email before login")
}

	// 4. verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}
func (s *Service) VerifyEmail(
	rawToken string,
) error {

	hashedToken := hashToken(rawToken)

	var verification EmailVerification

	err := s.DB.Where(
		"token_hash = ?",
		hashedToken,
	).First(&verification).Error

	if err != nil {
		return errors.New("invalid verification token")
	}

	if time.Now().After(verification.ExpiresAt) {
		return errors.New("verification token expired")
	}

	var user User

	if err := s.DB.First(
		&user,
		verification.UserID,
	).Error; err != nil {

		return errors.New("user not found")
	}

	user.IsVerified = true

	if err := s.DB.Save(&user).Error; err != nil {
		return errors.New("failed to verify user")
	}

	s.DB.Delete(&verification)

	return nil
}