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
	assert.Contains(t, err.Error(), "client_id")
}

func TestClientCredit_Validate_NegativeLimit(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      1,
		CreditLimit:   -10.0,
		CreditBalance: 0,
	}

	err := cc.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "credit_limit")
}

func TestClientCredit_Validate_NegativeBalance(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      1,
		CreditLimit:   100.0,
		CreditBalance: -5.0,
	}

	err := cc.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "credit_balance")
}

func TestClientCredit_Validate_BalanceGreaterThanLimit(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      1,
		CreditLimit:   100.0,
		CreditBalance: 200.0,
	}

	err := cc.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "credit_balance")
}

func TestClientCredit_ValidateForUpdate_Valid(t *testing.T) {
	oldCredit := &ClientCredit{
		ClientID:      1,
		CreditLimit:   1000.0,
		CreditBalance: 500.0,
	}

	newCredit := &ClientCredit{
		ClientID:      1,
		CreditLimit:   1500.0,
		CreditBalance: 500.0,
	}

	err := newCredit.ValidateForUpdate(oldCredit)
	assert.NoError(t, err)
}

func TestClientCredit_ValidateForUpdate_LimitBelowBalance(t *testing.T) {
	oldCredit := &ClientCredit{
		ClientID:      1,
		CreditLimit:   1000.0,
		CreditBalance: 800.0,
	}

	newCredit := &ClientCredit{
		ClientID:      1,
		CreditLimit:   700.0,
		CreditBalance: 500.0,
	}

	err := newCredit.ValidateForUpdate(oldCredit)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "credit_limit")
}

func TestClientCredit_CanUseCredit_Success(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      1,
		AllowCredit:   true,
		CreditLimit:   1000.0,
		CreditBalance: 500.0,
	}

	canUse, err := cc.CanUseCredit(200.0)
	assert.True(t, canUse)
	assert.NoError(t, err)
}

func TestClientCredit_CanUseCredit_CreditNotAllowed(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      1,
		AllowCredit:   false,
		CreditLimit:   1000.0,
		CreditBalance: 500.0,
	}

	canUse, err := cc.CanUseCredit(200.0)
	assert.False(t, canUse)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "allow_credit")
}

func TestClientCredit_CanUseCredit_InsufficientCredit(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      1,
		AllowCredit:   true,
		CreditLimit:   1000.0,
		CreditBalance: 900.0,
	}

	canUse, err := cc.CanUseCredit(200.0)
	assert.False(t, canUse)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "amount")
}

func TestClientCredit_CanUseCredit_InvalidAmount(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      1,
		AllowCredit:   true,
		CreditLimit:   1000.0,
		CreditBalance: 500.0,
	}

	canUse, err := cc.CanUseCredit(0)
	assert.False(t, canUse)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "amount")
}

func TestClientCredit_CanUseCredit_NegativeAmount(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      1,
		AllowCredit:   true,
		CreditLimit:   1000.0,
		CreditBalance: 500.0,
	}

	canUse, err := cc.CanUseCredit(-50.0)
	assert.False(t, canUse)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "amount")
}

func TestClientCredit_UseCredit_Success(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      1,
		AllowCredit:   true,
		CreditLimit:   1000.0,
		CreditBalance: 500.0,
	}

	initialBalance := cc.CreditBalance
	err := cc.UseCredit(200.0)
	assert.NoError(t, err)
	assert.Equal(t, initialBalance+200.0, cc.CreditBalance)
}

func TestClientCredit_UseCredit_Failure(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      1,
		AllowCredit:   false,
		CreditLimit:   1000.0,
		CreditBalance: 500.0,
	}

	initialBalance := cc.CreditBalance
	err := cc.UseCredit(200.0)
	assert.Error(t, err)
	assert.Equal(t, initialBalance, cc.CreditBalance)
}

func TestClientCredit_AddCredit_Success(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      1,
		AllowCredit:   true,
		CreditLimit:   1000.0,
		CreditBalance: 500.0,
	}

	initialBalance := cc.CreditBalance
	err := cc.AddCredit(200.0)
	assert.NoError(t, err)
	assert.Equal(t, initialBalance-200.0, cc.CreditBalance)
}

func TestClientCredit_AddCredit_InvalidAmount(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      1,
		AllowCredit:   true,
		CreditLimit:   1000.0,
		CreditBalance: 500.0,
	}

	initialBalance := cc.CreditBalance
	err := cc.AddCredit(0)
	assert.Error(t, err)
	assert.Equal(t, initialBalance, cc.CreditBalance)
}

func TestClientCredit_AddCredit_ExceedsBalance(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      1,
		AllowCredit:   true,
		CreditLimit:   1000.0,
		CreditBalance: 500.0,
	}

	initialBalance := cc.CreditBalance
	err := cc.AddCredit(600.0)
	assert.Error(t, err)
	assert.Equal(t, initialBalance, cc.CreditBalance)
}

func TestClientCredit_AvailableCredit(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      1,
		AllowCredit:   true,
		CreditLimit:   1000.0,
		CreditBalance: 500.0,
	}

	assert.Equal(t, 500.0, cc.AvailableCredit())
}

func TestClientCredit_AvailableCredit_NotAllowed(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      1,
		AllowCredit:   false,
		CreditLimit:   1000.0,
		CreditBalance: 500.0,
	}

	assert.Equal(t, 0.0, cc.AvailableCredit())
}

func TestClientCredit_IsCreditAvailable(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      1,
		AllowCredit:   true,
		CreditLimit:   1000.0,
		CreditBalance: 500.0,
	}

	assert.True(t, cc.IsCreditAvailable(300.0))
	assert.True(t, cc.IsCreditAvailable(500.0))
	assert.False(t, cc.IsCreditAvailable(600.0))
}

func TestClientCredit_SetCreditLimit_Success(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      1,
		CreditLimit:   1000.0,
		CreditBalance: 500.0,
	}

	err := cc.SetCreditLimit(1500.0)
	assert.NoError(t, err)
	assert.Equal(t, 1500.0, cc.CreditLimit)
}

func TestClientCredit_SetCreditLimit_BelowBalance(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      1,
		CreditLimit:   1000.0,
		CreditBalance: 500.0,
	}

	err := cc.SetCreditLimit(400.0)
	assert.Error(t, err)
	assert.Equal(t, 1000.0, cc.CreditLimit)
}

func TestClientCredit_SetCreditLimit_Negative(t *testing.T) {
	cc := &ClientCredit{
		ClientID:      1,
		CreditLimit:   1000.0,
		CreditBalance: 500.0,
	}

	err := cc.SetCreditLimit(-100.0)
	assert.Error(t, err)
	assert.Equal(t, 1000.0, cc.CreditLimit)
}
