# Backend
Backend repository for the Personal Expense Tracker final project @ A2SV


## Database Setup

1. Make sure PostgreSQL is installed and running.
2. Create the database:

```bash
createdb -U postgres Personal_Expense_tracker_DB


psql -U postgres -d Personal_Expense_tracker_DB -f infrastructure/db/schema.sql
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
Server listens on `:8080` (or `PORT` env var).

**API documentation (Swagger):** [http://localhost:8080/api-docs](http://localhost:8080/api-docs) — Team 2 (Expenses & Categories) endpoints are grouped in subsections there.

---

## Team 2: Expenses & Categories API

Expense and category endpoints **require JWT authentication** (same as User profile). Obtain a token via `POST /auth/login`, then send:

- **`Authorization: Bearer <token>`**

All expense and category operations are scoped to the authenticated user (ownership checks). `/api-docs` and `/` remain public.

### Expenses
| Method | Path | Description |
|--------|------|-------------|
| POST | `/expenses` | Create expense (body: amount, expense_date, category_id?, note?, etc.) |
| GET | `/expenses` | List expenses (query: `from_date`, `to_date`, `category_id`) |
| GET | `/expenses/:id` | Get one expense |
| PUT | `/expenses/:id` | Update expense |
| DELETE | `/expenses/:id` | Delete expense |

### Categories
| Method | Path | Description |
|--------|------|-------------|
| POST | `/categories` | Create category (body: name, user_id optional) |
| GET | `/categories` | List global + user's categories |
| GET | `/categories/:id` | Get category by ID |
| PUT | `/categories/:id` | Update category (own only) |
| DELETE | `/categories/:id` | Delete category (own only) |

### Example (curl)
```bash
# 1. Login to get token
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"yourpassword"}' | jq -r '.token')

# 2. Create category (user-defined)
curl -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"Food"}'

# 3. Create expense
curl -X POST http://localhost:8080/expenses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"amount":50,"expense_date":"2026-02-12","category_id":"<category-uuid>","note":"Lunch"}'

# 4. List expenses (optional filters)
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/expenses?from_date=2026-02-01&to_date=2026-02-28"
```