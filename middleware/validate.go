package middleware

import (
	"go-fiber-api/internal/shared/types"
	utils "go-fiber-api/utils/jwt"
	"go-fiber-api/utils/web"
	"log/slog"
	"net/http"
	"strings"
)

func ValidateRole(requiredRole ...types.Roles) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			token := strings.TrimPrefix(authHeader, "Bearer ")
			token = strings.TrimSpace(token)

			if token == "" {
				web.Err(w, web.NewHTTPError(http.StatusUnauthorized, "Unauthorized", web.ErrAuthentication))
				return
			}

			claims, err := utils.ParseJWTToken(token)
			if err != nil || claims == nil {
				web.Err(w, web.NewHTTPError(http.StatusUnauthorized, "Unauthorized", web.ErrAuthentication))
				return
			}

			slog.Info("Request info", "path", r.URL.Path, "method", r.Method, "userID", claims.ID, "role", claims.Role)

			// Role 0 (misal: Super Admin) bypass semua
			if claims.Role == 0 {
				next.ServeHTTP(w, r)
				return
			}

			userRole := claims.Role
			for _, role := range requiredRole {
				if userRole == int(role) {
					next.ServeHTTP(w, r)
					return
				}
			}

			web.Err(w, web.NewHTTPError(http.StatusForbidden, "Forbidden", web.ErrForbidden))
		})
	}
}

func ValidateRoleAdmin() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			token := strings.TrimPrefix(authHeader, "Bearer ")
			token = strings.TrimSpace(token)

			if token == "" {
				web.Err(w, web.NewHTTPError(http.StatusUnauthorized, "Unauthorized", web.ErrAuthentication))
				return
			}

			claims, err := utils.ParseJWTToken(token)
			if err != nil || claims == nil {
				web.Err(w, web.NewHTTPError(http.StatusUnauthorized, "Unauthorized", web.ErrAuthentication))
				return
			}

			slog.Info("Request info", "path", r.URL.Path, "method", r.Method, "userID", claims.ID, "role", claims.Role)

			// Cek apakah role-nya admin (misal: role == 1)
			if claims.Role != int(types.RoleAdmin) {
				web.Err(w, web.NewHTTPError(http.StatusForbidden, "Forbidden", web.ErrForbidden))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
