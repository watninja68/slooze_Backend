package handlers

import (
	"backend/internal/database"
	"encoding/json"
	//"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ─── Request / Response DTOs ──────────────────────────────────────────────────

type PaymentMethodCreateRequest struct {
	MethodType string `json:"method_type"`          // e.g. credit_card, upi
	Details    string `json:"details"`              // token or encrypted blob
	IsDefault  bool   `json:"is_default,omitempty"` // optional
}

type PaymentMethodCreateResponse struct {
	PaymentMethodID int64  `json:"payment_method_id"`
	Message         string `json:"message"`
}

type PaymentMethodUpdateRequest struct {
	MethodType string `json:"method_type"`
	Details    string `json:"details"`
	IsDefault  bool   `json:"is_default"`
}

// ─── Handlers ────────────────────────────────────────────────────────────────

// GET /users/{userId}/payment-methods
func ListPaymentMethods(w http.ResponseWriter, r *http.Request) {
	//userIDStr := chi.URLParam(r, "userId")
	userIDStr := "1"
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "userId must be an integer", http.StatusBadRequest)
		return
	}

	svc := database.New()
	methods, err := svc.ListPaymentMethodsDB(r.Context(), userID)
	if err != nil {
		http.Error(w, "could not list payment methods: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, methods)
}

// POST /users/{userId}/payment-methods
func AddPaymentMethod(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "userId must be an integer", http.StatusBadRequest)
		return
	}

	var req PaymentMethodCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.MethodType == "" || req.Details == "" {
		http.Error(w, "`method_type` and `details` are required", http.StatusBadRequest)
		return
	}

	svc := database.New()
	id, err := svc.AddPaymentMethodDB(r.Context(), database.PaymentMethod{
		UserID:    userID,
		Type:      req.MethodType,
		Details:   req.Details,
		IsDefault: req.IsDefault,
	})
	if err != nil {
		http.Error(w, "could not add payment method: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := PaymentMethodCreateResponse{
		PaymentMethodID: id,
		Message:         "payment method added successfully",
	}
	writeJSON(w, http.StatusCreated, resp)
}

// PUT /payment-methods/{methodId}
func UpdatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	methodIDStr := chi.URLParam(r, "methodId")
	methodID, err := strconv.ParseInt(methodIDStr, 10, 64)
	if err != nil {
		http.Error(w, "methodId must be an integer", http.StatusBadRequest)
		return
	}

	var req PaymentMethodUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.MethodType == "" || req.Details == "" {
		http.Error(w, "`method_type` and `details` are required", http.StatusBadRequest)
		return
	}

	svc := database.New()
	if err := svc.UpdatePaymentMethodDB(r.Context(), database.PaymentMethod{
		ID:        methodID,
		Type:      req.MethodType,
		Details:   req.Details,
		IsDefault: req.IsDefault,
	}); err != nil {
		http.Error(w, "could not update payment method: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "payment method updated successfully",
	})
}
