package mock

import "time"

type MockRow struct {
	Value interface{}
	Err   error
}

// No arquivo mockdb, adicione no método Scan do MockRow:
func (m *MockRow) Scan(dest ...interface{}) error {
	if m.Err != nil {
		return m.Err
	}

	if len(dest) > 0 {
		switch ptr := dest[0].(type) {
		case *time.Time:
			if value, ok := m.Value.(time.Time); ok {
				*ptr = value
			}
		case *int:
			if value, ok := m.Value.(int); ok {
				*ptr = value
			}
		case *int64:
			if value, ok := m.Value.(int64); ok {
				*ptr = value
			}
		case *bool:
			if value, ok := m.Value.(bool); ok {
				*ptr = value
			}
		case *string:
			if value, ok := m.Value.(string); ok {
				*ptr = value
			}
		}
	}
	return nil
}
