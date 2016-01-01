package main

// Item model
type Item struct {
	ID           int     `json:"id"`
	RestaurantID int     `json:"restaurant_id"`
	Name         string  `json:"name"`
	Image        string  `json:"image"`
	Price        float32 `json:"price"`
}
