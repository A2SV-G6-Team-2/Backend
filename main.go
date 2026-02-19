package main

import (
	"log"
	"net/http"
	"os"

	httpdelivery "expense_tracker/delivery/http"
	"expense_tracker/infrastructure/auth"
	"expense_tracker/infrastructure/db"
	infrarepo "expense_tracker/infrastructure/repository"
	"expense_tracker/infrastructure/repositoryPG"
	"expense_tracker/usecases"
)

func main() {
	log.Println("Starting Expense Tracker server...")

	if err := db.DB_Init(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	userRepo := repositoryPG.NewUserRepoPG(db.DB)
	expenseRepo := repositoryPG.NewExpenseRepoPG(db.DB)
	debtReportRepo := repositoryPG.NewDebtRepoPG(db.DB)
	debtRepo := infrarepo.NewDebtRepositoryPG(db.DB)

	hasher := auth.BcryptHasher{}
	jwtSvc := auth.NewJWTService(os.Getenv("JWT_SECRET"))

	authUC := usecases.NewAuthUsecase(userRepo, hasher, jwtSvc)
	userUC := usecases.NewUserUsecase(userRepo)
	reportUC := usecases.NewReportUsecase(expenseRepo, debtReportRepo)
	debtUsecase := usecases.NewDebtUsecase(debtRepo)

	authHandler := httpdelivery.NewAuthHandler(authUC)
	userHandler := httpdelivery.NewUserHandler(userUC, jwtSvc)
	reportHandler := httpdelivery.NewReportHandler(reportUC, jwtSvc)
	debtHandler := httpdelivery.NewDebtHandler(debtUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/login", authHandler.Login)
	mux.HandleFunc("/user/profile", userHandler.GetProfile)
	mux.HandleFunc("/user/update", userHandler.UpdateProfile)
	mux.HandleFunc("/reports/weekly", reportHandler.GetWeeklyReport)
	httpdelivery.RegisterDebtRoutes(mux, debtHandler)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
