package internal

import (
	"antifraud-demo-backend/internal/auth"
	"net/http"
)

func jwtClaimsFromContext(r *http.Request) (*auth.JwtClaims, bool) {
	v := r.Context().Value(auth.CtxKeyClaims)
	if v == nil {
		return nil, false
	}
	claims, ok := v.(*auth.JwtClaims)
	if !ok {
		return nil, false
	}
	return claims, true
}
