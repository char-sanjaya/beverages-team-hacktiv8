package mysql

import (
	"database/sql"
	"errors"

	"beverages-team-hacktiv8/domain"
)

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

// 1. Tarik semua data pesanan murni dari tabel 'order'
func (r *orderRepository) FetchAll() ([]domain.Order, error) {
	query := "SELECT id, invoice_number, user_id, total_price, status, created_at FROM `orders` ORDER BY created_at DESC"
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		// Scan HANYA 6 kolom asli sesuai ERD!
		if err := rows.Scan(&o.ID, &o.InvoiceNumber, &o.UserID, &o.TotalPrice, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

// 2. Tarik item pesanan murni dari tabel 'order_detail'
func (r *orderRepository) GetItemsByOrderID(orderID int) ([]domain.OrderDetail, error) {
	query := "SELECT id, order_id, product_id, quantity, price FROM order_detail WHERE order_id = ? ORDER BY id ASC"
	
	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.OrderDetail
	for rows.Next() {
		var it domain.OrderDetail
		// Scan HANYA 5 kolom asli sesuai ERD!
		if err := rows.Scan(&it.ID, &it.OrderID, &it.ProductID, &it.Quantity, &it.Price); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, nil
}

// 3. Update status pesanan
func (r *orderRepository) UpdateStatus(orderID int, status string) error {
	query := "UPDATE `orders` SET status = ? WHERE id = ?"
	res, err := r.db.Exec(query, status, orderID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("pesanan tidak ditemukan")
	}
	return nil
}

// 4. Hapus pesanan dengan Transaksi ACID
func (r *orderRepository) DeleteWithTransaction(orderID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Hapus order_detail
	_, err = tx.Exec("DELETE FROM order_detail WHERE order_id = ?", orderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Hapus order
	res, err := tx.Exec("DELETE FROM `orders` WHERE id = ?", orderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		tx.Rollback()
		return errors.New("pesanan tidak ditemukan")
	}

	return tx.Commit()
}