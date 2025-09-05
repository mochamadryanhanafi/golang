package middleware

import (
	"auth-service/utils"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey = contextKey("userID")

// JWTMiddleware memvalidasi token JWT dari header Authorization.
func JWTMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.WriteError(w, http.StatusUnauthorized, "authorization header missing")
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenStr == authHeader {
				utils.WriteError(w, http.StatusUnauthorized, "invalid token format")
				return
			}

			userID, err := utils.ValidateJWT(tokenStr, jwtSecret)
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, "invalid token: "+err.Error())
				return
			}

			// inject user ID ke dalam context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
