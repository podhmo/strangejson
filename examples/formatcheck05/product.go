package model

import "time"

// Product :
type Product struct {
	Name  string `json:"name" required:"true" minLength:"1" maxLength:"255"`
	Price int    `json:"price"`
}

// Item :
type Item struct {
	Product Product `json:"product"`
	Count   int     `json:"count"`
}

// Order :
type Order struct {
	OrderedAt time.Time `json:"orderedAt"`
	Items     []Item    `json:"items"`
}
