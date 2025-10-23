package dto

import (
	"testing"
	"time"

	model "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client_contact_relation"
	"github.com/stretchr/testify/assert"
)

func TestToClientContactRelationModel(t *testing.T) {
	input := ClientContactRelationDTO{
		ContactID: 1,
		ClientID:  2,
	}

	modelObj := ToClientContactRelationModel(input)

	assert.Equal(t, input.ContactID, modelObj.ContactID)
	assert.Equal(t, input.ClientID, modelObj.ClientID)
	assert.WithinDuration(t, time.Now(), modelObj.CreatedAt, time.Second)
}

func TestToClientContactRelationDTO(t *testing.T) {
	now := time.Now()
	modelObj := &model.ClientContactRelation{
		ContactID: 1,
		ClientID:  2,
		CreatedAt: now,
	}

	dtoObj := ToClientContactRelationDTO(modelObj)

	assert.Equal(t, modelObj.ContactID, dtoObj.ContactID)
	assert.Equal(t, modelObj.ClientID, dtoObj.ClientID)
	assert.Equal(t, now.Format(time.RFC3339), dtoObj.CreatedAt)
}

func TestToClientContactRelationDTO_Nil(t *testing.T) {
	dtoObj := ToClientContactRelationDTO(nil)
	assert.Equal(t, int64(0), dtoObj.ContactID)
	assert.Equal(t, int64(0), dtoObj.ClientID)
	assert.Equal(t, "", dtoObj.CreatedAt)
}
