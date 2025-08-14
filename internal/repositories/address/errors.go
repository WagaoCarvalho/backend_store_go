package repositories

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrCreateAddress     = errors.New("erro ao criar endereço")
	ErrFetchAddress      = errors.New("erro ao buscar endereço")
	ErrAddressNotFound   = errors.New("endereço não encontrado")
	ErrUpdateAddress     = errors.New("erro ao atualizar endereço")
	ErrDeleteAddress     = errors.New("erro ao excluir endereço")
	ErrInvalidForeignKey = errors.New("referência inválida: usuário, cliente ou fornecedor não existe")
)

func IsForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
}
