package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() *sql.DB{
	dns := os.Getenv("MYSQL_DSN_CLI")

	db, err := sql.Open("mysql", dns)
	if err != nil {
		fmt.Printf("gagal membuka driver database: %v", err.Error())
		panic(err)
	}
	
	err = db.Ping()
	if err != nil {
		fmt.Println("gagal terhubung ke database")
		panic(err)
	}

	fmt.Println("Berhasil terhubung ke database");
	return db
}