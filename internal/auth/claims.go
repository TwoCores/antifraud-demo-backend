package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	UserId string `json:"user_id,omitempty"`

	ExpiresAt int64    `json:"exp,omitempty"`
	IssuedAt  int64    `json:"iat,omitempty"`
	NotBefore int64    `json:"nbf,omitempty"`
	Issuer    string   `json:"iss,omitempty"`
	Subject   string   `json:"sub,omitempty"`
	Audience  []string `json:"aud,omitempty"`

	PhoneModel string `json:"phone_model,omitempty"`
	OS         string `json:"os,omitempty"`
}

func (c *JwtClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(c.ExpiresAt, 0)), nil
}

func (c *JwtClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(c.IssuedAt, 0)), nil
}

func (c *JwtClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(c.NotBefore, 0)), nil
}

func (c *JwtClaims) GetIssuer() (string, error) {
	return c.Issuer, nil
}

func (c *JwtClaims) GetSubject() (string, error) {
	return c.Subject, nil
}

func (c *JwtClaims) GetAudience() (jwt.ClaimStrings, error) {
	return c.Audience, nil
}
