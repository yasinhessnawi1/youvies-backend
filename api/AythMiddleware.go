package api

import (
	"context"
	"net/http"
	"strings"
	"youvies-backend/utils"
)

type ContextKey string

const UserKey ContextKey = "user"

func AuthMiddleware(next http.Handler, role string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Forbidden no key provided", http.StatusForbidden)
			return
		}

		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Forbidden, token not valid"+tokenString, http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), claims.Username, claims.Role)

		switch role {
		case "admin":
			if claims.Role != "admin" {
				http.Error(w, "Forbidden, admin role required", http.StatusForbidden)
				return
			}
		case "user":
			if claims.Role != "user" {
				if claims.Role != "admin" {
					http.Error(w, "Forbidden, user role required", http.StatusForbidden)
					return
				}
			}

		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
