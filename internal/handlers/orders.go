package handlers

import (
	"backend/internal/database"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ─── Request / Response DTOs ──────────────────────────────────────────────────

type OrderCreateRequest struct {
	UserID       int64   `json:"user_id"`
	RestaurantID int64   `json:"restaurant_id"`
	TotalPrice   float64 `json:"total_price"`
}

type OrderCreateResponse struct {
	OrderID int64  `json:"order_id"`
	Message string `json:"message"`
}

// ─── Handlers ────────────────────────────────────────────────────────────────

// POST /orders
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req OrderCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body: "+err.Error(), http.StatusBadRequest)
		return
	}

	svc := database.New()
	id, err := svc.CreateOrderDB(r.Context(), req.UserID, req.RestaurantID, req.TotalPrice)
	if err != nil {
		http.Error(w, "could not create order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := OrderCreateResponse{
		OrderID: id,
		Message: "order created successfully",
	}
	writeJSON(w, http.StatusCreated, resp)
}

// GET /orders
func ListOrders(w http.ResponseWriter, r *http.Request) {
	svc := database.New()
	orders, err := svc.ListOrdersDB(r.Context())
	if err != nil {
		http.Error(w, "could not list orders: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, orders)
}

// GET /orders/{orderId}
func GetOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "orderId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "orderId must be an integer", http.StatusBadRequest)
		return
	}

	svc := database.New()
	order, err := svc.GetOrderDB(r.Context(), id)
	if err != nil {
		http.Error(w, "could not fetch order: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, order)
}

// stubs (keep your existing TODO comments) ---------------------------
func AddItemToOrder(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
func CheckoutOrder(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
func CancelOrder(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// ─── helper ───────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

