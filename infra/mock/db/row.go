package mock

import "time"

type MockRow struct {
	Value  interface{}   // legado — usado por outros testes
	Values []interface{} // novo — array para Scan múltiplo
	Err    error
}

func (m *MockRow) Scan(dest ...interface{}) error {
	if m.Err != nil {
		return m.Err
	}

	// PRIORIDADE 1: comportamento original (1 valor só)
	if m.Values == nil {
		if len(dest) == 0 {
			return nil
		}
		d := dest[0]

		switch ptr := d.(type) {
		case *time.Time:
			if v, ok := m.Value.(time.Time); ok {
				*ptr = v
			}
		case *int:
			if v, ok := m.Value.(int); ok {
				*ptr = v
			}
		case *int64:
			if v, ok := m.Value.(int64); ok {
				*ptr = v
			}
		case *bool:
			if v, ok := m.Value.(bool); ok {
				*ptr = v
			}
		case *string:
			if v, ok := m.Value.(string); ok {
				*ptr = v
			}
		}

		return nil
	}

	// PRIORIDADE 2: modo tabela (Values)
	for i, d := range dest {
		if i >= len(m.Values) {
			break
		}

		switch ptr := d.(type) {
		case *int64:
			if v, ok := m.Values[i].(int64); ok {
				*ptr = v
			}
		case *int:
			if v, ok := m.Values[i].(int); ok {
				*ptr = v
			}
		case *float64:
			if v, ok := m.Values[i].(float64); ok {
				*ptr = v
			}
		case *string:
			if v, ok := m.Values[i].(string); ok {
				*ptr = v
			}
		case *bool:
			if v, ok := m.Values[i].(bool); ok {
				*ptr = v
			}
		case *time.Time:
			if v, ok := m.Values[i].(time.Time); ok {
				*ptr = v
			}
		}
	}

	return nil
}
