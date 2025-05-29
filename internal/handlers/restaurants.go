package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"backend/internal/database"
	"backend/internal/models"
)

func ListRestaurants(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)

	// Admins see all, others only their country
	var countryFilter *int64
	if user.Role != models.RoleAdmin {
		countryFilter = &user.CountryID
	}

	restaurants, err := database.New().ListRestaurantsDB(r.Context(), countryFilter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(restaurants)
}

func GetMenu(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)
	restaurantID, _ := strconv.ParseInt(chi.URLParam(r, "restaurantId"), 10, 64)

	// quick country boundary check for non-admins
	if user.Role != models.RoleAdmin {
		rest, err := database.New().ListRestaurantsDB(r.Context(), &user.CountryID) // returns only user’s country
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
