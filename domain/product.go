package domain

import "time"

// Entity product sebagai blueprint data product
type Product struct {
	ID        int
	Name      string
	Price     float64
	Stock     int
	CreatedAt time.Time
}

// kontrak fungsi database untuk product
type ProductRepository interface{
	FetchAll() ([]Product, error)
	GetByID(id int) (*Product, error)
	Create(p *Product) error
	Update(p *Product) error
	Delete(id int) error
	CountInOrderDetails(productID int) (int, error)
}

// Kontrak logika bisnis product
type ProductUsecase interface{
	GetAllProducts() ([]Product, error)
	GetProduct(id int) (*Product, error)
	AddProduct(name string, price float64, stock int) error
	EditProduct(id int, name string, price float64, stock int) error
	RemoveProduct(id int) error
}

