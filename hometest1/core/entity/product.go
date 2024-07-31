package entity

import "time"

type Product struct {
	ID        int64
	Serial    string
	Name      string
	Price     float64
	UpdatedAt time.Time
}

type ProductQuantity struct {
	ID        int64
	ProductID int64
	Quantity  int
	UpdatedAt time.Time
}
