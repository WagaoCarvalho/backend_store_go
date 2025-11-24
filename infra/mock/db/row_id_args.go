package mock

import (
	"errors"
	"time"
)

type MockRowWithIDArgs struct {
	Values []interface{}
	Err    error
}

// No arquivo infra/mock/db/db.go, no método Scan do MockRowWithIDArgs, adicione:
func (m *MockRowWithIDArgs) Scan(dest ...any) error {
	if m.Err != nil {
		return m.Err
	}

	if len(dest) != len(m.Values) {
		return errors.New("Scan: quantidade de valores não corresponde")
	}

	for i, d := range dest {
		switch ptr := d.(type) {
		case *int64:
			*ptr = m.Values[i].(int64)
		case **int64: // ADICIONAR ESTE CASE
			if m.Values[i] == nil {
				*ptr = nil
			} else {
				val := m.Values[i].(int64) // Valor direto, não ponteiro
				*ptr = &val
			}
		case *uint:
			*ptr = m.Values[i].(uint)
		case *int:
			*ptr = m.Values[i].(int)
		case *float64:
			*ptr = m.Values[i].(float64)
		case *string:
			*ptr = m.Values[i].(string)
		case *bool:
			*ptr = m.Values[i].(bool)
		case *time.Time:
			*ptr = m.Values[i].(time.Time)
		default:
			return errors.New("Scan: tipo não suportado")
		}
	}

	return nil
}
