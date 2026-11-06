package handler

import (
	"beverages-cli/entity"
	"beverages-cli/helper"
	"database/sql"
	"fmt"
	"os"
	"time"
)

// Storage temporari untuk menyimpan keranjang user selama aplikasi berjalan
var tempCart []entity.CartItem

// 2. Halaman Landing Page (Lobby User)
func UserLandingPage(db *sql.DB, user *entity.User) {
	for {
		// Selalu sinkronkan saldo terbaru dari DB
		_ = db.QueryRow("SELECT balance FROM user WHERE id = ?", user.ID).Scan(&user.Balance)

		fmt.Println("\n==========================================")
		fmt.Println("          DASHBOARD USER BEVERAGES        ")
		fmt.Println("==========================================")
		fmt.Printf("User  : %s\n", user.Name)
		fmt.Printf("Saldo : Rp %.2f\n", user.Balance)
		if len(tempCart) > 0 {
			fmt.Printf("Cart  : %d jenis produk tersimpan di keranjang\n", len(tempCart))
		}
		fmt.Println("------------------------------------------")
		fmt.Println("1. Top Up Saldo")
		fmt.Println("2. Order Minuman")
		fmt.Println("3. History Transaksi")
		fmt.Println("4. Keluar (Exit)")
		fmt.Println("==========================================")

		choice := helper.ReadPositiveInt("Pilih menu (1-4): ")

		switch choice {
		case 1:
			TopUpPage(db, user)
		case 2:
			OrderMenuPage(db, user)
		case 3:
			HistoryPage(db, user)
		case 4:
			fmt.Println("\nTerima kasih telah menggunakan aplikasi Beverages!")
			os.Exit(0)
		default:
			fmt.Println("[!] Pilihan tidak tersedia, silakan pilih 1 sampai 4.")
		}
	}
}

// 3. Halaman Top Up Saldo
func TopUpPage(db *sql.DB, user *entity.User) {
	for {
		fmt.Println("\n==========================================")
		fmt.Println("             HALAMAN TOP UP SALDO         ")
		fmt.Println("==========================================")
		fmt.Printf("Saldo Saat Ini: Rp %.2f\n", user.Balance)

		amount := helper.ReadPositiveFloat("Mau top up saldo berapa? Rp ")

		fmt.Println("\n[...] Saldo sedang diproses, mohon tunggu...")
		time.Sleep(1 * time.Second)

		_, err := db.Exec("UPDATE user SET balance = balance + ? WHERE id = ?", amount, user.ID)
		if err != nil {
			fmt.Printf("[!] Gagal melakukan top up: %v\n", err)
			return
		}

		user.Balance += amount

		fmt.Println("------------------------------------------")
		fmt.Println("             TOP UP BERHASIL!             ")
		fmt.Printf("Saldo berhasil ditambahkan : +Rp %.2f\n", amount)
		fmt.Printf("Total Saldo Anda sekarang  :  Rp %.2f\n", user.Balance)
		fmt.Println("------------------------------------------")

		for {
			fmt.Println("\nPilihan:")
			fmt.Println("1. Masih ingin Top Up lagi")
			fmt.Println("2. Exit (Kembali ke Dashboard)")
			subChoice := helper.ReadPositiveInt("Pilih (1/2): ")

			if subChoice == 1 {
				break // Lanjut loop untuk top up lagi
			} else if subChoice == 2 {
				return // Kembali ke halaman landing page / dashboard
			} else {
				fmt.Println("[!] Pilihan tidak valid! Masukkan 1 atau 2.")
			}
		}
	}
}

// 4. Halaman Menu & Order
func OrderMenuPage(db *sql.DB, user *entity.User) {
	// Jika storage temporari masih kosong, langsung minta user pilih produk pertama
	if len(tempCart) == 0 {
		berhasilTambah := tambahProdukKeCart(db)
		if !berhasilTambah {
			return
		}
	} else {
		fmt.Println("\n[i] Menampilkan kembali produk yang tersimpan di keranjang Anda...")
	}

	// Loop tampilan Rincian Order Sementara
	for {
		var totalOrderPrice float64
		fmt.Println("\n-------------------------------------------------------")
		fmt.Println("            RINCIAN ORDER SEMENTARA (CART)             ")
		fmt.Println("-------------------------------------------------------")
		if len(tempCart) == 0 {
			fmt.Println("          (Keranjang saat ini kosong)                  ")
		} else {
			for i, item := range tempCart {
				fmt.Printf("%d. %-22s x%d @ Rp %.2f = Rp %.2f\n",
					i+1, item.ProductName, item.Quantity, item.Price, item.Subtotal)
				totalOrderPrice += item.Subtotal
			}
		}
		fmt.Println("-------------------------------------------------------")
		fmt.Printf("TOTAL HARGA SELURUH ORDER : Rp %.2f\n", totalOrderPrice)
		fmt.Println("-------------------------------------------------------")

		fmt.Println("Pilihan:")
		fmt.Println("1. Tambah produk lagi")
		fmt.Println("2. Lanjut proses transaksi")
		fmt.Println("3. Hapus produk")
		fmt.Println("4. Exit")
		action := helper.ReadPositiveInt("Pilih aksi (1-4): ")

		switch action {
		case 1:
			tambahProdukKeCart(db)
		case 2:
			if len(tempCart) == 0 {
				fmt.Println("[!] Keranjang masih kosong! Silakan tambah produk terlebih dahulu.")
				continue
			}
			// Jika return false (karena pilih Top Up saat saldo kurang), kembali ke Dashboard
			// dengan kondisi tempCart tetap tersimpan utuh
			orderLagi := processCheckout(db, user, totalOrderPrice)
			if orderLagi {
				tambahProdukKeCart(db)
				continue
			}
			return
		case 3:
			if len(tempCart) == 0 {
				fmt.Println("[!] Tidak ada produk yang bisa dihapus karena keranjang kosong!")
				continue
			}
			hapusProdukDariCart()
		case 4:
			fmt.Println("\nTerima kasih, keluar dari program...")
			os.Exit(0)
		default:
			fmt.Println("[!] Pilihan tidak valid! Masukkan angka 1 sampai 4.")
		}
	}
}

// Fungsi untuk menampilkan menu DB, mengecek stok, dan menambah produk ke storage temporari
func tambahProdukKeCart(db *sql.DB) bool {
	fmt.Println("\n=======================================================")
	fmt.Println("                  DAFTAR MENU MINUMAN                  ")
	fmt.Println("=======================================================")

	rows, err := db.Query("SELECT id, name, price, stock FROM product ORDER BY id ASC")
	if err != nil {
		fmt.Printf("[!] Gagal memuat menu: %v\n", err)
		return false
	}

	products := make(map[int]entity.Product)
	adaStokTersedia := false

	fmt.Printf("%-4s | %-22s | %-12s | %-5s\n", "ID", "Nama Produk", "Harga", "Stok")
	fmt.Println("-------------------------------------------------------")
	for rows.Next() {
		var p entity.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err == nil {
			products[p.ID] = p
			if p.Stock > 0 {
				adaStokTersedia = true
			}
			fmt.Printf("%-4d | %-22s | Rp %9.2f | %-5d\n", p.ID, p.Name, p.Price, p.Stock)
		}
	}
	rows.Close()
	fmt.Println("=======================================================")

	if !adaStokTersedia {
		fmt.Println("[!] Maaf, seluruh stok produk sedang kosong.")
		return false
	}

	// TAHAP 1: Cek stok saat pilih menu
	var selectedProduct entity.Product
	var qtyInCart int

	for {
		prodID := helper.ReadPositiveInt("Pilih ID menu yang ingin dipesan: ")
		p, exists := products[prodID]
		if !exists {
			fmt.Println("[!] ID menu tidak ditemukan! Silakan pilih ID yang tertera.")
			continue
		}

		// Hitung jumlah produk ini yang sudah ada di storage temporari
		qtyInCart = 0
		for _, item := range tempCart {
			if item.ProductID == p.ID {
				qtyInCart = item.Quantity
				break
			}
		}

		// Cek apakah stok di DB 0 atau sudah habis diambil ke keranjang
		if p.Stock <= 0 {
			fmt.Printf("[!] Stok '%s' kosong! Silakan pilih menu lain.\n", p.Name)
			continue
		}
		if p.Stock-qtyInCart <= 0 {
			fmt.Printf("[!] Seluruh stok '%s' (%d item) sudah berada di keranjang Anda! Silakan pilih menu lain.\n", p.Name, p.Stock)
			continue
		}

		selectedProduct = p
		break
	}

	// TAHAP 2: Cek stok saat input quantity
	sisaStokBisaDiambil := selectedProduct.Stock - qtyInCart
	var qty int
	for {
		qty = helper.ReadPositiveInt(fmt.Sprintf("Masukkan quantity untuk %s (Stok tersedia: %d): ", selectedProduct.Name, sisaStokBisaDiambil))
		if qty > sisaStokBisaDiambil {
			fmt.Printf("[!] Stok tidak mencukupi! Maksimal quantity yang bisa ditambah adalah %d.\n", sisaStokBisaDiambil)
			continue
		}
		break
	}

	// Masukkan atau update produk di storage temporari (tempCart)
	found := false
	for i := range tempCart {
		if tempCart[i].ProductID == selectedProduct.ID {
			tempCart[i].Quantity += qty
			tempCart[i].Subtotal = float64(tempCart[i].Quantity) * tempCart[i].Price
			found = true
			break
		}
	}
	if !found {
		tempCart = append(tempCart, entity.CartItem{
			ProductID:   selectedProduct.ID,
			ProductName: selectedProduct.Name,
			Price:       selectedProduct.Price,
			Quantity:    qty,
			Subtotal:    float64(qty) * selectedProduct.Price,
		})
	}

	fmt.Printf("[+] Berhasil menambahkan %d %s ke keranjang.\n", qty, selectedProduct.Name)
	return true
}

// Fungsi tampilan hapus produk dari keranjang sementara
func hapusProdukDariCart() {
	fmt.Println("\n=======================================================")
	fmt.Println("              HALAMAN HAPUS PRODUK ORDER               ")
	fmt.Println("=======================================================")
	fmt.Printf("%-4s | %-22s | %-5s | %-12s\n", "No", "Nama Produk", "Qty", "Subtotal")
	fmt.Println("-------------------------------------------------------")
	for i, item := range tempCart {
		fmt.Printf("%-4d | %-22s | %-5d | Rp %9.2f\n", i+1, item.ProductName, item.Quantity, item.Subtotal)
	}
	fmt.Println("=======================================================")

	// Pilih nomor urut produk yang ingin dihapus
	var idx int
	for {
		nomor := helper.ReadPositiveInt(fmt.Sprintf("Pilih nomor produk yang ingin dihapus (1-%d): ", len(tempCart)))
		if nomor > len(tempCart) {
			fmt.Printf("[!] Nomor tidak valid! Pilih antara 1 sampai %d.\n", len(tempCart))
			continue
		}
		idx = nomor - 1
		break
	}

	targetItem := tempCart[idx]

	// Jika quantity lebih dari 1 (misal 2 atau lebih), tanyakan mau hapus berapa
	if targetItem.Quantity > 1 {
		delQty := helper.ReadPositiveInt(fmt.Sprintf("Produk '%s' memiliki quantity %d. Mau hapus berapa? ", targetItem.ProductName, targetItem.Quantity))

		// Jika input >= quantity yang ada (misal qty 2 tapi input 3), hapus produk tersebut tanpa sisa
		if delQty >= targetItem.Quantity {
			tempCart = append(tempCart[:idx], tempCart[idx+1:]...)
			fmt.Printf("\n[+] Produk '%s' berhasil dihapus seluruhnya dari keranjang!\n", targetItem.ProductName)
		} else {
			tempCart[idx].Quantity -= delQty
			tempCart[idx].Subtotal = float64(tempCart[idx].Quantity) * tempCart[idx].Price
			fmt.Printf("\n[+] Produk '%s' berhasil dihapus sebanyak %d! (Sisa quantity: %d)\n", targetItem.ProductName, delQty, tempCart[idx].Quantity)
		}
	} else {
		// Jika quantity hanya 1, langsung hapus dari slice
		tempCart = append(tempCart[:idx], tempCart[idx+1:]...)
		fmt.Printf("\n[+] Produk '%s' berhasil dihapus dari keranjang!\n", targetItem.ProductName)
	}
}

// Fungsi untuk memproses pembayaran ke database
func processCheckout(db *sql.DB, user *entity.User, totalOrderPrice float64) bool {
	// Ambil saldo terbaru user dari DB
	_ = db.QueryRow("SELECT balance FROM user WHERE id = ?", user.ID).Scan(&user.Balance)

	// Cek apakah saldo kurang dari total price
	if user.Balance < totalOrderPrice {
		fmt.Println("\n==========================================")
		fmt.Println("            SALDO ANDA KURANG!            ")
		fmt.Println("==========================================")
		fmt.Printf("Total Order : Rp %.2f\n", totalOrderPrice)
		fmt.Printf("Saldo Anda  : Rp %.2f\n", user.Balance)
		fmt.Printf("Kekurangan  : Rp %.2f\n", totalOrderPrice-user.Balance)
		fmt.Println("[i] Tenang, produk pilihan Anda tetap tersimpan di keranjang.")
		fmt.Println("------------------------------------------")
		fmt.Println("1. Top Up Saldo")
		fmt.Println("2. Exit (Keluar Program)")

		for {
			choice := helper.ReadPositiveInt("Pilih (1/2): ")
			if choice == 1 {
				TopUpPage(db, user)
				// Setelah selesai Top Up dan pilih Exit di halaman Top Up,
				// kembali ke Dashboard dengan tempCart yang masih tersimpan
				return false
			} else if choice == 2 {
				fmt.Println("\nKeluar dari program...")
				os.Exit(0)
			} else {
				fmt.Println("[!] Pilihan tidak valid! Masukkan 1 atau 2.")
			}
		}
	}

	// Jika saldo cukup, mulai Database Transaction
	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("[!] Gagal memulai transaksi: %v\n", err)
		return false
	}

	// TAHAP 3: Validasi akhir stok di DB sebelum eksekusi pembayaran
	for _, item := range tempCart {
		var currentStock int
		err := tx.QueryRow("SELECT stock FROM product WHERE id = ? FOR UPDATE", item.ProductID).Scan(&currentStock)
		if err != nil || currentStock < item.Quantity {
			tx.Rollback()
			fmt.Printf("\n[!] Transaksi dibatalkan! Stok '%s' di gudang tidak mencukupi (Sisa stok: %d).\n", item.ProductName, currentStock)
			return false
		}
	}

	invoiceNum := fmt.Sprintf("INV-%s-%03d", time.Now().Format("20060102"), time.Now().Second()+user.ID*10)

	// 1. Insert ke tabel `order`
	res, err := tx.Exec(
		"INSERT INTO `order` (invoice_number, user_id, total_price, status) VALUES (?, ?, ?, 'success')",
		invoiceNum, user.ID, totalOrderPrice,
	)
	if err != nil {
		tx.Rollback()
		fmt.Printf("[!] Gagal menyimpan order: %v\n", err)
		return false
	}

	orderID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		fmt.Printf("[!] Gagal mendapatkan ID order: %v\n", err)
		return false
	}

	// 2. Insert ke tabel `order_detail` & kurangi stok produk
	for _, item := range tempCart {
		_, err := tx.Exec(
			"INSERT INTO order_detail (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)",
			orderID, item.ProductID, item.Quantity, item.Price,
		)
		if err != nil {
			tx.Rollback()
			fmt.Printf("[!] Gagal menyimpan detail order: %v\n", err)
			return false
		}

		_, err = tx.Exec(
			"UPDATE product SET stock = stock - ? WHERE id = ?",
			item.Quantity, item.ProductID,
		)
		if err != nil {
			tx.Rollback()
			fmt.Printf("[!] Gagal memperbarui stok produk: %v\n", err)
			return false
		}
	}

	// 3. Kurangi saldo user
	_, err = tx.Exec("UPDATE user SET balance = balance - ? WHERE id = ?", totalOrderPrice, user.ID)
	if err != nil {
		tx.Rollback()
		fmt.Printf("[!] Gagal memotong saldo: %v\n", err)
		return false
	}

	if err := tx.Commit(); err != nil {
		fmt.Printf("[!] Gagal menyelesaikan transaksi: %v\n", err)
		return false
	}

	saldoAwal := user.Balance
	user.Balance -= totalOrderPrice

	// Tampilkan Invoice Transaksi Berhasil
	fmt.Println("\n=======================================================")
	fmt.Println("                  TRANSAKSI BERHASIL!                  ")
	fmt.Println("                    INVOICE PESANAN                    ")
	fmt.Println("=======================================================")
	fmt.Printf("No. Invoice : %s\n", invoiceNum)
	fmt.Printf("Pelanggan   : %s\n", user.Name)
	fmt.Printf("Tanggal     : %s\n", time.Now().Format("02-01-2006 15:04:05"))
	fmt.Println("-------------------------------------------------------")
	for _, item := range tempCart {
		fmt.Printf("- %-22s (x%d @ Rp %.2f) = Rp %.2f\n",
			item.ProductName, item.Quantity, item.Price, item.Subtotal)
	}
	fmt.Println("-------------------------------------------------------")
	fmt.Printf("Saldo Awal       : Rp %.2f\n", saldoAwal)
	fmt.Printf("Saldo Berkurang  : -Rp %.2f\n", totalOrderPrice)
	fmt.Printf("Sisa Saldo Anda  : Rp %.2f\n", user.Balance)
	fmt.Println("=======================================================")

	// Kosongkan storage temporari karena pesanan sudah berhasil dibayar
	tempCart = nil

	// Pilihan setelah transaksi berhasil
	for {
		fmt.Println("\nPilihan:")
		fmt.Println("1. Order Lagi")
		fmt.Println("2. Exit (Keluar Program)")
		postChoice := helper.ReadPositiveInt("Pilih (1/2): ")

		if postChoice == 1 {
			return true // Kembali ke menu order dengan keranjang baru
		} else if postChoice == 2 {
			fmt.Println("\nTerima kasih telah berbelanja! Keluar dari program...")
			os.Exit(0)
		} else {
			fmt.Println("[!] Pilihan tidak valid! Masukkan 1 atau 2.")
		}
	}
}

// 5. Halaman Histori Transaksi (Beserta Detail Barang yang Dibeli)
func HistoryPage(db *sql.DB, user *entity.User) {
	fmt.Println("\n===================================================================")
	fmt.Println("                       HISTORI TRANSAKSI ANDA                      ")
	fmt.Println("===================================================================")

	// 1. Ambil daftar transaksi utama milik user dari tabel `order`
	queryOrder := `SELECT id, invoice_number, total_price, status, created_at 
	               FROM ` + "`order`" + ` 
	               WHERE user_id = ? 
	               ORDER BY created_at DESC`

	rows, err := db.Query(queryOrder, user.ID)
	if err != nil {
		fmt.Printf("[!] Gagal mengambil histori transaksi: %v\n", err)
		return
	}

	var orders []entity.Order
	for rows.Next() {
		var o entity.Order
		if err := rows.Scan(&o.ID, &o.InvoiceNumber, &o.TotalPrice, &o.Status, &o.CreatedAt); err == nil {
			orders = append(orders, o)
		}
	}
	rows.Close() // Tutup cursor pertama sebelum melakukan query detail

	if len(orders) == 0 {
		fmt.Println("Belum ada histori transaksi untuk akun ini.")
		fmt.Println("===================================================================")
	} else {
		// 2. Untuk setiap order, ambil detail barang dari tabel order_detail JOIN product
		queryDetail := `SELECT p.name, od.quantity, od.price 
		                FROM order_detail od
		                JOIN product p ON od.product_id = p.id
		                WHERE od.order_id = ?`

		for i, o := range orders {
			fmt.Printf("[%d] No. Invoice : %s\n", i+1, o.InvoiceNumber)
			fmt.Printf("    Waktu       : %s\n", o.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("    Status      : %s\n", o.Status)
			fmt.Println("    Detail Barang yang Dibeli:")

			detailRows, err := db.Query(queryDetail, o.ID)
			if err != nil {
				fmt.Printf("      [!] Gagal memuat detail barang: %v\n", err)
			} else {
				for detailRows.Next() {
					var productName string
					var qty int
					var price float64

					if err := detailRows.Scan(&productName, &qty, &price); err == nil {
						subtotal := float64(qty) * price
						fmt.Printf("      - %-22s (x%d @ Rp %9.2f) = Rp %10.2f\n",
							productName, qty, price, subtotal)
					}
				}
				detailRows.Close()
			}

			fmt.Printf("    Total Harga : Rp %.2f\n", o.TotalPrice)
			fmt.Println("-------------------------------------------------------------------")
		}
	}

	for {
		fmt.Println("\nPilihan:")
		fmt.Println("1. Exit (Close Program)")
		fmt.Println("2. Kembali ke Landing Page")
		choice := helper.ReadPositiveInt("Pilih (1/2): ")

		if choice == 1 {
			fmt.Println("\nTerima kasih! Menutup program...")
			os.Exit(0)
		} else if choice == 2 {
			return
		} else {
			fmt.Println("[!] Pilihan tidak valid! Masukkan 1 atau 2.")
		}
	}
}
