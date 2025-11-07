package mock

type MockRowWithInt struct {
	IntValue int
	Err      error
}

func (m *MockRowWithInt) Scan(dest ...interface{}) error {
	if m.Err != nil {
		return m.Err
	}

	if len(dest) > 0 {
		if ptr, ok := dest[0].(*int); ok {
			*ptr = m.IntValue
		}
	}
	return nil
}
