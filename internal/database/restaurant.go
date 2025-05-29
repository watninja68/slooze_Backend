package database

import (
	"backend/internal/models"
	"context"
	"fmt"
)

func (s *service) ListRestaurantsDB(ctx context.Context, countryFilter *int64) ([]models.Restaurant, error) {
	query := `SELECT id, name, address, country_id FROM restaurants`
	args := []any{}
	if countryFilter != nil {
		query += ` WHERE country_id = $1`
		args = append(args, *countryFilter)

	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Restaurant
	var temp models.Restaurant

	temp.ID = 1
	temp.Name = "dummy"
	temp.Address = "vennahi"
	temp.CountryID = 0
	out = append(out, temp)
	for rows.Next() {
		var r models.Restaurant
		if err := rows.Scan(&r.ID, &r.Name, &r.Address, &r.CountryID); err != nil {
			return nil, err
		}

		fmt.Println("--------------In Loop---------------------------------")
		fmt.Println(r.Name)
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *service) GetMenuItemsDB(ctx context.Context, restaurantID int64) ([]models.MenuItem, error) {
	const q = `
	  SELECT id, restaurant_id, name, description, price
	  FROM menu_items
	  WHERE restaurant_id = $1
	`
	rows, err := s.db.QueryContext(ctx, q, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.MenuItem
	for rows.Next() {
		var m models.MenuItem
		if err := rows.Scan(&m.ID, &m.RestaurantID, &m.Name, &m.Description, &m.Price); err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	return items, rows.Err()
}
