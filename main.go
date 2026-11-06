package main

import (
	"beverages-cli/config"
	"beverages-cli/handler"
	"fmt"
	"log"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("[FATAL] Koneksi database gagal: %v", err)
	}
	defer db.Close()

	fmt.Println("Berhasil terhubung ke database 'beverages'.")

	// Masuk ke menu autentikasi awal (Masuk / Daftar)
	handler.AuthMenuPage(db)
}
