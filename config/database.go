package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() (*sql.DB, error) {
	// Format: username:password@tcp(host:port)/dbname?parseTime=true
	dsn := "root:1234@tcp(127.0.0.1:3306)/beverages?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database tidak merespons: %v", err)
	}

	return db, nil
}
