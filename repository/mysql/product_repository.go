package mysql

import (
	"beverages-team-hacktiv8/domain"
	"database/sql"
	"errors"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

// 1. Tarik seluruh data minuman
func (r *productRepository) FetchAll() ([]domain.Product, error) {
	query := "SELECT id, name, price, stock, created_at FROM product ORDER BY id ASC"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Wajib ditutup agar koneksi tidak bocor!

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// 2. Ambil 1 produk berdasarkan ID
func (r *productRepository) GetByID(id int) (*domain.Product, error) {
	query := "SELECT id, name, price, stock, created_at FROM product WHERE id = ?"
	var p domain.Product
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// 3. Tambah minuman baru
func (r *productRepository) Create(p *domain.Product) error {
	query := "INSERT INTO product (name, price, stock) VALUES (?, ?, ?)"
	res, err := r.db.Exec(query, p.Name, p.Price, p.Stock)
	if err != nil {
		return err
	}
	lastID, _ := res.LastInsertId()
	p.ID = int(lastID)
	return nil
}

// 4. Update data minuman
func (r *productRepository) Update(p *domain.Product) error {
	query := "UPDATE product SET name = ?, price = ?, stock = ? WHERE id = ?"
	res, err := r.db.Exec(query, p.Name, p.Price, p.Stock, p.ID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}
	return nil
}

// 5. Hapus minuman
func (r *productRepository) Delete(id int) error {
	res, err := r.db.Exec("DELETE FROM product WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}
	return nil
}

// 6. Cek apakah produk ini ada di riwayat pesanan (order_detail)
func (r *productRepository) CountInOrderDetails(productID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM order_detail WHERE product_id = ?"
	err := r.db.QueryRow(query, productID).Scan(&count)
	return count, err
}