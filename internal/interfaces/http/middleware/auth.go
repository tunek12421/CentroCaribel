package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/tunek/centro-caribel/internal/application/auth"
	apperrors "github.com/tunek/centro-caribel/pkg/errors"
	"github.com/tunek/centro-caribel/pkg/response"
)

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	RolNombreKey contextKey = "rol_nombre"
)

func AuthMiddleware(jwtSvc auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				response.Error(w, apperrors.NewUnauthorized("Token no proporcionado"))
				return
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				response.Error(w, apperrors.NewUnauthorized("Formato de token inválido"))
				return
			}

			claims, err := jwtSvc.ValidateToken(parts[1])
			if err != nil {
				response.Error(w, apperrors.NewUnauthorized("Token inválido o expirado"))
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, RolNombreKey, claims.RolNombre)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRoles(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rolNombre, ok := r.Context().Value(RolNombreKey).(string)
			if !ok {
				response.Error(w, apperrors.NewForbidden("Rol no encontrado en el contexto"))
				return
			}

			for _, allowed := range roles {
				if rolNombre == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}

			response.Error(w, apperrors.NewForbidden("No tiene permisos para esta acción"))
		})
	}
}

func GetUserID(ctx context.Context) string {
	if v, ok := ctx.Value(UserIDKey).(string); ok {
		return v
	}
	return ""
}

func GetRolNombre(ctx context.Context) string {
	if v, ok := ctx.Value(RolNombreKey).(string); ok {
		return v
	}
	return ""
}
