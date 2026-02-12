package main

import (
	deliveryhttp "expense_tracker/delivery/http"
	"expense_tracker/infrastructure/db"
	"expense_tracker/infrastructure/repository"
	"expense_tracker/usecases"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := db.DB_Init(); err != nil {
		log.Fatalf("db init: %v", err)
	}

	expenseRepo := repository.NewExpenseRepoPG(db.DB)
	categoryRepo := repository.NewCategoryRepoPG(db.DB)

	expenseUC := usecases.NewExpenseUseCase(expenseRepo)
	categoryUC := usecases.NewCategoryUseCase(categoryRepo)

	expenseHandler := deliveryhttp.NewExpenseHandler(expenseUC)
	categoryHandler := deliveryhttp.NewCategoryHandler(categoryUC)

	router := deliveryhttp.NewRouter(expenseHandler, categoryHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, router.Handler()); err != nil {
		log.Fatalf("server: %v", err)
	}
}
