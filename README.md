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

Replace **postgres** with your _PostgreSQL username_ if different.

## Environment Setup

```bash
cp .env.example .env
```

## How to run the application

```bash
go run main.go
```

Server listens on `:8080` (or `PORT` env var).

**API documentation (Swagger):** [http://localhost:8080/api-docs](http://localhost:8080/api-docs) — Team 2 (Expenses & Categories) endpoints are grouped in subsections there.

---

## Team 2: Expenses & Categories API

Expense and category endpoints use **mock user ID** for this week. Send the header:

- **`X-User-ID`**: UUID of the current user (e.g. `8f3b2c9e-6d1a-4a9f-bc21-4e9f0a2d7c33`)

All expense and category operations are scoped to this user (ownership checks).

### Expenses

| Method | Path            | Description                                                            |
| ------ | --------------- | ---------------------------------------------------------------------- |
| POST   | `/expenses`     | Create expense (body: amount, expense_date, category_id?, note?, etc.) |
| GET    | `/expenses`     | List expenses (query: `from_date`, `to_date`, `category_id`)           |
| GET    | `/expenses/:id` | Get one expense                                                        |
| PUT    | `/expenses/:id` | Update expense                                                         |
| DELETE | `/expenses/:id` | Delete expense                                                         |

### Categories

| Method | Path              | Description                                    |
| ------ | ----------------- | ---------------------------------------------- |
| POST   | `/categories`     | Create category (body: name, user_id optional) |
| GET    | `/categories`     | List global + user's categories                |
| GET    | `/categories/:id` | Get category by ID                             |
| PUT    | `/categories/:id` | Update category (own only)                     |
| DELETE | `/categories/:id` | Delete category (own only)                     |

### Example (curl)

```bash
# Create category (user-defined)
curl -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 8f3b2c9e-6d1a-4a9f-bc21-4e9f0a2d7c33" \
  -d '{"name":"Food"}'

# Create expense
curl -X POST http://localhost:8080/expenses \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 8f3b2c9e-6d1a-4a9f-bc21-4e9f0a2d7c33" \
  -d '{"amount":50,"expense_date":"2026-02-12","category_id":"<category-uuid>","note":"Lunch"}'

# List expenses (optional filters)
curl -H "X-User-ID: 8f3b2c9e-6d1a-4a9f-bc21-4e9f0a2d7c33" \
  "http://localhost:8080/expenses?from_date=2026-02-01&to_date=2026-02-28"
```
