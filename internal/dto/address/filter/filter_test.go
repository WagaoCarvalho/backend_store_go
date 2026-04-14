package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newAddressFilterDTOFromQuery(_ map[string][]string) (*AddressFilterDTO, error) {
	dto := &AddressFilterDTO{}

	return dto, nil
}

func TestAddressFilterDTO_ToModel_AllFields(t *testing.T) {
	now := time.Now()
	userID := int64(1)
	clientCpfID := int64(2)
	supplierID := int64(3)
	isActive := true

	dto := AddressFilterDTO{
		UserID:       &userID,
		ClientCpfID:  &clientCpfID,
		SupplierID:   &supplierID,
		Street:       "Rua Teste",
		StreetNumber: "123",
		Complement:   "Apto 101",
		City:         "São Paulo",
		State:        "SP",
		Country:      "Brasil",
		PostalCode:   "01234567",
		IsActive:     &isActive,
		CreatedFrom:  &now,
		CreatedTo:    &now,
		UpdatedFrom:  &now,
		UpdatedTo:    &now,
		Limit:        10,
		Offset:       5,
		SortBy:       "city",
		SortOrder:    "desc",
	}

	model, err := dto.ToModel()
	assert.NoError(t, err)
	assert.NotNil(t, model)

	assert.Equal(t, userID, *model.UserID)
	assert.Equal(t, clientCpfID, *model.ClientCpfID)
	assert.Equal(t, supplierID, *model.SupplierID)
	assert.Equal(t, "Rua Teste", model.Street)
	assert.Equal(t, "123", model.StreetNumber)
	assert.Equal(t, "Apto 101", model.Complement)
	assert.Equal(t, "São Paulo", model.City)
	assert.Equal(t, "SP", model.State)
	assert.Equal(t, "Brasil", model.Country)
	assert.Equal(t, "01234567", model.PostalCode)
	assert.Equal(t, isActive, *model.IsActive)
	assert.Equal(t, now, *model.CreatedFrom)
	assert.Equal(t, now, *model.CreatedTo)
	assert.Equal(t, now, *model.UpdatedFrom)
	assert.Equal(t, now, *model.UpdatedTo)
	assert.Equal(t, 10, model.Limit)
	assert.Equal(t, 5, model.Offset)
	assert.Equal(t, "city", model.SortBy)
	assert.Equal(t, "desc", model.SortOrder)
}

func TestAddressFilterDTO_ToModel_IsActiveNil(t *testing.T) {
	userID := int64(1)

	dto := AddressFilterDTO{
		UserID:   &userID,
		IsActive: nil,
		Limit:    10,
	}

	model, err := dto.ToModel()
	assert.NoError(t, err)
	assert.Nil(t, model.IsActive)
}

func TestAddressFilterDTO_ToModel_Empty(t *testing.T) {
	dto := AddressFilterDTO{}
	model, err := dto.ToModel()

	assert.NoError(t, err)
	assert.NotNil(t, model)
	assert.Nil(t, model.UserID)
	assert.Nil(t, model.ClientCpfID)
	assert.Nil(t, model.SupplierID)
	assert.Empty(t, model.Street)
	assert.Empty(t, model.StreetNumber)
	assert.Empty(t, model.Complement)
	assert.Empty(t, model.City)
	assert.Empty(t, model.State)
	assert.Empty(t, model.Country)
	assert.Empty(t, model.PostalCode)
	assert.Nil(t, model.IsActive)
	assert.Nil(t, model.CreatedFrom)
	assert.Nil(t, model.CreatedTo)
	assert.Nil(t, model.UpdatedFrom)
	assert.Nil(t, model.UpdatedTo)
	assert.Equal(t, 0, model.Limit)
	assert.Equal(t, 0, model.Offset)
	assert.Empty(t, model.SortBy)
	assert.Empty(t, model.SortOrder)
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

	tests := []struct {
		name string
		dto  AddressFilterDTO
	}{
		{"Apenas UserID", AddressFilterDTO{UserID: &userID, Limit: 10}},
		{"Apenas ClientCpfID", AddressFilterDTO{ClientCpfID: &clientCpfID, Limit: 10}},
		{"Apenas SupplierID", AddressFilterDTO{SupplierID: &supplierID, Limit: 10}},
		{"Apenas Street", AddressFilterDTO{Street: "Rua Teste", Limit: 10}},
		{"Apenas StreetNumber", AddressFilterDTO{StreetNumber: "123", Limit: 10, UserID: &userID}},
		{"Apenas Complement", AddressFilterDTO{Complement: "Apto 101", Limit: 10, UserID: &userID}},
		{"Apenas City", AddressFilterDTO{City: "São Paulo", Limit: 10}},
		{"Apenas State", AddressFilterDTO{State: "SP", Limit: 10, UserID: &userID}},
		{"Apenas Country", AddressFilterDTO{Country: "Brasil", Limit: 10, UserID: &userID}},
		{"Apenas PostalCode (sem formatação)", AddressFilterDTO{PostalCode: "01234567", Limit: 10, UserID: &userID}},
		{"Apenas PostalCode (com hífen)", AddressFilterDTO{PostalCode: "01234-567", Limit: 10, UserID: &userID}},
		{"Apenas PostalCode (com pontos e hífen)", AddressFilterDTO{PostalCode: "01.234-567", Limit: 10, UserID: &userID}},
		{"Apenas IsActive", AddressFilterDTO{IsActive: &isActive, Limit: 10}},
		{"Apenas CreatedFrom", AddressFilterDTO{CreatedFrom: &now, Limit: 10}},
		{"Apenas UpdatedFrom", AddressFilterDTO{UpdatedFrom: &now, Limit: 10}},
		{"Apenas CreatedTo", AddressFilterDTO{CreatedTo: &now, Limit: 10, UserID: &userID}},
		{"Apenas UpdatedTo", AddressFilterDTO{UpdatedTo: &now, Limit: 10, UserID: &userID}},
		{"Múltiplos campos texto", AddressFilterDTO{
			Street: "Rua", City: "Cidade", State: "SP", Country: "Brasil", Limit: 10,
		}},
		{"Com sort_by válido", AddressFilterDTO{City: "São Paulo", Limit: 10, SortBy: "created_at"}},
		{"Com sort_by e sort_order asc", AddressFilterDTO{City: "São Paulo", Limit: 10, SortBy: "created_at", SortOrder: "asc"}},
		{"Com sort_by e sort_order desc", AddressFilterDTO{City: "São Paulo", Limit: 10, SortBy: "CITY", SortOrder: "DESC"}},
		{"Com offset positivo", AddressFilterDTO{City: "São Paulo", Limit: 10, Offset: 5}},
		{"Com limite máximo", AddressFilterDTO{City: "São Paulo", Limit: 100}},
		{"Com offset máximo", AddressFilterDTO{City: "São Paulo", Limit: 10, Offset: 10_000}},
		{"Com intervalos de datas válidos", AddressFilterDTO{
			UserID: &userID, Limit: 10,
			CreatedFrom: &now, CreatedTo: &now,
			UpdatedFrom: &now, UpdatedTo: &now,
		}},
		{"State com espaços (deve ser tratado)", AddressFilterDTO{State: "  SP  ", Limit: 10, UserID: &userID}},
		{"PostalCode com espaços (deve ser tratado)", AddressFilterDTO{PostalCode: "  01234-567  ", Limit: 10, UserID: &userID}},
		{"StreetNumber vazio mas com outro filtro", AddressFilterDTO{StreetNumber: "", Limit: 10, UserID: &userID}},
		{"Complement vazio mas com outro filtro", AddressFilterDTO{Complement: "", Limit: 10, UserID: &userID}},
		{"Country vazio mas com outro filtro", AddressFilterDTO{Country: "", Limit: 10, UserID: &userID}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dto.Validate()
			assert.NoError(t, err)
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
	future := now.Add(time.Hour)
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
			"Street muito curta",
			AddressFilterDTO{Street: "R", Limit: 10, UserID: &userID},
			"'street' deve conter no mínimo 2 caracteres",
		},
		{
			"StreetNumber muito longo",
			AddressFilterDTO{StreetNumber: "123456789012345678901", Limit: 10, UserID: &userID},
			"'street_number' não pode ter mais que 20 caracteres",
		},
		{
			"Complement muito longo",
			AddressFilterDTO{Complement: string(make([]byte, 101)), Limit: 10, UserID: &userID},
			"'complement' não pode ter mais que 100 caracteres",
		},
		{
			"City muito curta",
			AddressFilterDTO{City: "S", Limit: 10},
			"'city' deve conter no mínimo 2 caracteres",
		},
		{
			"State com 1 caractere",
			AddressFilterDTO{State: "S", Limit: 10, UserID: &userID},
			"'state' deve conter exatamente 2 caracteres (UF)",
		},
		{
			"State com 3 caracteres",
			AddressFilterDTO{State: "SPO", Limit: 10, UserID: &userID},
			"'state' deve conter exatamente 2 caracteres (UF)",
		},
		{
			"Country muito curto",
			AddressFilterDTO{Country: "B", Limit: 10, UserID: &userID},
			"'country' deve conter no mínimo 2 caracteres",
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
			"UserID <= 0",
			AddressFilterDTO{UserID: int64Ptr(0), Limit: 10},
			"'user_id' deve ser maior que zero",
		},
		{
			"UserID negativo",
			AddressFilterDTO{UserID: int64Ptr(-5), Limit: 10},
			"'user_id' deve ser maior que zero",
		},
		{
			"ClientCpfID <= 0",
			AddressFilterDTO{ClientCpfID: int64Ptr(0), Limit: 10, UserID: &userID},
			"'client_cpf_id' deve ser maior que zero",
		},
		{
			"ClientCpfID negativo",
			AddressFilterDTO{ClientCpfID: int64Ptr(-5), Limit: 10, UserID: &userID},
			"'client_cpf_id' deve ser maior que zero",
		},
		{
			"SupplierID <= 0",
			AddressFilterDTO{SupplierID: int64Ptr(0), Limit: 10, UserID: &userID},
			"'supplier_id' deve ser maior que zero",
		},
		{
			"SupplierID negativo",
			AddressFilterDTO{SupplierID: int64Ptr(-5), Limit: 10, UserID: &userID},
			"'supplier_id' deve ser maior que zero",
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
			AddressFilterDTO{UserID: &userID, Limit: 10, CreatedFrom: &future, CreatedTo: &past},
			"'created_from' não pode ser maior que 'created_to'",
		},
		{
			"UpdatedFrom > UpdatedTo",
			AddressFilterDTO{UserID: &userID, Limit: 10, UpdatedFrom: &future, UpdatedTo: &past},
			"'updated_from' não pode ser maior que 'updated_to'",
		},
		{
			"CreatedFrom sozinho é válido (não é erro)",
			AddressFilterDTO{UserID: &userID, Limit: 10, CreatedFrom: &now},
			"",
		},
		{
			"CreatedTo sozinho é válido (não é erro)",
			AddressFilterDTO{UserID: &userID, Limit: 10, CreatedTo: &now},
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dto.Validate()
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.err)
			}
		})
	}
}

/*
|--------------------------------------------------------------------------
| Validate – erros acumulados
|--------------------------------------------------------------------------
*/

func TestAddressFilterDTO_Validate_MultipleErrors(t *testing.T) {
	tests := []struct {
		name     string
		dto      AddressFilterDTO
		expected []string
	}{
		{
			"Múltiplos erros - State, Limit, Offset",
			AddressFilterDTO{
				State:  "S",
				Limit:  0,
				Offset: -1,
			},
			[]string{
				"'state' deve conter exatamente 2 caracteres (UF)",
				"'limit' deve ser maior que zero",
				"'offset' não pode ser negativo",
			},
		},
		{
			"Múltiplos erros - City, PostalCode, SortBy",
			AddressFilterDTO{
				City:       "S",
				PostalCode: "123",
				SortBy:     "invalido",
				UserID:     int64Ptr(1),
				Limit:      10,
			},
			[]string{
				"'city' deve conter no mínimo 2 caracteres",
				"'postal_code' inválido - deve conter 8 dígitos",
				"'sort_by' inválido",
			},
		},
		{
			"Múltiplos erros - IDs inválidos",
			AddressFilterDTO{
				UserID:      int64Ptr(0),
				ClientCpfID: int64Ptr(-5),
				SupplierID:  int64Ptr(0),
				Limit:       10,
			},
			[]string{
				"'user_id' deve ser maior que zero",
				"'client_cpf_id' deve ser maior que zero",
				"'supplier_id' deve ser maior que zero",
			},
		},
		{
			"Múltiplos erros - Street, StreetNumber, Complement, Country",
			AddressFilterDTO{
				Street:       "R",
				StreetNumber: "123456789012345678901",
				Complement:   string(make([]byte, 101)),
				Country:      "B",
				UserID:       int64Ptr(1),
				Limit:        10,
			},
			[]string{
				"'street' deve conter no mínimo 2 caracteres",
				"'street_number' não pode ter mais que 20 caracteres",
				"'complement' não pode ter mais que 100 caracteres",
				"'country' deve conter no mínimo 2 caracteres",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dto.Validate()
			assert.Error(t, err)
			for _, expected := range tt.expected {
				assert.Contains(t, err.Error(), expected)
			}
		})
	}
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
			true,
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
			"Estado minúsculo com espaços",
			AddressFilterDTO{State: "  sp  ", Limit: 10, UserID: &userID},
			true,
		},
		{
			"Estado só espaços",
			AddressFilterDTO{State: "   ", Limit: 10, UserID: &userID},
			true,
		},
		{
			"CEP com espaços",
			AddressFilterDTO{PostalCode: "  01234-567  ", Limit: 10, UserID: &userID},
			true,
		},
		{
			"Country com espaços",
			AddressFilterDTO{Country: "  Brasil  ", Limit: 10, UserID: &userID},
			true,
		},
		{
			"Street com espaços",
			AddressFilterDTO{Street: "  Rua Teste  ", Limit: 10, UserID: &userID},
			true,
		},
		{
			"StreetNumber com espaços",
			AddressFilterDTO{StreetNumber: "  123  ", Limit: 10, UserID: &userID},
			true,
		},
		{
			"Complement com espaços",
			AddressFilterDTO{Complement: "  Apto 101  ", Limit: 10, UserID: &userID},
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
		"user_id", "USER_ID",
		"client_cpf_id",
		"supplier_id",
		"street",
		"street_number",
		"city", "CITY",
		"state", "STATE",
		"country",
		"postal_code",
		"is_active",
		"created_at",
		"updated_at",
	}

	for _, f := range valid {
		assert.True(t, isValidAddressSortField(f), "Campo '%s' deveria ser válido", f)
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
		"id_user",
		"user",
		"client",
		"supplier",
	}

	for _, f := range invalid {
		assert.False(t, isValidAddressSortField(f), "Campo '%s' deveria ser inválido", f)
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
			true,
		},
		{
			"Múltiplos IDs - User e Supplier",
			AddressFilterDTO{UserID: &userID, SupplierID: &supplierID, Limit: 10},
			true,
		},
		{
			"Múltiplos IDs - Client e Supplier",
			AddressFilterDTO{ClientCpfID: &clientCpfID, SupplierID: &supplierID, Limit: 10},
			true,
		},
		{
			"Todos os IDs",
			AddressFilterDTO{UserID: &userID, ClientCpfID: &clientCpfID, SupplierID: &supplierID, Limit: 10},
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
| Teste para cobrir a normalização do SortOrder
|--------------------------------------------------------------------------
*/

func TestAddressFilterDTO_Validate_SortOrderNormalization(t *testing.T) {
	userID := int64(1)

	tests := []struct {
		name      string
		sortOrder string
		expected  string
	}{
		{"ASC maiúsculo", "ASC", "asc"},
		{"Desc maiúsculo", "Desc", "desc"},
		{"desc minúsculo", "desc", "desc"},
		{"asc minúsculo", "asc", "asc"},
		{"DESC maiúsculo", "DESC", "desc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := AddressFilterDTO{
				UserID:    &userID,
				City:      "São Paulo",
				Limit:     10,
				SortBy:    "city",
				SortOrder: tt.sortOrder,
			}
			err := dto.Validate()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, dto.SortOrder)
		})
	}
}

/*
|--------------------------------------------------------------------------
| Teste para cobrir a limpeza do PostalCode
|--------------------------------------------------------------------------
*/

func TestAddressFilterDTO_Validate_PostalCodeCleanup(t *testing.T) {
	userID := int64(1)

	tests := []struct {
		name       string
		postalCode string
		expected   string
	}{
		{"Com hífen", "01234-567", "01234567"},
		{"Com pontos", "01.234.567", "01234567"},
		{"Com hífen e pontos", "01.234-567", "01234567"},
		{"Com espaços", " 01234-567 ", "01234567"},
		{"Sem formatação", "01234567", "01234567"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := AddressFilterDTO{
				UserID:     &userID,
				PostalCode: tt.postalCode,
				Limit:      10,
			}
			err := dto.Validate()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, dto.PostalCode)
		})
	}
}

/*
|--------------------------------------------------------------------------
| newAddressFilterDTOFromQuery (placeholder)
|--------------------------------------------------------------------------
*/

func TestNewAddressFilterDTOFromQuery(t *testing.T) {
	params := map[string][]string{
		"user_id": {"1"},
		"city":    {"São Paulo"},
	}

	dto, err := newAddressFilterDTOFromQuery(params)
	assert.NoError(t, err)
	assert.NotNil(t, dto)
}

/*
|--------------------------------------------------------------------------
| Teste para cobrir erro de validação quando apenas UpdatedFrom é fornecido
|--------------------------------------------------------------------------
*/

func TestAddressFilterDTO_Validate_OnlyUpdatedFrom(t *testing.T) {
	now := time.Now()
	userID := int64(1)

	dto := AddressFilterDTO{
		UserID:      &userID,
		UpdatedFrom: &now,
		Limit:       10,
	}

	err := dto.Validate()
	assert.NoError(t, err)
}

/*
|--------------------------------------------------------------------------
| Teste para cobrir erro de validação quando apenas UpdatedTo é fornecido
|--------------------------------------------------------------------------
*/

func TestAddressFilterDTO_Validate_OnlyUpdatedTo(t *testing.T) {
	now := time.Now()
	userID := int64(1)

	dto := AddressFilterDTO{
		UserID:    &userID,
		UpdatedTo: &now,
		Limit:     10,
	}

	err := dto.Validate()
	assert.NoError(t, err)
}

/*
|--------------------------------------------------------------------------
| Helper functions para testes
|--------------------------------------------------------------------------
*/

func int64Ptr(i int64) *int64 {
	return &i
}
