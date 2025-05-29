package middleware

import (
	"context"
	"net/http"
	"strings"

	"backend/internal/auth"
	"backend/internal/database"
)

type AuthUser struct {
	ID        int
	Role      string
	Country   string
	CountryID int64
}

// UserCtxKey / AuthUser unchanged
type UserCtxKey string

const CtxUserKey UserCtxKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Header
		h := r.Header.Get("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(h, "Bearer ")

		// 2. Verify/parse
		claims, err := auth.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// 3. Refresh user from DB (no stub!)
		svc := database.New()
		u, err := svc.GetUserByID(r.Context(), claims.UserID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		authUser := AuthUser{
			ID:        u.ID,
			Role:      u.Role.Name,
			Country:   u.Country.Name,
			CountryID: int64(u.Country.ID),
		}

		// 4. inject
		ctx := context.WithValue(r.Context(), CtxUserKey, authUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
