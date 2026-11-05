package domain

import "time"

// Pilihan ENUM Status
const (
	OrderStatusPending = "pending"
	OrderStatusSuccess = "success"
)

// 1. Entity Order: 100% persis kolom tabel 'order' di ERD
type Order struct {
	ID            int
	InvoiceNumber string
	UserID        int
	TotalPrice    float64
	Status        string
	CreatedAt     time.Time
}

// 2. Entity OrderDetail: 100% persis kolom tabel 'order_detail' di ERD
type OrderDetail struct {
	ID        int
	OrderID   int
	ProductID int
	Quantity  int
	Price     float64
}

// Kontrak fungsi database untuk pesanan
type OrderRepository interface {
	FetchAll() ([]Order, error)
	GetItemsByOrderID(orderID int) ([]OrderDetail, error)
	UpdateStatus(orderID int, status string) error
	DeleteWithTransaction(orderID int) error
}

// Kontrak aturan bisnis pesanan
type OrderUsecase interface {
	GetAllOrders() ([]Order, error)
	GetOrderDetails(orderID int) ([]OrderDetail, error)
	ChangeStatus(orderID int, status string) error
	DeleteOrder(orderID int) error
}