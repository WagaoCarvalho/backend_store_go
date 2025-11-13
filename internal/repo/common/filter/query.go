package repo

// QueryBuilder define o contrato mínimo que um builder deve implementar
// para receber filtros de maneira genérica.
type QueryBuilder interface {
	AddCondition(field string, value any)
	AddILIKECondition(field string, value string)
}

// ----------------------------
// Interface genérica de filtro
// ----------------------------

type FilterCondition interface {
	Apply(qb QueryBuilder)
}

// ----------------------------
// Filtros genéricos reutilizáveis
// ----------------------------

// TextFilter aplica um ILIKE genérico em um campo textual.
type TextFilter struct {
	Field string
	Value string
}

func (f TextFilter) Apply(qb QueryBuilder) {
	if f.Value != "" {
		qb.AddILIKECondition(f.Field, f.Value)
	}
}

// EqualFilter aplica uma igualdade genérica (campo = valor).
type EqualFilter[T any] struct {
	Field string
	Value *T
}

func (f EqualFilter[T]) Apply(qb QueryBuilder) {
	if f.Value != nil {
		qb.AddCondition(f.Field+" =", *f.Value)
	}
}

// RangeFilter aplica condições >= e <= genéricas para tipos comparáveis.
type RangeFilter[T any] struct {
	FieldMin string
	FieldMax string
	Min      *T
	Max      *T
}

func (f RangeFilter[T]) Apply(qb QueryBuilder) {
	if f.Min != nil {
		qb.AddCondition(f.FieldMin+" >=", *f.Min)
	}
	if f.Max != nil {
		qb.AddCondition(f.FieldMax+" <=", *f.Max)
	}
}
