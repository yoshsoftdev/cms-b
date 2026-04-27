package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	adminUser string
	passHash  []byte
	secret    string
}

func NewAuthService(adminUser, adminPassword, secret string) (*AuthService, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &AuthService{adminUser: adminUser, passHash: hash, secret: secret}, nil
}

func (s *AuthService) Login(username, password string) (string, error) {
	if username != s.adminUser {
		return "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword(s.passHash, []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	claims := jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *AuthService) Validate(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
}
