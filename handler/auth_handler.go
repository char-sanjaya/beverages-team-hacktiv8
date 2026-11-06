package handler

import (
	"beverages-cli/entity"
	"beverages-cli/helper"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Menu awal saat program pertama kali dijalankan
func AuthMenuPage(db *sql.DB) {
	for {
		fmt.Println("\n==========================================")
		fmt.Println("        SELAMAT DATANG DI BEVERAGES       ")
		fmt.Println("==========================================")
		fmt.Println("1. Masuk (Login)")
		fmt.Println("2. Daftar Akun Baru (Register)")
		fmt.Println("3. Keluar (Exit)")
		fmt.Println("==========================================")

		choice := helper.ReadPositiveInt("Pilih menu (1-3): ")

		switch choice {
		case 1:
			LoginPage(db)
		case 2:
			RegisterPage(db)
		case 3:
			fmt.Println("\nTerima kasih! Menutup program...")
			os.Exit(0)
		default:
			fmt.Println("[!] Pilihan tidak tersedia! Silakan pilih 1, 2, atau 3.")
		}
	}
}

// 1A. Halaman Login
func LoginPage(db *sql.DB) {
	for {
		fmt.Println("\n==========================================")
		fmt.Println("              HALAMAN LOGIN               ")
		fmt.Println("==========================================")

		email := helper.ReadString("Masukkan Email    : ")
		password := helper.ReadString("Masukkan Password : ")

		var user entity.User
		query := `SELECT id, name, email, password, role, balance, created_at 
		          FROM user 
		          WHERE email = ? AND password = ?`

		err := db.QueryRow(query, email, password).Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.Balance,
			&user.CreatedAt,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("\n[!] Email atau password salah!")
				fmt.Println("1. Coba Login Lagi")
				fmt.Println("2. Kembali ke Menu Awal")
				retry := helper.ReadPositiveInt("Pilih (1/2): ")
				if retry == 2 {
					return // Kembali ke AuthMenuPage
				}
				continue
			}
			fmt.Printf("\n[!] Terjadi kesalahan sistem: %v\n", err)
			return
		}

		// Cek role akun yang login
		if user.Role == "user" {
			fmt.Printf("\n[+] Login berhasil! Halo, %s.\n", user.Name)
			UserLandingPage(db, &user)
			return
		} else if user.Role == "admin" {
			fmt.Println("\n[i] Anda login sebagai Admin.")
			AdminLandingPage(db, &user)
			return
		}
	}
}

// 1B. Halaman Register (Daftar Pengguna Baru)
func RegisterPage(db *sql.DB) {
	fmt.Println("\n==========================================")
	fmt.Println("          HALAMAN DAFTAR AKUN BARU        ")
	fmt.Println("==========================================")

	// 1. Validasi Nama (Tidak boleh angka saja & maksimal 100 karakter sesuai VARCHAR(100))
	var name string
	for {
		name = helper.ReadString("Masukkan Nama Lengkap : ")
		if _, err := strconv.Atoi(name); err == nil {
			fmt.Println("[!] Input ditolak! Nama tidak boleh hanya berupa angka.")
			continue
		}
		if len(name) > 100 {
			fmt.Println("[!] Nama terlalu panjang (maksimal 100 karakter).")
			continue
		}
		break
	}

	// 2. Validasi Email (Format harus mengandung '@' dan '.', serta belum terdaftar di DB)
	var email string
	for {
		email = helper.ReadString("Masukkan Email        : ")
		if !strings.Contains(email, "@") || !strings.Contains(email, ".") || strings.Contains(email, " ") {
			fmt.Println("[!] Format email tidak valid! Contoh: budi@gmail.com")
			continue
		}
		if len(email) > 100 {
			fmt.Println("[!] Email terlalu panjang (maksimal 100 karakter).")
			continue
		}

		// Cek apakah email sudah ada di database (karena kolom email UNIQUE)
		var existingID int
		err := db.QueryRow("SELECT id FROM user WHERE email = ?", email).Scan(&existingID)
		if err == nil {
			fmt.Println("[!] Email sudah terdaftar! Silakan gunakan email lain.")
			continue
		} else if err != sql.ErrNoRows {
			fmt.Printf("[!] Gagal memeriksa email di database: %v\n", err)
			return
		}
		break
	}

	// 3. Validasi Password
	var password string
	for {
		password = helper.ReadString("Masukkan Password     : ")
		if len(password) < 6 {
			fmt.Println("[!] Password minimal harus 6 karakter!")
			continue
		}
		if len(password) > 100 {
			fmt.Println("[!] Password terlalu panjang (maksimal 100 karakter).")
			continue
		}
		break
	}

	// 4. Insert ke tabel user (role otomatis 'user', balance otomatis 0.00)
	query := `INSERT INTO user (name, email, password, role, balance) VALUES (?, ?, ?, 'user', 0.00)`
	_, err := db.Exec(query, name, email, password)
	if err != nil {
		fmt.Printf("\n[!] Gagal mendaftarkan akun: %v\n", err)
		return
	}

	fmt.Println("------------------------------------------")
	fmt.Println("[+] Pendaftaran berhasil! Silakan masuk menggunakan akun baru Anda.")
	fmt.Println("------------------------------------------")
}
