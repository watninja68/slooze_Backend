package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Placeholder: In a real app, this would fetch restaurants, applying country filters based on user role.
func ListRestaurants(w http.ResponseWriter, r *http.Request) {
	// TODO: Get user role/country from context (set by middleware)
	// TODO: Implement DB logic to fetch restaurants with filtering
	resp := map[string]string{"message": "List of restaurants (filtered by country if needed)"}
	json.NewEncoder(w).Encode(resp)
}

// Placeholder: In a real app, this would fetch menu items for a restaurant, applying country filters.
func GetMenu(w http.ResponseWriter, r *http.Request) {
	restaurantID := chi.URLParam(r, "restaurantId")
	// TODO: Get user role/country from context
	// TODO: Implement DB logic to fetch menu items with filtering
	resp := map[string]string{"message": "Menu for restaurant " + restaurantID + " (filtered by country if needed)"}
	json.NewEncoder(w).Encode(resp)
}
