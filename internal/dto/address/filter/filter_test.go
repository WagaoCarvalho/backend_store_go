package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
|--------------------------------------------------------------------------
| ToModel
|--------------------------------------------------------------------------
*/

func TestAddressFilterDTO_ToModel_AllFields(t *testing.T) {
	now := time.Now()
	userID := int64(1)
	clientCpfID := int64(2)
	supplierID := int64(3)
	isActive := true

	dto := AddressFilterDTO{
		UserID:      &userID,
		ClientCpfID: &clientCpfID,
		SupplierID:  &supplierID,
		City:        "São Paulo",
		State:       "SP",
		PostalCode:  "01234567",
		IsActive:    &isActive,
		CreatedFrom: &now,
		CreatedTo:   &now,
		UpdatedFrom: &now,
		UpdatedTo:   &now,
		Limit:       10,
		Offset:      5,
		SortBy:      "CITY",
		SortOrder:   "DESC",
	}

	model, err := dto.ToModel()
	assert.NoError(t, err)
	assert.NotNil(t, model)
}

func TestAddressFilterDTO_ToModel_Empty(t *testing.T) {
	dto := AddressFilterDTO{}
	model, err := dto.ToModel()

	assert.NoError(t, err)
	assert.NotNil(t, model)
}

/*
|--------------------------------------------------------------------------
| Validate – casos válidos
|--------------------------------------------------------------------------
*/

func TestAddressFilterDTO_Validate_ValidCases(t *testing.T) {
	now := time.Now()
	userID := int64(1)
	clientCpfID := int64(2)
	supplierID := int64(3)
	isActive := true

	tests := []AddressFilterDTO{
		{UserID: &userID, Limit: 10},
		{ClientCpfID: &clientCpfID, Limit: 10},
		{SupplierID: &supplierID, Limit: 10},
		{City: "São Paulo", Limit: 10},
		{State: "SP", Limit: 10},
		{PostalCode: "01234567", Limit: 10},
		{PostalCode: "01234-567", Limit: 10},  // Com hífen
		{PostalCode: "01.234-567", Limit: 10}, // Com pontos e hífen
		{IsActive: &isActive, Limit: 10},
		{CreatedFrom: &now, Limit: 10},
		{UpdatedFrom: &now, Limit: 10},
		{City: "São Paulo", Limit: 10, SortBy: "created_at"},
		{City: "São Paulo", Limit: 10, SortBy: "created_at", SortOrder: "asc"},
		{City: "São Paulo", Limit: 10, SortBy: "CITY", SortOrder: "DESC"},
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

func TestAddressFilterDTO_Validate_Errors(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)
	userID := int64(1)

	tests := []struct {
		name string
		dto  AddressFilterDTO
		err  string
	}{
		{
			"Sem filtros",
			AddressFilterDTO{Limit: 10},
			"pelo menos um filtro de busca deve ser fornecido",
		},
		{
			"Cidade curta",
			AddressFilterDTO{City: "S", Limit: 10},
			"'city' deve conter no mínimo 2 caracteres",
		},
		{
			"Estado com 1 caractere",
			AddressFilterDTO{State: "S", Limit: 10, UserID: &userID},
			"'state' deve conter exatamente 2 caracteres (UF)",
		},
		{
			"Estado com 3 caracteres",
			AddressFilterDTO{State: "SPO", Limit: 10, UserID: &userID},
			"'state' deve conter exatamente 2 caracteres (UF)",
		},
		{
			"CEP inválido - muito curto",
			AddressFilterDTO{PostalCode: "123", Limit: 10, UserID: &userID},
			"'postal_code' inválido - deve conter 8 dígitos",
		},
		{
			"CEP inválido - muito longo",
			AddressFilterDTO{PostalCode: "123456789", Limit: 10, UserID: &userID},
			"'postal_code' inválido - deve conter 8 dígitos",
		},
		{
			"Limit zero",
			AddressFilterDTO{City: "São Paulo", Limit: 0},
			"'limit' deve ser maior que zero",
		},
		{
			"Limit > 100",
			AddressFilterDTO{City: "São Paulo", Limit: 101},
			"'limit' não pode ser maior que 100",
		},
		{
			"Offset negativo",
			AddressFilterDTO{City: "São Paulo", Limit: 10, Offset: -1},
			"'offset' não pode ser negativo",
		},
		{
			"Offset excede o limite permitido",
			AddressFilterDTO{City: "São Paulo", Limit: 10, Offset: 10_001},
			"'offset' excede o limite permitido",
		},
		{
			"SortBy inválido",
			AddressFilterDTO{City: "São Paulo", Limit: 10, SortBy: "foo"},
			"'sort_by' inválido",
		},
		{
			"SortOrder inválido",
			AddressFilterDTO{City: "São Paulo", Limit: 10, SortBy: "city", SortOrder: "foo"},
			"'sort_order' inválido",
		},
		{
			"CreatedFrom > CreatedTo",
			AddressFilterDTO{UserID: &userID, Limit: 10, CreatedFrom: &now, CreatedTo: &past},
			"'created_from' não pode ser maior que 'created_to'",
		},
		{
			"UpdatedFrom > UpdatedTo",
			AddressFilterDTO{UserID: &userID, Limit: 10, UpdatedFrom: &now, UpdatedTo: &past},
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

func TestAddressFilterDTO_Validate_MultipleErrors(t *testing.T) {
	dto := AddressFilterDTO{
		Limit:  0,
		Offset: -1,
		State:  "S", // Isso é considerado um filtro válido para hasContentFilter
	}

	err := dto.Validate()
	assert.Error(t, err)

	// Não espera mais a mensagem de "pelo menos um filtro"
	assert.Contains(t, err.Error(), "'state' deve conter exatamente 2 caracteres (UF)")
	assert.Contains(t, err.Error(), "'limit' deve ser maior que zero")
	assert.Contains(t, err.Error(), "'offset' não pode ser negativo")

	// Verifica que NÃO contém a mensagem de filtro
	assert.NotContains(t, err.Error(), "pelo menos um filtro de busca deve ser fornecido")
}

/*
|--------------------------------------------------------------------------
| TrimSpace
|--------------------------------------------------------------------------
*/

func TestAddressFilterDTO_Validate_TrimSpace(t *testing.T) {
	userID := int64(1)

	tests := []struct {
		name string
		dto  AddressFilterDTO
		pass bool
	}{
		{
			"Cidade só espaços - considerada como sem filtro",
			AddressFilterDTO{City: "   ", Limit: 10, UserID: &userID},
			true, // Deve passar, pois vira string vazia e não é validada
		},
		{
			"Cidade 1 char após trim",
			AddressFilterDTO{City: "  S  ", Limit: 10, UserID: &userID},
			false,
		},
		{
			"Cidade válida após trim",
			AddressFilterDTO{City: "  São Paulo  ", Limit: 10, UserID: &userID},
			true,
		},
		{
			"Estado com espaços",
			AddressFilterDTO{State: "  SP  ", Limit: 10, UserID: &userID},
			true,
		},
		{
			"CEP com espaços",
			AddressFilterDTO{PostalCode: "  01234-567  ", Limit: 10, UserID: &userID},
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
| isValidAddressSortField – 100% branch coverage
|--------------------------------------------------------------------------
*/

func TestIsValidAddressSortField_AllBranches(t *testing.T) {
	valid := []string{
		"id", "ID",
		"user_id",
		"client_cpf_id",
		"supplier_id",
		"city",
		"state",
		"postal_code",
		"is_active",
		"created_at",
		"updated_at",
	}

	for _, f := range valid {
		assert.True(t, isValidAddressSortField(f))
	}

	invalid := []string{
		"",
		" ",
		"foo",
		"created",
		"updated",
		"cidade",
		"uf",
		"cep",
	}

	for _, f := range invalid {
		assert.False(t, isValidAddressSortField(f))
	}
}

/*
|--------------------------------------------------------------------------
| Validação de relacionamento exclusivo (opcional)
|--------------------------------------------------------------------------
*/

func TestAddressFilterDTO_Validate_ExclusiveRelations(t *testing.T) {
	userID := int64(1)
	clientCpfID := int64(2)
	supplierID := int64(3)

	tests := []struct {
		name string
		dto  AddressFilterDTO
		pass bool
	}{
		{
			"Apenas UserID",
			AddressFilterDTO{UserID: &userID, Limit: 10},
			true,
		},
		{
			"Apenas ClientCpfID",
			AddressFilterDTO{ClientCpfID: &clientCpfID, Limit: 10},
			true,
		},
		{
			"Apenas SupplierID",
			AddressFilterDTO{SupplierID: &supplierID, Limit: 10},
			true,
		},
		{
			"Múltiplos IDs - User e Client",
			AddressFilterDTO{UserID: &userID, ClientCpfID: &clientCpfID, Limit: 10},
			true, // O DTO não valida exclusividade, isso deve ser feito na camada de serviço
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
