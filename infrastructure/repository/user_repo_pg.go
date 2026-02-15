package repository

import "database/sql"

type UserRepoPG struct {
	DB *sql.DB
}
