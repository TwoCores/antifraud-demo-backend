package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var (
	signingKey []byte
)

type ctxKey string

const CtxKeyClaims ctxKey = "jwt_claims"

func init() {
	viper.AutomaticEnv()
	signingKey = []byte(viper.GetString("SECRET_KEY"))
}

func GenerateToken(uid string, phoneModel, osStr string) (string, error) {
	now := time.Now()
	claims := &JwtClaims{
		UserId:      uid,
		IsSuperuser: false,
		ExpiresAt:   now.Add(24 * time.Hour).Unix(),
		IssuedAt:    now.Unix(),
		NotBefore:   now.Unix(),
		Subject:     uid,
		PhoneModel:  phoneModel,
		OS:          osStr,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func GenerateSUToken(uid string) (string, error) {
	now := time.Now()
	claims := &JwtClaims{
		UserId:      uid,
		IsSuperuser: true,
		ExpiresAt:   now.Add(12 * time.Hour).Unix(),
		IssuedAt:    now.Unix(),
		NotBefore:   now.Unix(),
		Subject:     uid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func ValidateToken(ctx context.Context, tokenString string) (context.Context, bool, error) {
	tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer"))
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return ctx, false, nil
	}

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return signingKey, nil
	}

	tkn, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, keyFunc)
	if err != nil {
		return ctx, false, err
	}
	if !tkn.Valid {
		return ctx, false, nil
	}

	claims, ok := tkn.Claims.(*JwtClaims)
	if !ok {
		return ctx, false, nil
	}

	ctx = context.WithValue(ctx, CtxKeyClaims, claims)
	return ctx, true, nil
}
