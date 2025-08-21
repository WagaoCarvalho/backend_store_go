package models

import (
	"testing"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators"
	"github.com/stretchr/testify/assert"
)

func TestUserCategory_Validate(t *testing.T) {
	validName := "Categoria válida"
	longName := string(make([]byte, 101))
	for i := range longName {
		longName = longName[:i] + "a" + longName[i+1:]
	}

	t.Run("Valid", func(t *testing.T) {
		uc := &UserCategory{
			Name:        validName,
			Description: "Descrição válida",
		}
		err := uc.Validate()
		assert.NoError(t, err)
	})

	t.Run("Missing Name", func(t *testing.T) {
		uc := &UserCategory{
			Name:        "  ", // em branco
			Description: "Qualquer",
		}
		err := uc.Validate()
		assert.Error(t, err)
		verr, ok := err.(*validators.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "Name", verr.Field)
	})

	t.Run("Name Too Long", func(t *testing.T) {
		uc := &UserCategory{
			Name:        longName,
			Description: "Descrição",
		}
		err := uc.Validate()
		assert.Error(t, err)
		verr, ok := err.(*validators.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "Name", verr.Field)
	})

	t.Run("Description Too Long", func(t *testing.T) {
		longDesc := ""
		for i := 0; i < 260; i++ {
			longDesc += "x"
		}
		uc := &UserCategory{
			Name:        validName,
			Description: longDesc,
		}
		err := uc.Validate()
		assert.Error(t, err)
		verr, ok := err.(*validators.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "Description", verr.Field)
	})
}
