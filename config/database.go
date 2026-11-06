package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() (*sql.DB, error) {
	// Gunakan environment variable MYSQL_DSN_CLI jika tersedia, jika tidak gunakan default DSN
	dsn := os.Getenv("MYSQL_DSN_CLI")
	if dsn == "" {
		dsn = "root:1234@tcp(127.0.0.1:3306)/beverages?parseTime=true"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi database: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database tidak merespons: %v", err)
	}
	return db, nil
}

func InitDB() *sql.DB {
	db, err := ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	return db
}
