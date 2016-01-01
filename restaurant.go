package main

import (
	"database/sql"
	"fmt"
)

// Restaurant model
type Restaurant struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

// GetAllRestaurants returns a list of all restaurants
func GetAllRestaurants(db *sql.DB) ([]Restaurant, error) {
	var restaurants []Restaurant

	rs, err := db.Query("select * from restaurants")
	if err != nil {
		fmt.Printf("Query Error: %s", err.Error())
		return nil, err
	}
	defer rs.Close()

	for rs.Next() {
		restaurant := Restaurant{}
		err = rs.Scan(&restaurant.ID, &restaurant.Name, &restaurant.Image)
		if err != nil {
			fmt.Printf("DB Scan Error: %s\n", err.Error())
			return nil, err
		}
		restaurants = append(restaurants, restaurant)
	}

	return restaurants, nil
}
