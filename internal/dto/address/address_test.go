package dto_test

import (
	"testing"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/address"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	"github.com/stretchr/testify/assert"
)

func TestAddressDTO_ToModel(t *testing.T) {
	id := int64(10)
	userID := int64(1)
	clientID := int64(2)
	supplierID := int64(3)
	isActive := false

	tests := []struct {
		name string
		dto  dto.AddressDTO
		want models.Address
	}{
		{
			name: "All fields set",
			dto: dto.AddressDTO{
				ID:         &id,
				UserID:     &userID,
				ClientID:   &clientID,
				SupplierID: &supplierID,
				Street:     "Rua A",
				City:       "Cidade",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
				IsActive:   &isActive,
			},
			want: models.Address{
				ID:         10,
				UserID:     &userID,
				ClientID:   &clientID,
				SupplierID: &supplierID,
				Street:     "Rua A",
				City:       "Cidade",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
				IsActive:   false,
			},
		},
		{
			name: "Nil ID and IsActive",
			dto: dto.AddressDTO{
				Street:     "Rua B",
				City:       "Cidade",
				State:      "RJ",
				Country:    "Brasil",
				PostalCode: "87654321",
			},
			want: models.Address{
				ID:         0,
				Street:     "Rua B",
				City:       "Cidade",
				State:      "RJ",
				Country:    "Brasil",
				PostalCode: "87654321",
				IsActive:   true, // default
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dto.ToAddressModel(tt.dto)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.UserID, got.UserID)
			assert.Equal(t, tt.want.ClientID, got.ClientID)
			assert.Equal(t, tt.want.SupplierID, got.SupplierID)
			assert.Equal(t, tt.want.Street, got.Street)
			assert.Equal(t, tt.want.City, got.City)
			assert.Equal(t, tt.want.State, got.State)
			assert.Equal(t, tt.want.Country, got.Country)
			assert.Equal(t, tt.want.PostalCode, got.PostalCode)
			assert.Equal(t, tt.want.IsActive, got.IsActive)
		})
	}
}

func TestAddressDTO_FromModel(t *testing.T) {
	userID := int64(1)
	address := models.Address{
		ID:         5,
		UserID:     &userID,
		Street:     "Rua Teste",
		City:       "Cidade",
		State:      "SP",
		Country:    "Brasil",
		PostalCode: "12345000",
		IsActive:   true,
	}

	dtoResult := dto.ToAddressDTO(&address)

	assert.Equal(t, &address.ID, dtoResult.ID)
	assert.Equal(t, address.UserID, dtoResult.UserID)
	assert.Equal(t, address.Street, dtoResult.Street)
	assert.Equal(t, address.City, dtoResult.City)
	assert.Equal(t, address.State, dtoResult.State)
	assert.Equal(t, address.Country, dtoResult.Country)
	assert.Equal(t, address.PostalCode, dtoResult.PostalCode)
	assert.Equal(t, &address.IsActive, dtoResult.IsActive)
}

func TestToAddressDTOs(t *testing.T) {
	userID := int64(1)
	clientID := int64(2)

	address1 := &models.Address{
		ID:         1,
		UserID:     &userID,
		ClientID:   &clientID,
		Street:     "Rua 1",
		City:       "Cidade1",
		State:      "SP",
		Country:    "Brasil",
		PostalCode: "11111-111",
		IsActive:   true,
	}

	address2 := &models.Address{
		ID:         2,
		Street:     "Rua 2",
		City:       "Cidade2",
		State:      "RJ",
		Country:    "Brasil",
		PostalCode: "22222-222",
		IsActive:   false,
	}

	t.Run("Multiple valid addresses", func(t *testing.T) {
		input := []*models.Address{address1, address2}
		result := dto.ToAddressDTOs(input)

		assert.Len(t, result, 2)
		assert.Equal(t, address1.ID, *result[0].ID)
		assert.Equal(t, address2.ID, *result[1].ID)
		assert.Equal(t, address1.Street, result[0].Street)
		assert.Equal(t, address2.Street, result[1].Street)
	})

	t.Run("List with nil element", func(t *testing.T) {
		input := []*models.Address{address1, nil, address2}
		result := dto.ToAddressDTOs(input)

		assert.Len(t, result, 2) // ignora o nil
		assert.Equal(t, address1.ID, *result[0].ID)
		assert.Equal(t, address2.ID, *result[1].ID)
	})

	t.Run("Empty list", func(t *testing.T) {
		input := []*models.Address{}
		result := dto.ToAddressDTOs(input)

		assert.NotNil(t, result) // deve retornar slice vazio, não nil
		assert.Empty(t, result)
	})
}
