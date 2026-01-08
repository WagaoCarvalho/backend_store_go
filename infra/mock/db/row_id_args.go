package mock

import (
	"errors"
	"time"
)

type MockRowWithIDArgs struct {
	Values []interface{}
	Err    error
}

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
			// Para *int64, o valor pode ser int64 ou *int64
			if m.Values[i] == nil {
				*ptr = 0
			} else {
				switch v := m.Values[i].(type) {
				case int64:
					*ptr = v
				case *int64:
					if v == nil {
						*ptr = 0
					} else {
						*ptr = *v
					}
				default:
					return errors.New("Scan: tipo inválido para *int64")
				}
			}
		case **int64:
			// Para **int64, o valor pode ser nil, int64 ou *int64
			if m.Values[i] == nil {
				*ptr = nil
			} else {
				switch v := m.Values[i].(type) {
				case int64:
					val := v
					*ptr = &val
				case *int64:
					*ptr = v
				default:
					return errors.New("Scan: tipo inválido para **int64")
				}
			}
		case *uint:
			if m.Values[i] == nil {
				*ptr = 0
			} else {
				*ptr = m.Values[i].(uint)
			}
		case *int:
			if m.Values[i] == nil {
				*ptr = 0
			} else {
				*ptr = m.Values[i].(int)
			}
		case *float64:
			if m.Values[i] == nil {
				*ptr = 0
			} else {
				*ptr = m.Values[i].(float64)
			}
		case *string:
			if m.Values[i] == nil {
				*ptr = ""
			} else {
				*ptr = m.Values[i].(string)
			}
		case *bool:
			if m.Values[i] == nil {
				*ptr = false
			} else {
				*ptr = m.Values[i].(bool)
			}
		case *time.Time:
			if m.Values[i] == nil {
				*ptr = time.Time{}
			} else {
				*ptr = m.Values[i].(time.Time)
			}
		default:
			return errors.New("Scan: tipo não suportado")
		}
	}

	return nil
}
