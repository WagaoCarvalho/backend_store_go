package dto

import (
	"fmt"
	"time"

	modelFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	modelUser "github.com/WagaoCarvalho/backend_store_go/internal/model/user/filter"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

type UserFilterDTO struct {
	Username    string  `schema:"username"`
	Email       string  `schema:"email"`
	Status      *bool   `schema:"status"`
	CreatedFrom *string `schema:"created_from"`
	CreatedTo   *string `schema:"created_to"`
	UpdatedFrom *string `schema:"updated_from"`
	UpdatedTo   *string `schema:"updated_to"`
	Limit       int     `schema:"limit"`
	Offset      int     `schema:"offset"`
}

func (d *UserFilterDTO) ToModel() (*modelUser.UserFilter, error) {
	// Parse de datas
	parseDate := func(s *string, fieldName string) (*time.Time, error) {
		if s == nil || *s == "" {
			return nil, nil
		}
		t, err := time.Parse("2006-01-02", *s)
		if err != nil {
			return nil, fmt.Errorf(
				"%w: campo '%s' com valor inválido '%s' - formato esperado: YYYY-MM-DD",
				errMsg.ErrInvalidFilter, fieldName, *s,
			)
		}
		return &t, nil
	}

	// Validação de paginação
	if d.Limit < 1 {
		return nil, fmt.Errorf("%w: 'limit' deve ser maior que 0", errMsg.ErrInvalidFilter)
	}
	if d.Limit > 100 {
		return nil, fmt.Errorf("%w: 'limit' máximo é 100", errMsg.ErrInvalidFilter)
	}
	if d.Offset < 0 {
		return nil, fmt.Errorf("%w: 'offset' não pode ser negativo", errMsg.ErrInvalidFilter)
	}

	// Validação de email (se fornecido, deve ser válido)
	if d.Email != "" && !isValidEmailFormat(d.Email) {
		return nil, fmt.Errorf(
			"%w: campo 'email' com formato inválido '%s'",
			errMsg.ErrInvalidFilter, d.Email,
		)
	}
	// NOTA: Email vazio é permitido (não é obrigatório no filtro)

	// Validação básica de username (se fornecido, deve ter no mínimo 3 caracteres)
	if d.Username != "" && len(d.Username) < 3 {
		return nil, fmt.Errorf(
			"%w: 'username' deve ter no mínimo 3 caracteres",
			errMsg.ErrInvalidFilter,
		)
	}

	baseFilter := modelFilter.BaseFilter{
		Limit:  d.Limit,
		Offset: d.Offset,
	}

	createdFrom, err := parseDate(d.CreatedFrom, "created_from")
	if err != nil {
		return nil, err
	}

	createdTo, err := parseDate(d.CreatedTo, "created_to")
	if err != nil {
		return nil, err
	}

	if createdFrom != nil && createdTo != nil && createdFrom.After(*createdTo) {
		return nil, fmt.Errorf(
			"%w: 'created_from' não pode ser depois de 'created_to'",
			errMsg.ErrInvalidFilter,
		)
	}

	updatedFrom, err := parseDate(d.UpdatedFrom, "updated_from")
	if err != nil {
		return nil, err
	}

	updatedTo, err := parseDate(d.UpdatedTo, "updated_to")
	if err != nil {
		return nil, err
	}

	if updatedFrom != nil && updatedTo != nil && updatedFrom.After(*updatedTo) {
		return nil, fmt.Errorf(
			"%w: 'updated_from' não pode ser depois de 'updated_to'",
			errMsg.ErrInvalidFilter,
		)
	}

	filter := &modelUser.UserFilter{
		BaseFilter: baseFilter,

		Username: d.Username,
		Email:    d.Email,
		Status:   d.Status,

		CreatedFrom: createdFrom,
		CreatedTo:   createdTo,
		UpdatedFrom: updatedFrom,
		UpdatedTo:   updatedTo,
	}

	return filter, nil
}

// Função auxiliar para validação básica de formato de email
func isValidEmailFormat(email string) bool {
	// Email não pode ser vazio
	if email == "" {
		return false
	}

	// Verifica se tem exatamente um @
	atCount := 0
	atIndex := -1

	for i, c := range email {
		if c == '@' {
			atCount++
			atIndex = i
		}
	}

	// Deve ter exatamente um @
	if atCount != 1 {
		return false
	}

	// @ não pode estar no início ou no fim
	if atIndex == 0 || atIndex == len(email)-1 {
		return false
	}

	// Deve ter pelo menos um ponto após o @
	hasDotAfterAt := false
	for i := atIndex + 1; i < len(email); i++ {
		if email[i] == '.' {
			// O ponto não pode ser imediatamente após o @
			if i == atIndex+1 {
				return false
			}
			// O ponto não pode estar no final
			if i == len(email)-1 {
				return false
			}
			hasDotAfterAt = true
			break
		}
	}

	return hasDotAfterAt
}
