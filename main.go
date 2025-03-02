package main

import (
	"api-otto/database"
	"api-otto/internal/handler"
	"api-otto/internal/repository"
	"api-otto/internal/service"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	// Setup database connection
	db, err := database.NewPostgresConnection("postgres://postgres:postgres@localhost:5432/voucher_redemption?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize repositories
	brandRepo := repository.NewBrandRepository(db)
	voucherRepo := repository.NewVoucherRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// Initialize services
	brandService := service.NewBrandService(brandRepo)
	voucherService := service.NewVoucherService(voucherRepo, brandRepo)
	transactionService := service.NewTransactionService(transactionRepo, voucherRepo)

	// Initialize handlers
	brandHandler := handler.NewBrandHandler(brandService)
	voucherHandler := handler.NewVoucherHandler(voucherService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// Setup router
	router := httprouter.New()

	// Brand routes
	router.POST("/brand", brandHandler.Create)
	router.GET("/brand/:id", brandHandler.GetByID)
	router.GET("/brand", brandHandler.GetAll)

	// Voucher routes
	router.POST("/voucher", voucherHandler.Create)
	router.GET("/brand/:id/vouchers", voucherHandler.GetByBrandID)
	router.GET("/voucher/:id", voucherHandler.GetByID)

	// Transaction routes
	router.POST("/transaction/redemption", transactionHandler.CreateRedemption)
	router.GET("/transaction/redemption/:id", transactionHandler.GetTransactionByID)

	// Start server
	log.Println("Server starting on :3000")
	log.Fatal(http.ListenAndServe(":3000", router))
} 