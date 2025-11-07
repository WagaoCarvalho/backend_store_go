package mock

import "time"

type MockRow struct {
	Value interface{}
	Err   error
}

func (m *MockRow) Scan(dest ...interface{}) error {
	if m.Err != nil {
		return m.Err
	}

	if len(dest) > 0 {
		if ptr, ok := dest[0].(*time.Time); ok {
			if value, ok := m.Value.(time.Time); ok {
				*ptr = value
			}
		}
	}
	return nil
}
