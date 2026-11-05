package usecase

import (
	"errors"
	"fmt"
	"strings"

	"beverages-team-hacktiv8/domain"
)

type productUsecase struct {
	prodRepo domain.ProductRepository
}

func NewProductUsecase(repo domain.ProductRepository) domain.ProductUsecase {
	return &productUsecase{prodRepo: repo}
}

// 1. Ambil Semua Produk
func (p *productUsecase) GetAllProducts() ([]domain.Product, error) {
	return p.prodRepo.FetchAll()
}

// 2. Ambil 1 Produk by ID
func (p *productUsecase) GetProduct(id int) (*domain.Product, error) {
	return p.prodRepo.GetByID(id)
}

// 3. Tambah Minuman Baru (Dengan Validasi)
func (p *productUsecase) AddProduct(name string, price float64, stock int) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("nama minuman tidak boleh kosong")
	}
	if price <= 0 {
		return errors.New("harga minuman harus lebih dari 0")
	}
	if stock < 0 {
		return errors.New("stok minuman tidak boleh bernilai negatif")
	}

	newProduct := &domain.Product{
		Name:  name,
		Price: price,
		Stock: stock,
	}
	return p.prodRepo.Create(newProduct)
}

// 4. Ubah Minuman (Dengan Validasi)
func (p *productUsecase) EditProduct(id int, name string, price float64, stock int) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("nama minuman tidak boleh kosong")
	}
	if price <= 0 {
		return errors.New("harga minuman harus lebih dari 0")
	}
	if stock < 0 {
		return errors.New("stok minuman tidak boleh bernilai negatif")
	}

	updated := &domain.Product{
		ID:    id,
		Name:  name,
		Price: price,
		Stock: stock,
	}
	return p.prodRepo.Update(updated)
}

// 5. Hapus Minuman (Dengan Pengecekan Relasi Pesanan)
func (p *productUsecase) RemoveProduct(id int) error {
	// Cek apakah minuman pernah dibeli di order_detail
	count, err := p.prodRepo.CountInOrderDetails(id)
	if err == nil && count > 0 {
		return fmt.Errorf("produk ID %d tidak boleh dihapus karena tersimpan pada %d riwayat transaksi pesanan", id, count)
	}
	return p.prodRepo.Delete(id)
}