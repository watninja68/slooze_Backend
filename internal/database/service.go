package database

import (
	"context"
	_ "database/sql"
	"time"
)

// ─── Domain Model ─────────────────────────────────────────────────────────────

type PaymentMethod struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Type      string    `json:"method_type"`
	Details   string    `json:"details"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ─── CRUD Helpers ─────────────────────────────────────────────────────────────

// ListPaymentMethodsDB returns all methods for a given user (default first).
func (s *service) ListPaymentMethodsDB(ctx context.Context, userID int64) ([]PaymentMethod, error) {
	const q = `
		SELECT id, user_id, method_type, details, is_default, created_at, updated_at
		FROM payment_methods
		WHERE user_id = $1
		ORDER BY is_default DESC, created_at ASC;`
	rows, err := s.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []PaymentMethod
	for rows.Next() {
		var pm PaymentMethod
		if err := rows.Scan(&pm.ID, &pm.UserID, &pm.Type, &pm.Details,
			&pm.IsDefault, &pm.CreatedAt, &pm.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, pm)
	}
	return out, rows.Err()
}

// AddPaymentMethodDB inserts a new record and returns its id.
// If IsDefault=true it atomically clears other defaults first.
func (s *service) AddPaymentMethodDB(ctx context.Context, pm PaymentMethod) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	if pm.IsDefault {
		if _, err := tx.ExecContext(ctx,
			`UPDATE payment_methods SET is_default = FALSE WHERE user_id = $1`, pm.UserID); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	const ins = `
		INSERT INTO payment_methods (user_id, method_type, details, is_default)
		VALUES ($1,$2,$3,$4) RETURNING id`
	var id int64
	if err := tx.QueryRowContext(ctx, ins, pm.UserID, pm.Type, pm.Details, pm.IsDefault).Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, tx.Commit()
}

// UpdatePaymentMethodDB updates an existing record by id.
func (s *service) UpdatePaymentMethodDB(ctx context.Context, pm PaymentMethod) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if pm.IsDefault {
		// keep only one default per user
		if _, err := tx.ExecContext(ctx, `
			UPDATE payment_methods
			SET is_default = FALSE
			WHERE user_id = (SELECT user_id FROM payment_methods WHERE id = $1)`, pm.ID); err != nil {
			tx.Rollback()
			return err
		}
	}
	const upd = `
		UPDATE payment_methods
		SET method_type = $1,
		    details     = $2,
		    is_default  = $3,
		    updated_at  = NOW()
		WHERE id = $4`
	if _, err := tx.ExecContext(ctx, upd, pm.Type, pm.Details, pm.IsDefault, pm.ID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
