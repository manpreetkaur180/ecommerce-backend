package user

import (
	"errors"

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
func (s *Service) SendOTP(identifier string) (string, error) {
	otp := generateOTP()

	err := s.RDB.Set(
		config.Ctx,
		"otp:"+identifier,
		otp,
		5*time.Minute,
	).Err()

	if err != nil {
		return "", err
	}

	return otp, nil
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

	// 2. validate input
	if req.Name == "" {
		return nil, errors.New("name is required")
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

	if req.Password != req.ConfirmPassword {
		return nil, errors.New("passwords do not match")
	}

	// 3. check existing email
	var existing User
	err := s.DB.Where("email = ?", req.Email).First(&existing).Error
	if err == nil {
		return nil, errors.New("email already exists")
	}

	// 4. hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// 5. create user
	user := User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(hashed),
	}

	err = s.DB.Create(&user).Error
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	return &user, nil
}

func (s *Service) Login(req LoginRequest) (*User, error) {

	// 1. require email or phone
	if req.Email == "" && req.Phone == "" {
		return nil, errors.New("email or phone is required")
	}

	if req.Password == "" {
		return nil, errors.New("password is required")
	}

	// 2. find user
	var user User
	var err error

	if req.Email != "" {
		err = s.DB.Where("email = ?", req.Email).First(&user).Error
	} else {
		err = s.DB.Where("phone = ?", req.Phone).First(&user).Error
	}

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 3. compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}