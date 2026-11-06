package handler

import (
	"beverages-cli/cli"
	"beverages-cli/domain"
	"beverages-cli/entity"
	repoMysql "beverages-cli/repository/mysql"
	"beverages-cli/usecase"
	"database/sql"
)

func AdminLandingPage(db *sql.DB, admin *entity.User) {
	// Kirim db ke repository
	userRepo := repoMysql.NewUserRepository(db)
	productRepo := repoMysql.NewProductRepository(db)
	orderRepo := repoMysql.NewOrderRepository(db)

	// Kirim repository ke usecase
	userUsecase := usecase.NewUserUsecase(userRepo)
	productUsecase := usecase.NewProductUsecase(productRepo)
	orderUsecase := usecase.NewOrderUsecase(orderRepo)

	// Kirim usecase ke CLI ConsoleHandler
	consoleHandler := cli.NewConsoleHandler(userUsecase, productUsecase, orderUsecase)

	domainAdmin := &domain.User{
		ID:        admin.ID,
		Name:      admin.Name,
		Email:     admin.Email,
		Password:  admin.Password,
		Role:      admin.Role,
		Balance:   admin.Balance,
		CreatedAt: admin.CreatedAt,
	}

	consoleHandler.StartAdminSession(domainAdmin)
}
