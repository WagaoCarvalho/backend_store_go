package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientCredit_Validate_Success(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      10,
		AllowCredit:   true,
		CreditLimit:   1000.0,
		CreditBalance: 500.0,
	}

	err := cc.Validate()
	assert.NoError(t, err)
}

func TestClientCredit_Validate_InvalidClientID(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      0,
		CreditLimit:   1000.0,
		CreditBalance: 500.0,
	}

	err := cc.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ClientID")
}

func TestClientCredit_Validate_NegativeLimit(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      1,
		CreditLimit:   -10.0,
		CreditBalance: 0,
	}

	err := cc.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CreditLimit")
}

func TestClientCredit_Validate_NegativeBalance(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      1,
		CreditLimit:   100.0,
		CreditBalance: -5.0,
	}

	err := cc.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CreditBalance")
}

func TestClientCredit_Validate_BalanceGreaterThanLimit(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      1,
		CreditLimit:   100.0,
		CreditBalance: 200.0,
	}

	err := cc.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CreditBalance")
}
