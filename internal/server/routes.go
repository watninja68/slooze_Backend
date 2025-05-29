package server

import (
	"encoding/json"
	"log"
	"net/http"

	"backend/internal/handlers"
	"backend/internal/middleware" // Import the middleware package

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer) // Add a recoverer middleware

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}, // Added X-CSRF-Token
		ExposedHeaders:   []string{"Link"},                                                    // Added Link
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// --- Public Routes ---
	r.Get("/", s.HelloWorldHandler)
	r.Get("/health", s.healthHandler)
	r.Get("/users/{userId}/payment-methods", handlers.ListPaymentMethods)
	r.Post("/users/{userId}/payment-methods", handlers.AddPaymentMethod)
	r.Put("/payment-methods/{methodId}", handlers.UpdatePaymentMethod)
	// --- API v1 Group with Authentication ---
	r.Route("/api/v1", func(r chi.Router) {
		// Apply AuthMiddleware to all /api/v1 routes
		r.Use(middleware.AuthMiddleware)
		// Apply CountryCheck (or handle country logic in handlers)
		r.Use(middleware.CountryCheck)

		// --- Restaurants & Menus (All roles, country-filtered) --- [cite: 4, 8]
		r.Get("/restaurants", handlers.ListRestaurants)
		r.Get("/restaurants/{restaurantId}/menu", handlers.GetMenu)

		r.Route("/orders", func(r chi.Router) {
			r.Post("/", handlers.CreateOrder)
			r.Get("/", handlers.ListOrders)
			r.Get("/{orderId}", handlers.GetOrder)
			r.Post("/{orderId}/items", handlers.AddItemToOrder)
			r.Post("/{orderId}/checkout", handlers.CheckoutOrder)
			r.Post("/{orderId}/cancel", handlers.CancelOrder)
		})
		// --- Payment Methods ---
		r.Route("/payment-methods", func(r chi.Router) {
			// Update payment method (Admin only) [cite: 10]
			r.With(middleware.AdminOnly).Put("/{methodId}", handlers.UpdatePaymentMethod)
			// TODO: Add DELETE and GET for /me/payment-methods
		})

		// --- User-Specific Payment Methods ---
		// We'll add a /me route for users to manage their own stuff
		r.Route("/me", func(r chi.Router) {
			// TODO: Implement GET /me/payment-methods (User's own)
			// TODO: Implement POST /me/payment-methods (User's own)
		})

		// --- Admin-Specific User Payment Methods ---
		r.With(middleware.AdminOnly).Route("/users/{userId}/payment-methods", func(r chi.Router) {
			r.Get("/", handlers.ListPaymentMethods)
			r.Post("/", handlers.AddPaymentMethod)
			// Note: Update is Admin-only but a global route above based on strict table interpretation
		})
	})

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World - Slooze Food Order App"

	w.Header().Set("Content-Type", "application/json") // Set content type
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("error handling JSON marshal. Err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	health := s.db.Health()
	w.Header().Set("Content-Type", "application/json") // Set content type
	jsonResp, err := json.Marshal(health)
	if err != nil {
		log.Printf("error handling JSON marshal. Err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonResp)
}
