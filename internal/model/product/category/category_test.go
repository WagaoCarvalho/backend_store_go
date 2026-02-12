package model

import (
	"testing"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
	"github.com/stretchr/testify/assert"
)

func TestProductCategory_Validate(t *testing.T) {
	validName := "Categoria válida"

	// Nome com 256 caracteres (excede limite)
	longName256 := string(make([]byte, 256))
	for i := range longName256 {
		longName256 = longName256[:i] + "a" + longName256[i+1:]
	}

	// Descrição com 256 caracteres (excede limite)
	longDesc256 := string(make([]byte, 256))
	for i := range longDesc256 {
		longDesc256 = longDesc256[:i] + "x" + longDesc256[i+1:]
	}

	t.Run("Valid - com descrição", func(t *testing.T) {
		pc := &ProductCategory{
			Name:        validName,
			Description: "Descrição válida",
		}
		err := pc.Validate()
		assert.NoError(t, err)
	})

	t.Run("Valid - sem descrição", func(t *testing.T) {
		pc := &ProductCategory{
			Name:        validName,
			Description: "", // Descrição vazia permitida
		}
		err := pc.Validate()
		assert.NoError(t, err)
	})

	t.Run("Valid - nome com exatamente 255 caracteres", func(t *testing.T) {
		exactName255 := string(make([]byte, 255))
		for i := range exactName255 {
			exactName255 = exactName255[:i] + "a" + exactName255[i+1:]
		}

		pc := &ProductCategory{
			Name:        exactName255,
			Description: "Descrição",
		}
		err := pc.Validate()
		assert.NoError(t, err)
	})

	t.Run("Missing Name - string vazia", func(t *testing.T) {
		pc := &ProductCategory{
			Name:        "", // String vazia
			Description: "Qualquer",
		}
		err := pc.Validate()
		assert.Error(t, err)

		var verrs validators.ValidationErrors
		assert.ErrorAs(t, err, &verrs)
		assert.Equal(t, "name", verrs[0].Field)
		assert.Contains(t, verrs[0].Message, "obrigatório")
	})

	t.Run("Missing Name - apenas espaços", func(t *testing.T) {
		pc := &ProductCategory{
			Name:        "   ", // Apenas espaços
			Description: "Qualquer",
		}
		err := pc.Validate()
		assert.Error(t, err)

		var verrs validators.ValidationErrors
		assert.ErrorAs(t, err, &verrs)
		assert.Equal(t, "name", verrs[0].Field)
	})

	t.Run("Name Too Short - 1 caractere", func(t *testing.T) {
		pc := &ProductCategory{
			Name:        "A", // Apenas 1 caractere
			Description: "Descrição",
		}
		err := pc.Validate()
		assert.Error(t, err)

		var verrs validators.ValidationErrors
		assert.ErrorAs(t, err, &verrs)
		assert.Equal(t, "name", verrs[0].Field)
		assert.Contains(t, verrs[0].Message, "mínimo")
	})

	t.Run("Name Too Long - 256 caracteres", func(t *testing.T) {
		pc := &ProductCategory{
			Name:        longName256,
			Description: "Descrição",
		}
		err := pc.Validate()
		assert.Error(t, err)

		var verrs validators.ValidationErrors
		assert.ErrorAs(t, err, &verrs)
		assert.Equal(t, "name", verrs[0].Field)
		assert.Contains(t, verrs[0].Message, "255 caracteres")
	})

	t.Run("Description Too Long", func(t *testing.T) {
		pc := &ProductCategory{
			Name:        validName,
			Description: longDesc256,
		}
		err := pc.Validate()
		assert.Error(t, err)

		var verrs validators.ValidationErrors
		assert.ErrorAs(t, err, &verrs)
		assert.Equal(t, "description", verrs[0].Field)
		assert.Contains(t, verrs[0].Message, "255 caracteres")
	})

	t.Run("Description com espaços - deve ser válido", func(t *testing.T) {
		pc := &ProductCategory{
			Name:        validName,
			Description: "   Descrição com espaços   ",
		}
		err := pc.Validate()
		assert.NoError(t, err)
		// Espaços devem ser removidos pelo TrimSpace
	})

	t.Run("Múltiplos erros - nome vazio e descrição longa", func(t *testing.T) {
		pc := &ProductCategory{
			Name:        "",          // Erro 1
			Description: longDesc256, // Erro 2
		}
		err := pc.Validate()
		assert.Error(t, err)

		var verrs validators.ValidationErrors
		assert.ErrorAs(t, err, &verrs)
		assert.Len(t, verrs, 2)

		// Verifica que contém ambos os erros
		fields := []string{verrs[0].Field, verrs[1].Field}
		assert.Contains(t, fields, "name")
		assert.Contains(t, fields, "description")
	})
}
