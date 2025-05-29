package middleware

import (
	"backend/internal/models"
	"net/http"
	"strings"
)

// AdminOnly – short-circuit if role ≠ ADMIN
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, ok := r.Context().Value(CtxUserKey).(AuthUser)
		if !ok || strings.ToUpper(u.Role) != models.RoleAdmin {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CountryCheck – for bonus objective: manager/member can only touch own country.
// (Handlers that already check countryID will pass straight through.)
func CountryCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// At this stage we only need the user in context;
		// real resource-level checks live in the handlers.
		if _, ok := r.Context().Value(CtxUserKey).(AuthUser); !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
