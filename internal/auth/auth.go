package auth

import (
	"context"
	"fmt"
	"net/http"
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
		UserId:     uid,
		ExpiresAt:  now.Add(24 * time.Hour).Unix(),
		IssuedAt:   now.Unix(),
		NotBefore:  now.Unix(),
		Subject:    uid,
		PhoneModel: phoneModel,
		OS:         osStr,
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

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "missing Authorization", http.StatusUnauthorized)
			return
		}
		ctx, ok, err := ValidateToken(r.Context(), token)
		if err != nil {
			http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}
		if !ok {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
