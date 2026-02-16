package main

import (
	"log"
	"net/http"
	"os"

	delivery "expense_tracker/delivery/http"
	"expense_tracker/infrastructure/auth"
	"expense_tracker/infrastructure/db"
	"expense_tracker/infrastructure/repositoryPG"

	"expense_tracker/usecases"
)

func main() {
	log.Println("Starting Expense Tracker server...")

	// Initialize DB
	if err := db.DB_Init(); err != nil {
		log.Fatal(err)
	}

	// Repositories
	userRepo := repositoryPG.NewUserRepoPG(db.DB)
	expenseRepo := repositoryPG.NewExpenseRepoPG(db.DB)
	debtRepo := repositoryPG.NewDebtRepoPG(db.DB)

	// Infrastructure services
	hasher := auth.BcryptHasher{}
	jwtSvc := auth.NewJWTService(os.Getenv("JWT_SECRET"))

	// Usecases
	authUC := usecases.NewAuthUsecase(userRepo, hasher, jwtSvc)
	userUC := usecases.NewUserUsecase(userRepo)
	reportUC := usecases.NewReportUsecase(expenseRepo, debtRepo)

	// Handlers
	authHandler := delivery.NewAuthHandler(authUC)
	userHandler := delivery.NewUserHandler(userUC, jwtSvc)
	reportHandler := delivery.NewReportHandler(reportUC, jwtSvc)

	// Routes
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/login", authHandler.Login)
	mux.HandleFunc("/user/profile", userHandler.GetProfile)
	mux.HandleFunc("/user/update", userHandler.UpdateProfile)
	mux.HandleFunc("/reports/weekly", reportHandler.GetWeeklyReport)

	// Start server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
