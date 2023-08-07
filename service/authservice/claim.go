package authservice

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	jwt.RegisteredClaims
	UserID uint `json:"user_id"`
}

func (c Claims) GetExpirationTime() (*jwt.NumericDate, error) {
	return c.ExpiresAt, nil
}

func (c Claims) GetIssuedAt() (*jwt.NumericDate, error) {
	return c.IssuedAt, nil
}

func (c Claims) GetNotBefore() (*jwt.NumericDate, error) {
	return c.NotBefore, nil
}

func (c Claims) GetIssuer() (string, error) {
	return c.Issuer, nil
}

func (c Claims) GetSubject() (string, error) {
	return c.Subject, nil
}

func (c Claims) GetAudience() (jwt.ClaimStrings, error) {
	return c.Audience, nil
}
