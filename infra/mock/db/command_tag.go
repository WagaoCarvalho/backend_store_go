// infra/mock/db/mock_commandtag.go
package mock

type MockCommandTag struct {
	RowsAffectedCount int64 // Renomeado para evitar conflito
}

func (m MockCommandTag) String() string {
	return "MOCK COMMAND TAG"
}

func (m MockCommandTag) RowsAffected() int64 {
	return m.RowsAffectedCount
}

func (m MockCommandTag) Insert() bool {
	return false
}

func (m MockCommandTag) Update() bool {
	return false
}

func (m MockCommandTag) Delete() bool {
	return false
}

func (m MockCommandTag) Select() bool {
	return false
}
