package err

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

func IsForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
}

func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func IsCheckViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23514"
}

func IsDuplicateKey(err error) bool {
	if IsUniqueViolation(err) {
		return true
	}
	return err != nil && strings.Contains(err.Error(), "duplicate key")
}

func NewForeignKeyViolation(constraint string) *pgconn.PgError {
	return &pgconn.PgError{
		Code:           "23503",
		ConstraintName: constraint,
		Message:        "violação de chave estrangeira",
	}
}

func NewUniqueViolation(constraint string) *pgconn.PgError {
	return &pgconn.PgError{
		Code:           "23505",
		ConstraintName: constraint,
		Message:        "violação de constraint única",
	}
}

func NewCheckViolation(constraint string) *pgconn.PgError {
	return &pgconn.PgError{
		Code:           "23514",
		ConstraintName: constraint,
		Message:        "violação de constraint check",
	}
}
