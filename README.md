# Backend
Backend repository for the Personal Expense Tracker final project @ A2SV


## Database Setup

1. Make sure PostgreSQL is installed and running.
2. Create the database:

```bash
createdb -U postgres Personal_Expense_tracker_DB


psql -U postgres -d Personal_Expense_tracker_DB -f infrastructure/db/schema.sql

Replace **postgres** with your *PostgreSQL username* if different.
