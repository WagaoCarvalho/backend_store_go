package repositories

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrRelationExists                 = errors.New("relação já existe")
	ErrRelationNotFound               = errors.New("relação supplier-categoria não encontrada")
	ErrCreateRelation                 = errors.New("erro ao criar relação")
	ErrCheckRelation                  = errors.New("erro ao verificar existência da relação")
	ErrGetRelationsBySupplier         = errors.New("erro ao buscar relações do fornecedor")
	ErrGetRelationsByCategory         = errors.New("erro ao buscar relações da categoria")
	ErrScanRelationRow                = errors.New("erro ao ler relação")
	ErrDeleteRelation                 = errors.New("erro ao deletar relação")
	ErrDeleteAllRelationsBySupplier   = errors.New("erro ao deletar todas as relações do fornecedor")
	ErrSupplierCategoryRelationUpdate = errors.New("erro ao atualizar a relação de categoria do fornecedor")
	ErrVersionConflict                = errors.New("conflito de versão: o registro foi modificado por outro processo")
	ErrSupplierNotFound               = errors.New("erro ao buscar fornecedor")
	ErrGetSuppliers                   = errors.New("erro ao buscar fornecedores")
	ErrSupplierUpdate                 = errors.New("erro ao atualizar fornecedor")
	ErrInvalidForeignKey              = errors.New("erro chave estrangeira inválida")
)

// IsForeignKeyViolation verifica se o erro é violação de chave estrangeira.
func IsForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23503" // código padrão PostgreSQL para FK violation
	}
	return false
}

// IsUniqueViolation verifica se o erro é violação de restrição única (chave duplicada).
func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" // código padrão PostgreSQL para unique violation
	}
	return false
}

// IsDuplicateKey verifica se o erro é relacionado a chave duplicada (fallback para erros em texto)
func IsDuplicateKey(err error) bool {
	if IsUniqueViolation(err) {
		return true
	}

	// fallback para erros que contenham a frase "duplicate key"
	return err != nil && strings.Contains(err.Error(), "duplicate key")
}
