package handlers

import (
	"context"
	"net/http"

	"storeSystem/internal/auth"
)

type contextKey string

const UserClaimsKey contextKey = "userClaims"

func (h *Handlers) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "токен отсутствует")
			return
		}

		if cookie.Value == "" {
			respondWithError(w, http.StatusUnauthorized, "пустой токен")
			return
		}

		claims, err := auth.ParseToken(cookie.Value)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "невалидный токен")
			return
		}

		ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserClaimsFromContext(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(UserClaimsKey).(*auth.Claims)
	return claims, ok
}

func RequireAdmin() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetUserClaimsFromContext(r.Context())
			if !ok {
				respondWithError(w, http.StatusUnauthorized, "не авторизован")
				return
			}

			if claims.RoleID == 1 {
				next.ServeHTTP(w, r)
				return
			}
			respondWithError(w, http.StatusForbidden, "доступ запрещён")
		})
	}
}
