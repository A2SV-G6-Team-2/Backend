package main

import (
	httpdelivery "expense_tracker/delivery/http"
	"expense_tracker/infrastructure/db"
	infrarepo "expense_tracker/infrastructure/repository"
	"expense_tracker/usecases"
	"log"
	"net/http"
)

func main() {
	if err := db.DB_Init(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	debtRepo := infrarepo.NewDebtRepositoryPG(db.DB)
	debtUsecase := usecases.NewDebtUsecase(debtRepo)
	debtHandler := httpdelivery.NewDebtHandler(debtUsecase)

	mux := http.NewServeMux()
	httpdelivery.RegisterDebtRoutes(mux, debtHandler)

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
