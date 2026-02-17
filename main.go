package main

import (
	"log"
	"net/http"
	"os"

	delivery "expense_tracker/delivery/http"
	"expense_tracker/infrastructure/auth"
	"expense_tracker/infrastructure/db"
	infraRepo "expense_tracker/infrastructure/repository"
	"expense_tracker/infrastructure/repositoryPG"
	"expense_tracker/usecases"
)

func main() {
	log.Println("Starting Expense Tracker server...")

	if err := db.DB_Init(); err != nil {
		log.Fatal(err)
	}

	// Repositories
	userRepo := repositoryPG.NewUserRepoPG(db.DB)
	expenseRepo := infraRepo.NewExpenseRepoPG(db.DB)
	categoryRepo := infraRepo.NewCategoryRepoPG(db.DB)

	// Infrastructure
	hasher := auth.BcryptHasher{}
	jwtSvc := auth.NewJWTService(os.Getenv("JWT_SECRET"))

	// Use cases
	authUC := usecases.NewAuthUsecase(userRepo, hasher, jwtSvc)
	userUC := usecases.NewUserUsecase(userRepo)
	expenseUC := usecases.NewExpenseUseCase(expenseRepo)
	categoryUC := usecases.NewCategoryUseCase(categoryRepo)

	// Handlers
	authHandler := delivery.NewAuthHandler(authUC)
	userHandler := delivery.NewUserHandler(userUC, jwtSvc)
	expenseHandler := delivery.NewExpenseHandler(expenseUC)
	categoryHandler := delivery.NewCategoryHandler(categoryUC)

	// Router (expenses, categories, api-docs) with JWT auth for expense/category
	router := delivery.NewRouter(expenseHandler, categoryHandler, jwtSvc)
	team2Handler := router.Handler()

	// Routes
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/login", authHandler.Login)
	mux.HandleFunc("/user/profile", userHandler.GetProfile)
	mux.HandleFunc("/user/update", userHandler.UpdateProfile)
	mux.Handle("/expenses", team2Handler)
	mux.Handle("/expenses/", team2Handler)
	mux.Handle("/categories", team2Handler)
	mux.Handle("/categories/", team2Handler)
	mux.Handle("/api-docs", team2Handler)
	mux.Handle("/api-docs/", team2Handler)
	mux.Handle("/", team2Handler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
