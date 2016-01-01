package main

import (
	"database/sql"
	"fmt"
)

// Item model
type Item struct {
	ID           int     `json:"id"`
	RestaurantID int     `json:"restaurant_id"`
	Name         string  `json:"name"`
	Image        string  `json:"image"`
	Price        float32 `json:"price"`
}

// GetItemsForRestaurantID returns a list of all restaurants
func GetItemsForRestaurantID(db *sql.DB, restaurantID int) ([]Item, error) {
	var items []Item

	statement := fmt.Sprintf("select * from items where restaurantID = %d", restaurantID)
	rs, err := db.Query(statement)
	if err != nil {
		fmt.Printf("Query Error: %s", err.Error())
		return nil, err
	}
	defer rs.Close()

	for rs.Next() {
		item := Item{}
		err = rs.Scan(&item.ID, &item.RestaurantID, &item.Name, &item.Image, &item.Price)
		if err != nil {
			fmt.Printf("DB Scan Error: %s\n", err.Error())
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
