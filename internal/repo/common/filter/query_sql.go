package repo

import (
	"fmt"
	"strings"
)

// SQLQueryBuilder é a implementação genérica do QueryBuilder para SQL.
type SQLQueryBuilder struct {
	tableName string
	columns   []string
	where     strings.Builder
	args      []any
	pos       int
	orderBy   string
}

// NewSQLQueryBuilder cria um builder genérico para qualquer tabela.
// Exemplo: NewSQLQueryBuilder("products", []string{"id", "name", "price"}, "created_at DESC")
func NewSQLQueryBuilder(tableName string, columns []string, orderBy string) *SQLQueryBuilder {
	qb := &SQLQueryBuilder{
		tableName: tableName,
		columns:   columns,
		orderBy:   orderBy,
		args:      []any{},
		pos:       1,
	}
	qb.where.WriteString(" WHERE 1=1")
	return qb
}

// AddCondition adiciona uma condição do tipo "campo operador valor" (ex: "price >= $1").
func (qb *SQLQueryBuilder) AddCondition(condition string, value any) {
	qb.where.WriteString(fmt.Sprintf(" AND %s $%d", condition, qb.pos))
	qb.args = append(qb.args, value)
	qb.pos++
}

// AddILIKECondition adiciona uma busca textual com ILIKE (ex: "name ILIKE %$1%").
func (qb *SQLQueryBuilder) AddILIKECondition(field string, value string) {
	qb.where.WriteString(fmt.Sprintf(" AND %s ILIKE '%%' || $%d || '%%'", field, qb.pos))
	qb.args = append(qb.args, value)
	qb.pos++
}

// Build finaliza a query aplicando LIMIT, OFFSET e ORDER BY.
func (qb *SQLQueryBuilder) Build(limit, offset int) (string, []any) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		%s
		ORDER BY %s
		LIMIT %d OFFSET %d
	`,
		strings.Join(qb.columns, ", "),
		qb.tableName,
		qb.where.String(),
		qb.orderBy,
		limit,
		offset,
	)
	return strings.TrimSpace(query), qb.args
}
