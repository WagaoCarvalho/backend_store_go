package dto

import (
	"testing"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_contact_relation"
	"github.com/stretchr/testify/assert"
)

func TestToContactRelationModel(t *testing.T) {
	dto := UserContactRelationDTO{
		ContactID: 1,
		UserID:    2,
	}

	model := ToContactRelationModel(dto)

	assert.Equal(t, int64(1), model.ContactID)
	assert.Equal(t, int64(2), model.UserID)
	assert.WithinDuration(t, time.Now(), model.CreatedAt, time.Second, "CreatedAt deve ser próximo ao tempo atual")
}

func TestToContactRelationDTO(t *testing.T) {
	now := time.Now()
	model := &models.UserContactRelation{
		ContactID: 1,
		UserID:    2,
		CreatedAt: now,
	}

	dto := ToContactRelationDTO(model)

	assert.Equal(t, int64(1), dto.ContactID)
	assert.Equal(t, int64(2), dto.UserID)
	assert.Equal(t, now.Format(time.RFC3339), dto.CreatedAt)
}

func TestToContactRelationDTO_NilModel(t *testing.T) {
	var model *models.UserContactRelation
	dto := ToContactRelationDTO(model)

	assert.Equal(t, UserContactRelationDTO{}, dto)
}
