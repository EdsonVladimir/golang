package middleware

import (
	"context"
	"edson.com/go/rest-ws/models"
	"edson.com/go/rest-ws/server"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

var (
	NO_AUTH_NEEDED = []string{
		"login",
		"signup",
	}
)

func shouldCheckToken(route string) bool {
	for _, p := range NO_AUTH_NEEDED {
		if strings.Contains(route, p) {
			return false
		}
	}
	return true
}

func CheckAuthMiddleware(s server.Server) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verificamos si debemos revisar el token para esta ruta
			if !shouldCheckToken(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
			claims := &models.AppClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(s.Config().JWTSecret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "auth_user_id", claims.UserId)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
