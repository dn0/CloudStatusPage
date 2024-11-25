package db

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type FieldIsRequiredError struct {
	Field string
}

type ObjectNotFoundError struct {
	Object string
}

func (e *FieldIsRequiredError) Error() string {
	return "field is required: " + e.Field
}

func (e *ObjectNotFoundError) Error() string {
	return "not found: " + e.Object
}

func ErrorIsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation
}
