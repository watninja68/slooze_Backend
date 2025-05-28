package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Placeholder for listing payment methods
func ListPaymentMethods(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	// TODO: Check if current user can view this (self or Admin)
	// TODO: Implement DB logic
	resp := map[string]string{"message": "Payment methods for user " + userID}
	json.NewEncoder(w).Encode(resp)
}

// Placeholder for adding a payment method
func AddPaymentMethod(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	// TODO: Check if current user can add (self or Admin)
	// TODO: Decode body
	// TODO: Implement DB logic
	resp := map[string]string{"message": "Payment method added for user " + userID}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// Placeholder for updating a payment method
func UpdatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	methodID := chi.URLParam(r, "methodId")
	// TODO: Check if user is Admin (Strict interpretation)
	// TODO: Decode body
	// TODO: Implement DB logic
	resp := map[string]string{"message": "Payment method " + methodID + " updated"}
	json.NewEncoder(w).Encode(resp)
}
