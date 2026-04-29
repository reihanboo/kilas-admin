package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/reihanboo/kilas-admin/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) Login(email, password string) (*model.User, string, error) {
	var user model.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("invalid credentials")
		}
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := s.generateAccessToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func (s *AuthService) Me(userID uint) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) generateAccessToken(userID uint) (string, error) {
	secret := os.Getenv("ACCESS_TOKEN_SECRET")
	if secret == "" {
		return "", errors.New("missing ACCESS_TOKEN_SECRET")
	}

	expiresIn := 24 * time.Hour
	claims := jwt.MapClaims{
		"sub":  userID,
		"type": "access",
		"exp":  time.Now().Add(expiresIn).Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
