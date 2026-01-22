package dto

import (
	"testing"
	"time"

	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

/*
|--------------------------------------------------------------------------
| ToModel
|--------------------------------------------------------------------------
*/

func TestClientFilterDTO_ToModel_AllFields(t *testing.T) {
	now := time.Now()
	status := true
	version := 1

	dto := ClientFilterDTO{
		Name:        "Cliente XPTO",
		Email:       "teste@teste.com",
		CPF:         "12345678901",    // 11
		CNPJ:        "12345678000199", // 14
		Status:      &status,
		Version:     &version,
		CreatedFrom: &now,
		CreatedTo:   &now,
		UpdatedFrom: &now,
		UpdatedTo:   &now,
		Limit:       10,
		Offset:      5,
		SortBy:      "NAME",
		SortOrder:   "DESC",
	}

	model, err := dto.ToModel()
	assert.NoError(t, err)
	assert.NotNil(t, model)
}

func TestClientFilterDTO_ToModel_Empty(t *testing.T) {
	dto := ClientFilterDTO{}
	model, err := dto.ToModel()

	assert.NoError(t, err)
	assert.NotNil(t, model)
}

/*
|--------------------------------------------------------------------------
| Validate – casos válidos
|--------------------------------------------------------------------------
*/

func TestClientFilterDTO_Validate_ValidCases(t *testing.T) {
	now := time.Now()

	tests := []ClientFilterDTO{
		{Name: "Cliente", Limit: 10},
		{Email: "teste@teste.com", Limit: 10},
		{CPF: "12345678901", Limit: 10},     // válido
		{CNPJ: "12345678000199", Limit: 10}, // válido
		{Status: utils.BoolPtr(true), Limit: 10},
		{Version: utils.IntPtr(1), Limit: 10},
		{CreatedFrom: &now, Limit: 10},
		{UpdatedFrom: &now, Limit: 10},
		{Name: "Cliente", Limit: 10, SortBy: "created_at"},
		{Name: "Cliente", Limit: 10, SortBy: "created_at", SortOrder: "asc"},
		{Name: "Cliente", Limit: 10, SortBy: "CREATED_AT", SortOrder: "DESC"},
	}

	for i, dto := range tests {
		t.Run(string(rune(i)), func(t *testing.T) {
			assert.NoError(t, dto.Validate())
		})
	}
}

/*
|--------------------------------------------------------------------------
| Validate – erros individuais
|--------------------------------------------------------------------------
*/

func TestClientFilterDTO_Validate_Errors(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)

	tests := []struct {
		name string
		dto  ClientFilterDTO
		err  string
	}{
		{
			"Sem filtros",
			ClientFilterDTO{Limit: 10},
			"pelo menos um filtro de busca deve ser fornecido",
		},
		{
			"Nome curto",
			ClientFilterDTO{Name: "ab", Limit: 10},
			"'name' deve conter no mínimo 3 caracteres",
		},
		{
			"Email curto",
			ClientFilterDTO{Email: "a@", Limit: 10},
			"'email' deve conter no mínimo 5 caracteres",
		},
		{
			"CPF inválido",
			ClientFilterDTO{CPF: "123", Limit: 10},
			"'cpf' inválido",
		},
		{
			"CNPJ inválido",
			ClientFilterDTO{CNPJ: "123", Limit: 10},
			"'cnpj' inválido",
		},
		{
			"Limit zero",
			ClientFilterDTO{Name: "Teste", Limit: 0},
			"'limit' deve ser maior que zero",
		},
		{
			"Limit > 100",
			ClientFilterDTO{Name: "Teste", Limit: 101},
			"'limit' não pode ser maior que 100",
		},
		{
			"Offset negativo",
			ClientFilterDTO{Name: "Teste", Limit: 10, Offset: -1},
			"'offset' não pode ser negativo",
		},
		{
			"Offset excede o limite permitido",
			ClientFilterDTO{Name: "Teste", Limit: 10, Offset: 10_001},
			"'offset' excede o limite permitido",
		},
		{
			"SortBy inválido",
			ClientFilterDTO{Name: "Teste", Limit: 10, SortBy: "foo"},
			"'sort_by' inválido",
		},
		{
			"SortOrder inválido",
			ClientFilterDTO{Name: "Teste", Limit: 10, SortBy: "name", SortOrder: "foo"},
			"'sort_order' inválido",
		},
		{
			"CreatedFrom > CreatedTo",
			ClientFilterDTO{Name: "Teste", Limit: 10, CreatedFrom: &now, CreatedTo: &past},
			"'created_from' não pode ser maior que 'created_to'",
		},
		{
			"UpdatedFrom > UpdatedTo",
			ClientFilterDTO{Name: "Teste", Limit: 10, UpdatedFrom: &now, UpdatedTo: &past},
			"'updated_from' não pode ser maior que 'updated_to'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dto.Validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.err)
		})
	}
}

/*
|--------------------------------------------------------------------------
| Validate – erros acumulados
|--------------------------------------------------------------------------
*/

func TestClientFilterDTO_Validate_MultipleErrors(t *testing.T) {
	dto := ClientFilterDTO{
		Limit:  0,
		Offset: -1,
	}

	err := dto.Validate()
	assert.Error(t, err)

	assert.Contains(t, err.Error(), "pelo menos um filtro de busca deve ser fornecido")
	assert.Contains(t, err.Error(), "'limit' deve ser maior que zero")
	assert.Contains(t, err.Error(), "'offset' não pode ser negativo")
}

/*
|--------------------------------------------------------------------------
| TrimSpace
|--------------------------------------------------------------------------
*/

func TestClientFilterDTO_Validate_TrimSpace(t *testing.T) {
	tests := []struct {
		name string
		dto  ClientFilterDTO
		pass bool
	}{
		{
			"Só espaços",
			ClientFilterDTO{Name: "   ", Limit: 10},
			false,
		},
		{
			"1 char após trim",
			ClientFilterDTO{Name: "  a  ", Limit: 10},
			false,
		},
		{
			"Válido após trim",
			ClientFilterDTO{Name: "  João Silva  ", Limit: 10},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dto.Validate()
			if tt.pass {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

/*
|--------------------------------------------------------------------------
| isValidSortField – 100% branch coverage
|--------------------------------------------------------------------------
*/

func TestIsValidSortField_AllBranches(t *testing.T) {
	valid := []string{
		"id", "ID",
		"name",
		"email",
		"status",
		"version",
		"created_at",
		"updated_at",
	}

	for _, f := range valid {
		assert.True(t, isValidSortField(f))
	}

	invalid := []string{
		"",
		" ",
		"foo",
		"created",
		"updated",
		"nome",
	}

	for _, f := range invalid {
		assert.False(t, isValidSortField(f))
	}
}
