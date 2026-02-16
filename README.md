# Backend
Backend repository for the Personal Expense Tracker final project @ A2SV


## Database Setup

1. Make sure PostgreSQL is installed and running.
2. Create the database:

```bash
createdb -U postgres Personal_Expense_tracker_DB
```

3. Install Goose:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

4. Run migrations (from repo root):

```bash
export DB_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"
goose -dir infrastructure/db/migrations postgres "$DB_URL" up
```

Replace **postgres** with your *PostgreSQL username* if different.



## Environment Setup
```bash
cp .env.example .env
```

## How to run the application
```bash
go run main.go
```
Migrations run automatically on startup using the same DB settings in `.env`.