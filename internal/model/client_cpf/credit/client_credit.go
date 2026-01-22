package model

import (
	"time"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type ClientCredit struct {
	ID            int64
	ClientID      int64
	AllowCredit   bool
	CreditLimit   float64
	CreditBalance float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (cc *ClientCredit) Validate() error {
	var validationErrors []validators.ValidationError

	// Validar ClientID
	if cc.ClientID <= 0 {
		validationErrors = append(validationErrors, validators.ValidationError{
			Field:   "client_id",
			Message: "ID do cliente é obrigatório",
		})
	}

	// Validar CreditLimit
	if cc.CreditLimit < 0 {
		validationErrors = append(validationErrors, validators.ValidationError{
			Field:   "credit_limit",
			Message: "limite de crédito não pode ser negativo",
		})
	}

	// Validar CreditBalance
	if cc.CreditBalance < 0 {
		validationErrors = append(validationErrors, validators.ValidationError{
			Field:   "credit_balance",
			Message: "saldo de crédito não pode ser negativo",
		})
	}

	// Validar relação entre balance e limite (apenas se ambos forem válidos)
	if cc.CreditLimit >= 0 && cc.CreditBalance >= 0 {
		if cc.CreditBalance > cc.CreditLimit {
			validationErrors = append(validationErrors, validators.ValidationError{
				Field:   "credit_balance",
				Message: "saldo não pode ser maior que o limite de crédito",
			})
		}
	}

	// Retornar erros usando a função helper do seu pacote
	return validators.NewValidationErrors(validationErrors)
}

// ValidateForUpdate validações específicas para atualização
func (cc *ClientCredit) ValidateForUpdate(oldCredit *ClientCredit) error {
	// Validação básica primeiro
	if err := cc.Validate(); err != nil {
		return err
	}

	var validationErrors []validators.ValidationError

	// Validar que não está reduzindo o limite abaixo do saldo atual
	if oldCredit != nil && cc.CreditLimit < oldCredit.CreditBalance {
		validationErrors = append(validationErrors, validators.ValidationError{
			Field:   "credit_limit",
			Message: "novo limite não pode ser menor que o saldo atual",
		})
	}

	return validators.NewValidationErrors(validationErrors)
}

// CanUseCredit verifica se pode usar crédito
func (cc *ClientCredit) CanUseCredit(amount float64) (bool, error) {
	var validationErrors []validators.ValidationError

	if !cc.AllowCredit {
		validationErrors = append(validationErrors, validators.ValidationError{
			Field:   "allow_credit",
			Message: "cliente não tem crédito habilitado",
		})
	}

	if amount <= 0 {
		validationErrors = append(validationErrors, validators.ValidationError{
			Field:   "amount",
			Message: "valor deve ser positivo",
		})
	}

	available := cc.CreditLimit - cc.CreditBalance
	if amount > available {
		validationErrors = append(validationErrors, validators.ValidationError{
			Field:   "amount",
			Message: "valor excede o crédito disponível",
		})
	}

	if err := validators.NewValidationErrors(validationErrors); err != nil {
		return false, err
	}

	return true, nil
}

// UseCredit utiliza crédito
func (cc *ClientCredit) UseCredit(amount float64) error {
	canUse, err := cc.CanUseCredit(amount)
	if !canUse {
		return err
	}

	cc.CreditBalance += amount
	cc.UpdatedAt = time.Now()

	return nil
}

// AddCredit adiciona crédito (pagamento, por exemplo)
func (cc *ClientCredit) AddCredit(amount float64) error {
	var validationErrors []validators.ValidationError

	if amount <= 0 {
		validationErrors = append(validationErrors, validators.ValidationError{
			Field:   "amount",
			Message: "valor deve ser positivo",
		})
	}

	newBalance := cc.CreditBalance - amount
	if newBalance < 0 {
		validationErrors = append(validationErrors, validators.ValidationError{
			Field:   "amount",
			Message: "valor excede o saldo atual",
		})
	}

	if err := validators.NewValidationErrors(validationErrors); err != nil {
		return err
	}

	cc.CreditBalance = newBalance
	cc.UpdatedAt = time.Now()

	return nil
}

// AvailableCredit retorna o crédito disponível
func (cc *ClientCredit) AvailableCredit() float64 {
	if !cc.AllowCredit {
		return 0
	}
	return cc.CreditLimit - cc.CreditBalance
}

// IsCreditAvailable verifica se há crédito disponível para um valor específico
func (cc *ClientCredit) IsCreditAvailable(amount float64) bool {
	if !cc.AllowCredit {
		return false
	}
	return cc.AvailableCredit() >= amount
}

// SetCreditLimit define um novo limite de crédito com validação
func (cc *ClientCredit) SetCreditLimit(newLimit float64) error {
	var validationErrors []validators.ValidationError

	if newLimit < 0 {
		validationErrors = append(validationErrors, validators.ValidationError{
			Field:   "credit_limit",
			Message: "limite de crédito não pode ser negativo",
		})
	}

	// Verificar se o novo limite é menor que o saldo atual
	if newLimit < cc.CreditBalance {
		validationErrors = append(validationErrors, validators.ValidationError{
			Field:   "credit_limit",
			Message: "novo limite não pode ser menor que o saldo atual",
		})
	}

	if err := validators.NewValidationErrors(validationErrors); err != nil {
		return err
	}

	cc.CreditLimit = newLimit
	cc.UpdatedAt = time.Now()

	return nil
}
