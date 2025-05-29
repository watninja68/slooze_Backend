package middleware

import (
	"context"
	"net/http"
	"strings"

	"backend/internal/auth"
	"backend/internal/models" // change to your DAL package
)

// UserCtxKey / AuthUser unchanged
type UserCtxKey string

const CtxUserKey UserCtxKey = "user"

// AuthMiddleware validates JWT, refreshes user from DB, sets context.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// ---- 1. Extract "Authorization: Bearer <token>" header ----
		h := r.Header.Get("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(h, "Bearer ")

		// ---- 2. Validate + get claims ----
		claims, err := auth.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// ---- 3. Load the latest user from DB (optional but recommended) ----
		u, err := models.GetUserByID(r.Context(), claims.UserID) // implement in DAL
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		authUser := AuthUser{
			ID:      u.ID,
			Role:    u.Role.Name,    // assuming joined role
			Country: u.Country.Name, // assuming joined country
		}

		// ---- 4. Inject into request context ----
		ctx := context.WithValue(r.Context(), CtxUserKey, authUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
