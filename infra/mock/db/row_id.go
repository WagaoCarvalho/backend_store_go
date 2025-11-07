package mock

import "time"

type MockRowWithID struct {
	IDValue   int64
	TimeValue time.Time
	Err       error
}

func (m *MockRowWithID) Scan(dest ...interface{}) error {
	if m.Err != nil {
		return m.Err
	}

	for _, d := range dest {
		switch ptr := d.(type) {
		case *int64:
			*ptr = m.IDValue
		case *time.Time:
			*ptr = m.TimeValue
		}

	}
	return nil
}
