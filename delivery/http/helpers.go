package http

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

func parseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

func isValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

func isErrNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
