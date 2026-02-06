package main

import (
	"expense_tracker/infrastructure/db"
	"fmt"
)

func main() {
	fmt.Println("Mock server running ...")
	db.DB_Init()
}
