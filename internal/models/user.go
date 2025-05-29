package models

import "context"

// ----- tiny models -----
type Role struct {
	ID   int
	Name string
}

type Country struct {
	ID   int
	Name string
}

type User struct {
	ID      int
	Name    string
	Email   string
	Role    Role
	Country Country
}

// GetUserByID is a placeholder so the project builds.
// Replace with real SQL / GORM / sqlc call later.
func GetUserByID(ctx context.Context, id int) (User, error) {
	return User{
		ID:      id,
		Name:    "Stub User",
		Role:    Role{Name: "ADMIN"},
		Country: Country{Name: "India"},
	}, nil
}
