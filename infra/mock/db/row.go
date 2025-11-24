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

	// --- MODO SIMPLES (Value único - legado) ---
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
		case *uint:
			switch v := m.Value.(type) {
			case uint:
				*ptr = v
			case uint64:
				*ptr = uint(v)
			case int64:
				*ptr = uint(v)
			}
		case *uint64:
			switch v := m.Value.(type) {
			case uint64:
				*ptr = v
			case uint:
				*ptr = uint64(v)
			case int64:
				*ptr = uint64(v)
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

	// --- MODO TABELA (Values) ---
	for i, d := range dest {
		if i >= len(m.Values) {
			break
		}

		switch ptr := d.(type) {
		case *int64:
			switch v := m.Values[i].(type) {
			case int64:
				*ptr = v
			case int:
				*ptr = int64(v)
			case uint:
				*ptr = int64(v)
			case uint64:
				*ptr = int64(v)
			}

		case *int:
			switch v := m.Values[i].(type) {
			case int:
				*ptr = v
			case int64:
				*ptr = int(v)
			case uint:
				*ptr = int(v)
			case uint64:
				*ptr = int(v)
			}

		case *uint:
			switch v := m.Values[i].(type) {
			case uint:
				*ptr = v
			case uint64:
				*ptr = uint(v)
			case int:
				*ptr = uint(v)
			case int64:
				*ptr = uint(v)
			}

		case *uint64:
			switch v := m.Values[i].(type) {
			case uint64:
				*ptr = v
			case uint:
				*ptr = uint64(v)
			case int:
				*ptr = uint64(v)
			case int64:
				*ptr = uint64(v)
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
