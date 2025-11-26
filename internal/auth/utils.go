package auth

import (
	"net/http"
)

func JwtClaimsFromContext(r *http.Request) (*JwtClaims, bool) {
	v := r.Context().Value(CtxKeyClaims)
	if v == nil {
		return nil, false
	}

	claims, ok := v.(*JwtClaims)
	if !ok {
		return nil, false
	}

	return claims, true
}
