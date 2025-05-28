package middleware

import (
	"context"
	"net/http"
)

// UserCtxKey is a key for storing user info in the request context.
type UserCtxKey string

const CtxUserKey UserCtxKey = "user"

// Placeholder User struct (adapt as needed)
type AuthUser struct {
	ID      int
	Role    string // e.g., "ADMIN", "MANAGER", "MEMBER"
	Country string // e.g., "India", "America"
}

// Placeholder Authentication Middleware
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement actual token validation (e.g., JWT)
		// 1. Get token from header
		// 2. Validate token
		// 3. Fetch user details (ID, role, country) from DB or token claims
		// 4. If valid, add user to context and call next.ServeHTTP
		// 5. If invalid, return http.StatusUnauthorized

		// --- Placeholder ---
		// For demonstration, let's assume a valid user is always present
		// In a real app, you MUST implement real validation.
		// You might fetch this based on an API key or token.
		// We'll hardcode a Manager from India for now.
		user := AuthUser{ID: 2, Role: "MANAGER", Country: "India"}

		// Add user to context
		ctx := context.WithValue(r.Context(), CtxUserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper to get user from context
func GetUserFromContext(r *http.Request) (AuthUser, bool) {
	user, ok := r.Context().Value(CtxUserKey).(AuthUser)
	return user, ok
}

// Placeholder RBAC Middleware - Admin Only
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok || user.Role != "ADMIN" {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Placeholder RBAC Middleware - Manager or Admin
func ManagerOrAdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok || (user.Role != "MANAGER" && user.Role != "ADMIN") {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Placeholder Country Check Middleware (Can be integrated into handlers or separate)
// This is more complex as it depends on the resource being accessed.
// Often, it's better handled within the handler or database layer.
// Example:
func CountryCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// If user is Admin, they can access anything
		if user.Role == "ADMIN" {
			next.ServeHTTP(w, r)
			return
		}

		// TODO: For Manager/Member, check the resource's country.
		// This requires knowing WHAT is being accessed (e.g., restaurantId, orderId)
		// and fetching its country from the DB, then comparing.
		// This is a simplified example and needs real implementation.
		// resourceCountry := "India" // Fetch this based on URL params/DB
		// if user.Country == resourceCountry {
		// 	next.ServeHTTP(w, r)
		// } else {
		// 	http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		// }

		// For now, let's just pass through and handle in handlers.
		next.ServeHTTP(w, r)
	})
}
