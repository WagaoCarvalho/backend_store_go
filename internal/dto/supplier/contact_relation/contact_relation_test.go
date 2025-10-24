package dto

import (
	"testing"
	"time"

	model "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/contact_relation"
	"github.com/stretchr/testify/assert"
)

func TestToContactSupplierRelationModel(t *testing.T) {
	input := ContactSupplierRelationDTO{
		ContactID:  1,
		SupplierID: 2,
	}

	modelObj := ToContactSupplierRelationModel(input)

	assert.Equal(t, input.ContactID, modelObj.ContactID)
	assert.Equal(t, input.SupplierID, modelObj.SupplierID)
	assert.WithinDuration(t, time.Now(), modelObj.CreatedAt, time.Second)
}

func TestToContactSupplierRelationDTO(t *testing.T) {
	now := time.Now()
	modelObj := &model.SupplierContactRelation{
		ContactID:  1,
		SupplierID: 2,
		CreatedAt:  now,
	}

	dtoObj := ToContactSupplierRelationDTO(modelObj)

	assert.Equal(t, modelObj.ContactID, dtoObj.ContactID)
	assert.Equal(t, modelObj.SupplierID, dtoObj.SupplierID)
	assert.Equal(t, now.Format(time.RFC3339), dtoObj.CreatedAt)
}

func TestToContactSupplierRelationDTO_Nil(t *testing.T) {
	dtoObj := ToContactSupplierRelationDTO(nil)
	assert.Equal(t, int64(0), dtoObj.ContactID)
	assert.Equal(t, int64(0), dtoObj.SupplierID)
	assert.Equal(t, "", dtoObj.CreatedAt)
}
