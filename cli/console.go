package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"beverages-cli/domain" // sesuaikan nama modul Anda
)

type ConsoleHandler struct {
	userUsecase    domain.UserUsecase
	productUsecase domain.ProductUsecase
	orderUsecase   domain.OrderUsecase
	scanner        *bufio.Scanner
	currentAdmin   *domain.User // Menyimpan session admin yang sedang login
}

// Constructor: Menerima 3 usecase sekaligus
func NewConsoleHandler(
	uu domain.UserUsecase,
	pu domain.ProductUsecase,
	ou domain.OrderUsecase,
) *ConsoleHandler {
	return &ConsoleHandler{
		userUsecase:    uu,
		productUsecase: pu,
		orderUsecase:   ou,
		scanner:        bufio.NewScanner(os.Stdin),
	}
}

// ======================== HELPER INPUT TERMINAL ======================== //

// Membaca satu baris utuh beserta spasi
func (h *ConsoleHandler) inputString(label string) string {
	fmt.Print(label)
	if h.scanner.Scan() {
		return strings.TrimSpace(h.scanner.Text())
	}
	return ""
}

// Membaca input angka bulat
func (h *ConsoleHandler) inputInt(label string) (int, error) {
	val, err := strconv.Atoi(h.inputString(label))
	if err != nil {
		return 0, fmt.Errorf("input harus berupa angka")
	}
	return val, nil
}

// Membaca input nominal harga
func (h *ConsoleHandler) inputFloat(label string) (float64, error) {
	val, err := strconv.ParseFloat(h.inputString(label), 64)
	if err != nil {
		return 0, fmt.Errorf("input harus berupa nominal angka valid")
	}
	return val, nil
}

// ======================== LOOP UTAMA PROGRAM ======================== //

func (h *ConsoleHandler) Start() {
	fmt.Println("====================================================")
	fmt.Println("   BEVERAGE E-COMMERCE (GOLANG CLEAN ARCHITECTURE)  ")
	fmt.Println("====================================================")

	for {
		if h.currentAdmin == nil {
			h.showAuthMenu()
		} else {
			h.showAdminDashboard()
		}
	}
}

// StartAdminSession menjalankan dashboard admin langsung dari sesi login terpusat
func (h *ConsoleHandler) StartAdminSession(admin *domain.User) {
	h.currentAdmin = admin
	fmt.Printf("[SUKSES] Selamat datang, %s! (ID: %d)\n", admin.Name, admin.ID)
	for h.currentAdmin != nil {
		h.showAdminDashboard()
	}
}

// ======================== 1. MENU AUTENTIKASI ======================== //

func (h *ConsoleHandler) showAuthMenu() {
	fmt.Println("\n---------------- MENU AUTENTIKASI ADMIN ----------------")
	fmt.Println("1. Login Akun Admin")
	fmt.Println("2. Registrasi Admin Baru")
	fmt.Println("0. Keluar dari Program")
	fmt.Println("--------------------------------------------------------")

	switch h.inputString("Pilih menu [0-2]: ") {
	case "1":
		h.handleLogin()
	case "2":
		h.handleRegister()
	case "0":
		fmt.Println("\nTerima kasih telah menggunakan sistem ini. Sampai jumpa!")
		os.Exit(0)
	default:
		fmt.Println("[Peringatan] Pilihan tidak valid.")
	}
}

func (h *ConsoleHandler) handleLogin() {
	fmt.Println("\n>>> FORM LOGIN ADMIN <<<")
	email := h.inputString("Email    : ")
	password := h.inputString("Password : ")

	admin, err := h.userUsecase.LoginAdmin(email, password)
	if err != nil {
		fmt.Printf("[GAGAL LOGIN] %v\n", err)
		return
	}

	h.currentAdmin = admin
	fmt.Printf("\n[SUKSES] Selamat datang, %s! (ID: %d)\n", admin.Name, admin.ID)
}

func (h *ConsoleHandler) handleRegister() {
	fmt.Println("\n>>> REGISTRASI ADMIN BARU <<<")
	name := h.inputString("Nama     : ")
	email := h.inputString("Email    : ")
	password := h.inputString("Password : ")

	err := h.userUsecase.RegisterAdmin(name, email, password)
	if err != nil {
		fmt.Printf("[GAGAL] %v\n", err)
		return
	}

	fmt.Println("[SUKSES] Admin baru berhasil didaftarkan. Silakan login pada menu 1.")
}

// ======================== 2. DASHBOARD SETELAH LOGIN ======================== //

func (h *ConsoleHandler) showAdminDashboard() {
	fmt.Println("\n====================================================")
	fmt.Printf(" DASHBOARD ADMIN: %s (%s)\n", h.currentAdmin.Name, h.currentAdmin.Email)
	fmt.Println("====================================================")
	fmt.Println(" [PRODUK MINUMAN]")
	fmt.Println("  1. Lihat Katalog Produk (Read)")
	fmt.Println("  2. Tambah Minuman Baru (Create)")
	fmt.Println("  3. Ubah Minuman (Update)")
	fmt.Println("  4. Hapus Minuman (Delete)")
	fmt.Println("\n [PESANAN MASUK (TRANSAKSI USER)]")
	fmt.Println("  5. Lihat Semua Pesanan (Read)")
	fmt.Println("  6. Lihat Detail Item Pesanan (Read Detail)")
	fmt.Println("  7. Update Status Pesanan [pending/success] (Update)")
	fmt.Println("  8. Hapus Pesanan Transaksi (Delete)")
	fmt.Println("\n [MANAJEMEN AKUN ADMIN SENDIRI]")
	fmt.Println("  9. Lihat Profil Akun Saya (Read)")
	fmt.Println("  10. Ubah Profil Saya (Update)")
	fmt.Println("  11. Hapus Akun Saya Sendiri (Delete)")
	fmt.Println("  12. Registrasi Admin Baru (Create)")
	fmt.Println("\n [SISTEM]")
	fmt.Println("  13. Logout")
	fmt.Println("  0. Keluar dari Program")
	fmt.Println("====================================================")

	switch h.inputString("Pilih menu [0-13]: ") {
	case "1":
		h.listProducts()
	case "2":
		h.createProduct()
	case "3":
		h.updateProduct()
	case "4":
		h.deleteProduct()
	case "5":
		h.listOrders()
	case "6":
		h.viewOrderDetail()
	case "7":
		h.updateOrderStatus()
	case "8":
		h.deleteOrder()
	case "9":
		h.viewProfile()
	case "10":
		h.updateProfile()
	case "11":
		h.deleteMyAccount()
	case "12":
		h.handleRegister()
	case "13":
		h.currentAdmin = nil
		fmt.Println("\n[INFO] Berhasil logout.")
	case "0":
		fmt.Println("\nProgram dihentikan. Sampai jumpa!")
		os.Exit(0)
	default:
		fmt.Println("[Peringatan] Pilihan tidak valid.")
	}
}

// ======================== FITUR PRODUK MINUMAN ======================== //

func (h *ConsoleHandler) listProducts() {
	products, err := h.productUsecase.GetAllProducts()
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return
	}

	fmt.Println("\n======================== DAFTAR KATALOG MINUMAN ========================")
	fmt.Printf("%-5s | %-30s | %-15s | %-8s\n", "ID", "NAMA MINUMAN", "HARGA (RP)", "STOK")
	fmt.Println("------------------------------------------------------------------------")
	for _, p := range products {
		fmt.Printf("%-5d | %-30s | Rp %-12.2f | %-8d\n", p.ID, p.Name, p.Price, p.Stock)
	}
	fmt.Println("------------------------------------------------------------------------")
}

func (h *ConsoleHandler) createProduct() {
	fmt.Println("\n>>> TAMBAH MINUMAN BARU <<<")
	name := h.inputString("Nama Minuman : ")
	price, err := h.inputFloat("Harga Satuan (Rp) : ")
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return
	}
	stock, err := h.inputInt("Jumlah Stok Awal  : ")
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return
	}

	err = h.productUsecase.AddProduct(name, price, stock)
	if err != nil {
		fmt.Printf("[GAGAL] %v\n", err)
		return
	}
	fmt.Println("[SUKSES] Minuman baru berhasil ditambahkan!")
}

func (h *ConsoleHandler) updateProduct() {
	fmt.Println("\n>>> UBAH DATA PRODUK MINUMAN <<<")
	id, err := h.inputInt("ID Minuman yang diubah: ")
	if err != nil {
		fmt.Println("[ERROR]", err)
		return
	}

	existing, err := h.productUsecase.GetProduct(id)
	if err != nil {
		fmt.Printf("[GAGAL] %v\n", err)
		return
	}

	fmt.Printf("Data Lama -> Nama: %s | Harga: Rp %.2f | Stok: %d\n", existing.Name, existing.Price, existing.Stock)
	newName := h.inputString(fmt.Sprintf("Nama Baru [%s]: ", existing.Name))
	if newName == "" {
		newName = existing.Name
	}

	newPriceStr := h.inputString(fmt.Sprintf("Harga Baru [%.2f]: ", existing.Price))
	newPrice := existing.Price
	if newPriceStr != "" {
		if val, err := strconv.ParseFloat(newPriceStr, 64); err == nil {
			newPrice = val
		}
	}

	newStockStr := h.inputString(fmt.Sprintf("Stok Baru [%d]: ", existing.Stock))
	newStock := existing.Stock
	if newStockStr != "" {
		if val, err := strconv.Atoi(newStockStr); err == nil {
			newStock = val
		}
	}

	err = h.productUsecase.EditProduct(id, newName, newPrice, newStock)
	if err != nil {
		fmt.Printf("[GAGAL] %v\n", err)
		return
	}
	fmt.Println("[SUKSES] Data minuman berhasil diperbarui!")
}

func (h *ConsoleHandler) deleteProduct() {
	id, err := h.inputInt("ID Minuman yang akan dihapus: ")
	if err != nil {
		fmt.Println("[ERROR]", err)
		return
	}

	if strings.ToLower(h.inputString("Yakin ingin menghapus? [y/N]: ")) != "y" {
		fmt.Println("Dibatalkan.")
		return
	}

	err = h.productUsecase.RemoveProduct(id)
	if err != nil {
		fmt.Printf("[GAGAL] %v\n", err)
		return
	}
	fmt.Println("[SUKSES] Minuman berhasil dihapus.")
}

// ======================== FITUR PESANAN USER ======================== //

func (h *ConsoleHandler) listOrders() {
	orders, err := h.orderUsecase.GetAllOrders()
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return
	}

	fmt.Println("\n============================= DAFTAR PESANAN MASUK =============================")
	fmt.Printf("%-4s | %-18s | %-10s | %-14s | %-9s\n", "ID", "INVOICE NUMBER", "USER ID", "TOTAL (RP)", "STATUS")
	fmt.Println("--------------------------------------------------------------------------------")
	for _, o := range orders {
		fmt.Printf("%-4d | %-18s | User #%-5d | Rp %-11.2f | %-9s\n", o.ID, o.InvoiceNumber, o.UserID, o.TotalPrice, o.Status)
	}
	fmt.Println("--------------------------------------------------------------------------------")
}

func (h *ConsoleHandler) viewOrderDetail() {
	id, err := h.inputInt("Masukkan ID Pesanan: ")
	if err != nil {
		fmt.Println("[ERROR]", err)
		return
	}

	items, err := h.orderUsecase.GetOrderDetails(id)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return
	}

	if len(items) == 0 {
		fmt.Println("Tidak ada rincian item untuk pesanan ini.")
		return
	}

	fmt.Printf("\n--- Rincian Minuman untuk Order #%d ---\n", id)
	fmt.Printf("%-4s | %-12s | %-8s | %-12s | %-12s\n", "NO", "PRODUCT ID", "QTY", "HARGA SATUAN", "SUBTOTAL")
	fmt.Println("-------------------------------------------------------------------------")
	var grandTotal float64
	for i, it := range items {
		// PERHITUNGAN SUBTOTAL: Murni hasil perkalian quantity * price (tanpa kolom palsu di database!)
		subtotal := float64(it.Quantity) * it.Price
		grandTotal += subtotal
		fmt.Printf("%-4d | Produk #%-6d | %-8d | Rp %-9.2f | Rp %-9.2f\n", i+1, it.ProductID, it.Quantity, it.Price, subtotal)
	}
	fmt.Println("-------------------------------------------------------------------------")
	fmt.Printf("TOTAL PEMBELIAN: Rp %.2f\n", grandTotal)
}

func (h *ConsoleHandler) updateOrderStatus() {
	id, err := h.inputInt("Masukkan ID Pesanan: ")
	if err != nil {
		fmt.Println("[ERROR]", err)
		return
	}

	fmt.Println("Pilihan status baru:")
	fmt.Println("1. success (Pesanan selesai / dibayar)")
	fmt.Println("2. pending (Pesanan menunggu proses)")
	pilihan := h.inputString("Pilih [1/2]: ")

	var newStatus string
	if pilihan == "1" {
		newStatus = "success"
	} else if pilihan == "2" {
		newStatus = "pending"
	} else {
		fmt.Println("[GAGAL] Pilihan status tidak valid.")
		return
	}

	err = h.orderUsecase.ChangeStatus(id, newStatus)
	if err != nil {
		fmt.Printf("[GAGAL] %v\n", err)
		return
	}
	fmt.Printf("[SUKSES] Status pesanan ID #%d diubah ke '%s'!\n", id, newStatus)
}

func (h *ConsoleHandler) deleteOrder() {
	id, err := h.inputInt("Masukkan ID Pesanan: ")
	if err != nil {
		fmt.Println("[ERROR]", err)
		return
	}

	if strings.ToLower(h.inputString("Yakin ingin menghapus Order ID ini? [y/N]: ")) != "y" {
		fmt.Println("Dibatalkan.")
		return
	}

	err = h.orderUsecase.DeleteOrder(id)
	if err != nil {
		fmt.Printf("[GAGAL] %v\n", err)
		return
	}
	fmt.Printf("[SUKSES] Pesanan ID #%d berhasil dihapus via Transaksi ACID!\n", id)
}

// ======================== FITUR AKUN ADMIN ======================== //

func (h *ConsoleHandler) viewProfile() {
	profile, err := h.userUsecase.GetProfile(h.currentAdmin.ID)
	if err != nil {
		fmt.Println("[ERROR]", err)
		return
	}
	h.currentAdmin = profile
	fmt.Println("\n======================= PROFIL AKUN ADMIN =======================")
	fmt.Printf(" ID Akun        : %d\n", profile.ID)
	fmt.Printf(" Nama Lengkap   : %s\n", profile.Name)
	fmt.Printf(" Email Akun     : %s\n", profile.Email)
	fmt.Printf(" Peran (Role)   : %s\n", profile.Role)
	fmt.Printf(" Saldo Sistem   : Rp %.2f\n", profile.Balance)
	fmt.Printf(" Terdaftar Pada : %s\n", profile.CreatedAt.Format("02 January 2006 15:04:05"))
	fmt.Println("==================================================================")
}

func (h *ConsoleHandler) updateProfile() {
	newName := h.inputString(fmt.Sprintf("Nama Baru [%s]: ", h.currentAdmin.Name))
	if newName == "" {
		newName = h.currentAdmin.Name
	}
	newEmail := h.inputString(fmt.Sprintf("Email Baru [%s]: ", h.currentAdmin.Email))
	if newEmail == "" {
		newEmail = h.currentAdmin.Email
	}
	newPass := h.inputString("Password Baru (kosongkan jika tidak diubah): ")
	if newPass == "" {
		newPass = h.currentAdmin.Password
	}

	err := h.userUsecase.UpdateProfile(h.currentAdmin.ID, newName, newEmail, newPass)
	if err != nil {
		fmt.Printf("[GAGAL] %v\n", err)
		return
	}
	h.currentAdmin.Name = newName
	h.currentAdmin.Email = newEmail
	h.currentAdmin.Password = newPass
	fmt.Println("[SUKSES] Profil berhasil diperbarui!")
}

func (h *ConsoleHandler) deleteMyAccount() {
	if h.inputString("Ketik 'HAPUS' untuk konfirmasi hapus akun: ") != "HAPUS" {
		fmt.Println("Dibatalkan.")
		return
	}

	err := h.userUsecase.DeleteAccount(h.currentAdmin.ID)
	if err != nil {
		fmt.Printf("[GAGAL] %v\n", err)
		return
	}
	fmt.Println("[SUKSES] Akun admin Anda telah dihapus.")
	h.currentAdmin = nil
}
