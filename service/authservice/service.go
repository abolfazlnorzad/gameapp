package authservice

import (
	"fmt"
	"gameapp/entity"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

type Service struct {
	config Config
}

type Config struct {
	SignKey               string
	AccessExpirationTime  time.Duration
	RefreshExpirationTime time.Duration
	AccessSubject         string
	RefreshSubject        string
}

func New(c Config) Service {
	return Service{
		config: c,
	}
}

func (s Service) GenerateAccessToken(u entity.User) (string, error) {
	token, err := s.generateNewJwtToken(u.ID, s.config.AccessSubject, s.config.AccessExpirationTime)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s Service) GenerateRefreshToken(u entity.User) (string, error) {
	token, err := s.generateNewJwtToken(u.ID, s.config.RefreshSubject, s.config.RefreshExpirationTime)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s Service) VerifyToken(bearerToken string) (*Claims, error) {
	tokenStr := strings.Replace(bearerToken, "Bearer ", "", 1)
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.SignKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func (s Service) generateNewJwtToken(userID uint, subject string, expireDuration time.Duration) (string, error) {
	claim := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration))},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signingString, err := accessToken.SignedString([]byte(s.config.SignKey))
	if err != nil {
		return "", fmt.Errorf("err in singkey : %w", err)
	}
	return signingString, nil
}
