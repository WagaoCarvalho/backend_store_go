package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationError_Error(t *testing.T) {
	err := ValidationError{Field: "email", Message: "campo obrigatório"}
	assert.Equal(t, "erro no campo 'email': campo obrigatório", err.Error())
}

func TestValidationErrors_Error(t *testing.T) {
	errs := ValidationErrors{
		{Field: "email", Message: "campo obrigatório"},
		{Field: "password", Message: "mínimo de 6 caracteres"},
	}
	expected := "erro no campo 'email': campo obrigatório; erro no campo 'password': mínimo de 6 caracteres; "
	assert.Equal(t, expected, errs.Error())
}

func TestValidationErrors_HasErrors(t *testing.T) {
	var empty ValidationErrors
	assert.False(t, empty.HasErrors())

	errs := ValidationErrors{
		{Field: "email", Message: "campo obrigatório"},
	}
	assert.True(t, errs.HasErrors())
}

func TestIsBlank(t *testing.T) {
	assert.True(t, IsBlank(""))
	assert.True(t, IsBlank("   "))
	assert.False(t, IsBlank("a"))
	assert.False(t, IsBlank("  a  "))
}

func TestValidateSingleNonNil(t *testing.T) {
	a, b, c := int64(1), int64(2), int64(3)

	// exatamente 1 não-nulo → true
	assert.True(t, ValidateSingleNonNil(&a))
	assert.True(t, ValidateSingleNonNil(&b))
	assert.True(t, ValidateSingleNonNil(&c))

	// zero ou mais de 1 não-nulo → false
	assert.False(t, ValidateSingleNonNil(nil))
	assert.False(t, ValidateSingleNonNil(nil, nil))
	assert.False(t, ValidateSingleNonNil(&a, &b))
	assert.False(t, ValidateSingleNonNil(&a, &b, &c))
	assert.False(t, ValidateSingleNonNil(nil, &a, &b))
	assert.False(t, ValidateSingleNonNil(nil, nil, &c, &b))
}
