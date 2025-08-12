package repositories

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrCreateProduct       = errors.New("erro ao criar produto")
	ErrGetProduct          = errors.New("erro ao buscar produto")
	ErrGetProducts         = errors.New("erro ao buscar produtos")
	ErrUpdateProduct       = errors.New("erro ao atualizar produto")
	ErrDeleteProduct       = errors.New("erro ao excluir produto")
	ErrInvalidForeignKey   = errors.New("chave estrangeira inválida")
	ErrProductNotFound     = errors.New("produto não encontrado")
	ErrFetchProductVersion = errors.New("erro ao buscar versão do produto")
	ErrDisableProduct      = errors.New("erro ao desabilitar produto")
	ErrEnableProduct       = errors.New("erro ao ativar produto")
	ErrVersionConflict     = errors.New("conflito de versão")
)

func IsForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
}
