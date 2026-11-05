package main

import (
	"beverages-team-hacktiv8/cli"
	"beverages-team-hacktiv8/config"
	repoMysql "beverages-team-hacktiv8/repository/mysql"
	"beverages-team-hacktiv8/usecase"
)

func main(){
	// inisialisasi terhubung ke database
	db := config.InitDB()
	defer db.Close()

	// kirim db ke repository
	userRepo := repoMysql.NewUserRepository(db)
	productRepo := repoMysql.NewProductRepository(db)
	orderRepo := repoMysql.NewOrderRepository(db)

	// Kirim Repository ke Usecase
	userUsecase := usecase.NewUserUsecase(userRepo)
	productUsecase := usecase.NewProductUsecase(productRepo)
	orderUsecase := usecase.NewOrderUsecase(orderRepo)

	// Kirim Usecase ke CLI
	consoleHandler := cli.NewConsoleHandler(userUsecase, productUsecase, orderUsecase)

	// running CLI
	consoleHandler.Start()

}