package dto

import (
	client_credit "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/credit"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

type ClientCreditDTO struct {
	ID            *int64  `json:"id,omitempty"`
	ClientID      int64   `json:"client_id"`
	AllowCredit   bool    `json:"allow_credit"`
	CreditLimit   float64 `json:"credit_limit"`
	CreditBalance float64 `json:"credit_balance"`
}

func ToClientCreditModel(dto ClientCreditDTO) *client_credit.ClientCredit {
	return &client_credit.ClientCredit{
		ID:            utils.NilToZero(dto.ID),
		ClientID:      dto.ClientID,
		AllowCredit:   dto.AllowCredit,
		CreditLimit:   dto.CreditLimit,
		CreditBalance: dto.CreditBalance,
	}
}

func ToClientCreditDTO(m *client_credit.ClientCredit) ClientCreditDTO {
	if m == nil {
		return ClientCreditDTO{}
	}

	return ClientCreditDTO{
		ID:            &m.ID,
		ClientID:      m.ClientID,
		AllowCredit:   m.AllowCredit,
		CreditLimit:   m.CreditLimit,
		CreditBalance: m.CreditBalance,
	}
}
