package authservice

import (
	"fmt"
	"gameapp/entity"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Service struct {
	signKey               string
	accessExpirationTime  time.Duration
	refreshExpirationTime time.Duration
	accessSubject         string
	refreshSubject        string
}

func New(sk string, as string, rs string, ae time.Duration, re time.Duration) Service {
	return Service{
		signKey:               sk,
		accessExpirationTime:  ae,
		refreshExpirationTime: re,
		accessSubject:         as,
		refreshSubject:        rs,
	}
}

func (s Service) GenerateAccessToken(u entity.User) (string, error) {
	token, err := s.generateNewJwtToken(u.ID, s.accessSubject, s.accessExpirationTime)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s Service) GenerateRefreshToken(u entity.User) (string, error) {
	token, err := s.generateNewJwtToken(u.ID, s.refreshSubject, s.refreshExpirationTime)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s Service) VerifyToken(bearerToken string) (entity.User, error) {
	return entity.User{}, nil
}

func (s Service) generateNewJwtToken(userID uint, subject string, expireDuration time.Duration) (string, error) {
	claim := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration))},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signingString, err := accessToken.SignedString([]byte(s.signKey))
	if err != nil {
		return "", fmt.Errorf("err in singkey : %w", err)
	}
	return signingString, nil
}
