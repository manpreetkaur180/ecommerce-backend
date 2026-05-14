package user

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"user-service/config"
	"user-service/pkg/utils"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"time"
)

type Service struct {
	Repo *Repository
	RDB  *redis.Client
}

var ErrEmailAlreadyVerified = errors.New("email already verified")

func NewService(repo *Repository, rdb *redis.Client) *Service {
	return &Service{
		Repo: repo,
		RDB:  rdb,
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

func publicUserServiceURL() string {
	baseURL := os.Getenv("PUBLIC_USER_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3001"
	}
	return strings.TrimRight(baseURL, "/")
}

func (s *Service) SendOTP(
	user *User,
	identifier string,
	channel string,
) error {

	otp := generateOTP()

	err := s.RDB.Set(
		config.Ctx,
		"otp:"+identifier,
		otp,
		5*time.Minute,
	).Err()

	if err != nil {
		return err
	}

	if channel == "email" {
		if err := utils.PublishOTPEmail(
			user.Name,
			user.Email,
			otp,
		); err != nil {
			return err
		}

		log.Println("Published OTP email event for", user.Email)
		return nil
	}

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

	s.RDB.Del(config.Ctx, key)

	return nil
}

func (s *Service) FindByIdentifier(email, phone string) (*User, error) {
	u, err := s.Repo.FirstUserForLogin(email, phone)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return u, nil
}

func (s *Service) FindByID(userID uint) (*User, error) {
	u, err := s.Repo.FirstUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (s *Service) Register(req RegisterRequest) (*User, error) {

	req.Email = utils.NormalizeEmail(req.Email)
	req.Phone = utils.NormalizePhone(req.Phone)

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

	if err := utils.ValidateEmail(req.Email); err != nil {
		return nil, err
	}

	if err := utils.ValidatePhone(req.Phone); err != nil {
		return nil, err
	}

	if err := utils.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	if _, err := s.Repo.FirstUserByEmail(req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	if _, err := s.Repo.FirstUserByPhone(req.Phone); err == nil {
		return nil, errors.New("phone already exists")
	}

	if req.Password != req.ConfirmPassword {
		return nil, errors.New("passwords do not match")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(hashed),

		Role:       RoleBuyer,
		IsVerified: false,
	}

	if err := s.Repo.CreateUser(&user); err != nil {
		return nil, errors.New("failed to create user")
	}

	go func(user User) {
		if err := utils.PublishUserRegistered(
			user.Name,
			user.Email,
			user.Phone,
		); err != nil {
			log.Println("Failed to publish user registered event:", err)
			return
		}

		log.Println("Published user registered event for", user.Email)

		time.Sleep(25 * time.Second)

		rawToken := generateVerificationToken()

		hashedToken := hashToken(rawToken)

		verification := EmailVerification{
			UserID:    user.ID,
			TokenHash: hashedToken,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		if err := s.Repo.CreateEmailVerification(&verification); err != nil {
			log.Println("Failed to create email verification token:", err)
			return
		}

		verifyLink := fmt.Sprintf(
			"%s/api/v1/user/verify-email?token=%s",
			publicUserServiceURL(),
			rawToken,
		)
		if err := utils.PublishVerificationEmail(
			user.Name,
			user.Email,
			verifyLink,
		); err != nil {
			log.Println("Failed to send verification email:", err)
		}

		log.Println("Published verification email event for", user.Email)
	}(user)
	return &user, nil
}

func (s *Service) Login(req LoginRequest) (*User, error) {

	req.Email = utils.NormalizeEmail(req.Email)
	req.Phone = utils.NormalizePhone(req.Phone)

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

	user, err := s.Repo.FirstUserForLogin(req.Email, req.Phone)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if !user.IsVerified {
		return nil, errors.New("please verify your email before login")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *Service) VerifyEmail(
	rawToken string,
) error {

	hashedToken := hashToken(rawToken)

	verification, err := s.Repo.FirstEmailVerificationByTokenHash(hashedToken)
	if err != nil {
		return errors.New("invalid verification token")
	}

	if verification.Used {
		return ErrEmailAlreadyVerified
	}

	if time.Now().After(verification.ExpiresAt) {
		return errors.New("verification token expired")
	}

	user, err := s.Repo.FirstUserByID(verification.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	if user.IsVerified {
		verification.Used = true
		if err := s.Repo.SaveEmailVerification(verification); err != nil {
			return errors.New("failed to verify user")
		}

		return ErrEmailAlreadyVerified
	}

	user.IsVerified = true

	if err := s.Repo.SaveUser(user); err != nil {
		return errors.New("failed to verify user")
	}

	verification.Used = true
	if err := s.Repo.SaveEmailVerification(verification); err != nil {
		return errors.New("failed to verify user")
	}

	return nil
}

func (s *Service) ForgotPassword(email string) error {
	email = utils.NormalizeEmail(email)
	if err := utils.ValidateEmail(email); err != nil {
		return err
	}

	user, err := s.Repo.FirstUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	resetToken, err := utils.GeneratePasswordResetJWT(user.ID)
	if err != nil {
		return errors.New("failed to create reset token")
	}

	resetLink := fmt.Sprintf(
		"%s/api/v1/user/reset-password?token=%s",
		publicUserServiceURL(),
		resetToken,
	)

	if err := utils.SendResetPasswordEmail(
		user.Name,
		user.Email,
		resetLink,
	); err != nil {
		return errors.New("failed to send reset email")
	}

	return nil
}

func (s *Service) ResetPassword(rawToken, newPassword, confirmPassword string) error {
	if rawToken == "" {
		return errors.New("token required")
	}
	if newPassword == "" {
		return errors.New("new password is required")
	}
	if confirmPassword == "" {
		return errors.New("confirm password is required")
	}
	if newPassword != confirmPassword {
		return errors.New("passwords do not match")
	}
	if err := utils.ValidatePassword(newPassword); err != nil {
		return err
	}

	claims, err := utils.ParsePasswordResetJWT(rawToken)
	if err != nil {
		return errors.New("invalid reset token")
	}

	user, err := s.Repo.FirstUserByID(claims.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	if claims.IssuedAt == nil || claims.IssuedAt.Time.Before(user.UpdatedAt.Truncate(time.Second)) {
		return errors.New("invalid reset token")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.Password = string(hashed)

	if err := s.Repo.SaveUser(user); err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

func (s *Service) RequestUpdatePassword(email string) error {
	email = utils.NormalizeEmail(email)
	if err := utils.ValidateEmail(email); err != nil {
		return err
	}

	user, err := s.Repo.FirstUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	if !user.IsVerified {
		return errors.New("please verify your email before updating password")
	}

	rawToken := generateVerificationToken()
	hashedToken := hashToken(rawToken)

	update := PasswordUpdate{
		UserID:    user.ID,
		TokenHash: hashedToken,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := s.Repo.CreatePasswordUpdate(&update); err != nil {
		return errors.New("failed to create update token")
	}

	updateLink := fmt.Sprintf(
		"%s/api/v1/user/update-password?token=%s",
		publicUserServiceURL(),
		rawToken,
	)

	if err := utils.SendUpdatePasswordEmail(
		user.Name,
		user.Email,
		updateLink,
	); err != nil {
		return errors.New("failed to send update password email")
	}

	return nil
}

func (s *Service) ConfirmUpdatePassword(rawToken, oldPassword, newPassword, confirmNewPassword string) error {
	if rawToken == "" {
		return errors.New("token required")
	}
	if oldPassword == "" {
		return errors.New("old password is required")
	}
	if newPassword == "" {
		return errors.New("new password is required")
	}
	if confirmNewPassword == "" {
		return errors.New("confirm new password is required")
	}
	if newPassword != confirmNewPassword {
		return errors.New("passwords do not match")
	}
	if err := utils.ValidatePassword(newPassword); err != nil {
		return err
	}

	hashedToken := hashToken(rawToken)

	update, err := s.Repo.FirstPasswordUpdateByTokenHash(hashedToken)
	if err != nil {
		return errors.New("invalid update token")
	}
	if time.Now().After(update.ExpiresAt) {
		return errors.New("update token expired")
	}

	user, err := s.Repo.FirstUserByID(update.UserID)
	if err != nil {
		return errors.New("user not found")
	}
	if !user.IsVerified {
		return errors.New("please verify your email before updating password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("old password is incorrect")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.Password = string(hashed)
	if err := s.Repo.SaveUser(user); err != nil {
		return errors.New("failed to update password")
	}

	s.Repo.DeletePasswordUpdate(update)
	return nil
}
