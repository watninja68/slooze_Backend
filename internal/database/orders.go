package database

import (
	"context"
	"fmt"
)

// Order is a slim projection used by the service layer.
// Extend as you need, or create separate DTOs for handlers.
type Order struct {
	ID           int64   `json:"id"`
	UserID       int64   `json:"user_id"`
	RestaurantID int64   `json:"restaurant_id"`
	TotalPrice   float64 `json:"total_price"`
	Status       string  `json:"status"`
	CreatedAt    string  `json:"created_at"`
}

// CreateOrderDB inserts a new row into the orders table and
// returns the auto-generated primary key.
func (s *service) CreateOrderDB(
	ctx context.Context,
	userID int64,
	restaurantID int64,
	totalPrice float64,
) (int64, error) {

	const q = `
		INSERT INTO orders (user_id, restaurant_id, total_price)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	var id int64
	if err := s.db.QueryRowContext(ctx, q, userID, restaurantID, totalPrice).Scan(&id); err != nil {
		return 0, fmt.Errorf("insert order: %w", err)
	}

	return id, nil
}

// ListOrdersDB fetches every order (demo version – add paging / filters later).
func (s *service) ListOrdersDB(ctx context.Context) ([]Order, error) {
	const q = `
		SELECT id, user_id, restaurant_id, total_price, status, created_at
		FROM orders
		ORDER BY created_at DESC;
	`
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}
	defer rows.Close()

	var out []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.RestaurantID,
			&o.TotalPrice,
			&o.Status,
			&o.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

// GetOrderDB returns one order by id.
func (s *service) GetOrderDB(ctx context.Context, id int64) (Order, error) {
	const q = `
		SELECT id, user_id, restaurant_id, total_price, status, created_at
		FROM orders
		WHERE id = $1;
	`
	var o Order
	if err := s.db.QueryRowContext(ctx, q, id).Scan(
		&o.ID,
		&o.UserID,
		&o.RestaurantID,
		&o.TotalPrice,
		&o.Status,
		&o.CreatedAt,
	); err != nil {
		return Order{}, fmt.Errorf("get order %d: %w", id, err)
	}
	return o, nil
}
