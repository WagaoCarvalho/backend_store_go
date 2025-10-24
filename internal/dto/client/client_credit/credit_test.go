package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"

	client_credit "github.com/WagaoCarvalho/backend_store_go/internal/model/client/credit"
)

func TestToClientCreditModel(t *testing.T) {
	id := int64(1)
	dtoInput := ClientCreditDTO{
		ID:            &id,
		ClientID:      42,
		AllowCredit:   true,
		CreditLimit:   1000.00,
		CreditBalance: 500.00,
	}

	model := ToClientCreditModel(dtoInput)

	assert.Equal(t, int64(1), model.ID)
	assert.Equal(t, int64(42), model.ClientID)
	assert.True(t, model.AllowCredit)
	assert.Equal(t, 1000.00, model.CreditLimit)
	assert.Equal(t, 500.00, model.CreditBalance)
}

func TestToClientCreditDTO(t *testing.T) {
	modelInput := &client_credit.ClientCredit{
		ID:            1,
		ClientID:      42,
		AllowCredit:   true,
		CreditLimit:   2000.00,
		CreditBalance: 1500.00,
	}

	dtoOutput := ToClientCreditDTO(modelInput)

	assert.NotNil(t, dtoOutput.ID)
	assert.Equal(t, int64(1), *dtoOutput.ID)
	assert.Equal(t, int64(42), dtoOutput.ClientID)
	assert.True(t, dtoOutput.AllowCredit)
	assert.Equal(t, 2000.00, dtoOutput.CreditLimit)
	assert.Equal(t, 1500.00, dtoOutput.CreditBalance)
}

func TestToClientCreditDTO_NilInput(t *testing.T) {
	dtoOutput := ToClientCreditDTO(nil)
	assert.Equal(t, ClientCreditDTO{}, dtoOutput)
}
