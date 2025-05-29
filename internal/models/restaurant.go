package models

type Restaurant struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address,omitempty"`
	CountryID int64  `json:"country_id"`
}

type MenuItem struct {
	ID           int64   `json:"id"`
	RestaurantID int64   `json:"-"`
	Name         string  `json:"name"`
	Description  string  `json:"description,omitempty"`
	Price        float64 `json:"price"`
}
