package models

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
