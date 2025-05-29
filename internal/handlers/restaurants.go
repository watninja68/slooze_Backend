package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings" // Added for strings.ToUpper

	"github.com/go-chi/chi/v5"

	"backend/internal/database"
	"backend/internal/middleware" // Added for middleware.CtxUserKey and middleware.AuthUser
	"backend/internal/models"
)

func ListRestaurants(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(middleware.CtxUserKey).(middleware.AuthUser) // Changed user to u, and type to middleware.AuthUser

	// Admins see all, others only their country
	var countryFilter *int64
	if strings.ToUpper(u.Role) != models.RoleAdmin { // Changed user.Role to strings.ToUpper(u.Role)
		countryFilter = &u.CountryID // Changed user.CountryID to u.CountryID
	}

	restaurants, err := database.New().ListRestaurantsDB(r.Context(), countryFilter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(restaurants)
}

func GetMenu(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(middleware.CtxUserKey).(middleware.AuthUser) // Changed user to u, and type to middleware.AuthUser
	restaurantID, _ := strconv.ParseInt(chi.URLParam(r, "restaurantId"), 10, 64)

	// quick country boundary check for non-admins
	if strings.ToUpper(u.Role) != models.RoleAdmin { // Changed user.Role to strings.ToUpper(u.Role)
		rest, err := database.New().ListRestaurantsDB(r.Context(), &u.CountryID) // Changed user.CountryID to u.CountryID
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		valid := false
		for _, rs := range rest {
			if rs.ID == restaurantID {
				valid = true
				break
			}
		}
		if !valid {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
	}

	menu, err := database.New().GetMenuItemsDB(r.Context(), restaurantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(menu)
}

