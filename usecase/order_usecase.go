package usecase

import (
	"errors"
	"strings"

	"beverages-cli/domain"
)

type orderUsecase struct {
	orderRepo domain.OrderRepository
}

func NewOrderUsecase(repo domain.OrderRepository) domain.OrderUsecase {
	return &orderUsecase{orderRepo: repo}
}

// 1. Tampilkan Semua Pesanan
func (o *orderUsecase) GetAllOrders() ([]domain.Order, error) {
	return o.orderRepo.FetchAll()
}

// 2. Tampilkan Item di dalam Pesanan
func (o *orderUsecase) GetOrderDetails(orderID int) ([]domain.OrderDetail, error) {
	return o.orderRepo.GetItemsByOrderID(orderID)
}

// 3. Ubah Status Pesanan (Dengan Validasi ENUM!)
func (o *orderUsecase) ChangeStatus(orderID int, status string) error {
	status = strings.ToLower(strings.TrimSpace(status))

	// ATURAN BISNIS: Hanya izinkan status 'pending' atau 'success'
	if status != domain.OrderStatusPending && status != domain.OrderStatusSuccess {
		return errors.New("status tidak valid! Hanya boleh 'pending' atau 'success'")
	}

	return o.orderRepo.UpdateStatus(orderID, status)
}

// 4. Hapus Pesanan Transaksi
func (o *orderUsecase) DeleteOrder(orderID int) error {
	return o.orderRepo.DeleteWithTransaction(orderID)
}
