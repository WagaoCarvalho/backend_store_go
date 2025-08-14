package repositories

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrRelationNotFound       = errors.New("relação usuário-categoria não encontrada")
	ErrRelationExists         = errors.New("relação já existe")
	ErrCreateRelation         = errors.New("erro ao criar relação")
	ErrGetRelationsByUser     = errors.New("erro ao buscar relações por usuário")
	ErrGetRelationsByCategory = errors.New("erro ao buscar relações por categoria")
	ErrScanRelation           = errors.New("erro ao ler relação")
	ErrIterateRelations       = errors.New("erro após ler relações")
	ErrDeleteRelation         = errors.New("erro ao deletar relação")
	ErrDeleteAllUserRelations = errors.New("erro ao deletar todas as relações do usuário")
	ErrInvalidForeignKey      = errors.New("usuário ou categoria inválido")
	ErrCheckRelationExists    = errors.New("erro ao verificar se a relação entre usuário e categoria existe")
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
