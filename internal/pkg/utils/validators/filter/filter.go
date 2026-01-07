package validators

import (
	"time"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type FilterValidator interface {
	Validate() error
}

type DateRangeValidator struct {
	FromField string
	ToField   string
	From      *time.Time
	To        *time.Time
}

func (v *DateRangeValidator) Validate() error {
	if v.From != nil && v.To != nil && v.From.After(*v.To) {
		return &validators.ValidationError{
			Field:   v.FromField + "/" + v.ToField,
			Message: "intervalo inválido",
		}
	}
	return nil
}
