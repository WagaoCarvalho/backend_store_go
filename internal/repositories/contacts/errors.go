package repositories

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrContactNotFound         = errors.New("contato não encontrado")
	ErrCreateContact           = errors.New("erro ao criar contato")
	ErrFetchContact            = errors.New("erro ao buscar contato")
	ErrFetchContactsByUser     = errors.New("erro ao buscar contatos por user_id")
	ErrFetchContactsByClient   = errors.New("erro ao buscar contatos por client_id")
	ErrFetchContactsBySupplier = errors.New("erro ao buscar contatos por supplier_id")
	ErrScanContact             = errors.New("erro ao escanear contato")
	ErrUpdateContact           = errors.New("erro ao atualizar contato")
	ErrDeleteContact           = errors.New("erro ao deletar contato")
	ErrInvalidForeignKey       = errors.New("referência inválida: usuário, cliente ou fornecedor não existe")
)

func IsForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
}
