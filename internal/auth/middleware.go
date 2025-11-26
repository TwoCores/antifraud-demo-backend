package auth

import "net/http"

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

func RequireSuperuserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := JwtClaimsFromContext(r)
		if !ok || !claims.IsSuperuser {
			http.Error(w, "superuser access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
