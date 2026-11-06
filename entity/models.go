package entity

import "time"

type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	Role      string
	Balance   float64
	CreatedAt time.Time
}

type Product struct {
	ID        int
	Name      string
	Price     float64
	Stock     int
	CreatedAt time.Time
}

type Order struct {
	ID            int
	InvoiceNumber string
	UserID        int
	TotalPrice    float64
	Status        string
	CreatedAt     time.Time
}

// CartItem menyimpan rincian order sementara di memori CLI sebelum masuk DB
type CartItem struct {
	ProductID   int
	ProductName string
	Price       float64
	Quantity    int
	Subtotal    float64
}
