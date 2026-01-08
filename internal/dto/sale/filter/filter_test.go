package dto

import (
	"testing"

	modelFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaleFilterDTO_ToModel(t *testing.T) {
	// Testes de sucesso
	t.Run("Converte todos os campos preenchidos corretamente", func(t *testing.T) {
		clientID := int64(100)
		userID := int64(50)
		minTotalAmount := "150.75"
		maxTotalAmount := "500.00"
		minItemsAmount := "100.00"
		maxItemsAmount := "300.00"
		minItemsDiscount := "5.0"
		maxItemsDiscount := "15.5"
		minSaleDiscount := "10.0"
		maxSaleDiscount := "20.0"
		saleDateFrom := "2024-01-01"
		saleDateTo := "2024-01-31"
		createdFrom := "2024-01-01"
		createdTo := "2024-12-31"
		updatedFrom := "2024-06-01"
		updatedTo := "2024-06-30"

		dto := SaleFilterDTO{
			ClientID:         &clientID,
			UserID:           &userID,
			PaymentType:      "credit",
			Status:           "completed",
			MinTotalAmount:   &minTotalAmount,
			MaxTotalAmount:   &maxTotalAmount,
			MinItemsAmount:   &minItemsAmount,
			MaxItemsAmount:   &maxItemsAmount,
			MinItemsDiscount: &minItemsDiscount,
			MaxItemsDiscount: &maxItemsDiscount,
			MinSaleDiscount:  &minSaleDiscount,
			MaxSaleDiscount:  &maxSaleDiscount,
			SaleDateFrom:     &saleDateFrom,
			SaleDateTo:       &saleDateTo,
			CreatedFrom:      &createdFrom,
			CreatedTo:        &createdTo,
			UpdatedFrom:      &updatedFrom,
			UpdatedTo:        &updatedTo,
			Limit:            20,
			Offset:           10,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		require.NotNil(t, model)

		// Verificar campos básicos
		assert.Equal(t, &clientID, model.ClientID)
		assert.Equal(t, &userID, model.UserID)
		assert.Equal(t, "credit", model.PaymentType)
		assert.Equal(t, "completed", model.Status)
		assert.Equal(t, modelFilter.BaseFilter{Limit: 20, Offset: 10}, model.BaseFilter)

		// Verificar valores float
		require.NotNil(t, model.MinTotalAmount)
		require.NotNil(t, model.MaxTotalAmount)
		assert.Equal(t, 150.75, *model.MinTotalAmount)
		assert.Equal(t, 500.00, *model.MaxTotalAmount)

		// Verificar outros valores float
		require.NotNil(t, model.MinTotalItemsAmount)
		require.NotNil(t, model.MaxTotalItemsAmount)
		assert.Equal(t, 100.00, *model.MinTotalItemsAmount)
		assert.Equal(t, 300.00, *model.MaxTotalItemsAmount)

		require.NotNil(t, model.MinTotalItemsDiscount)
		require.NotNil(t, model.MaxTotalItemsDiscount)
		assert.Equal(t, 5.0, *model.MinTotalItemsDiscount)
		assert.Equal(t, 15.5, *model.MaxTotalItemsDiscount)

		require.NotNil(t, model.MinTotalSaleDiscount)
		require.NotNil(t, model.MaxTotalSaleDiscount)
		assert.Equal(t, 10.0, *model.MinTotalSaleDiscount)
		assert.Equal(t, 20.0, *model.MaxTotalSaleDiscount)

		// Verificar datas
		require.NotNil(t, model.SaleDateFrom)
		require.NotNil(t, model.SaleDateTo)
		require.NotNil(t, model.CreatedFrom)
		require.NotNil(t, model.CreatedTo)
		require.NotNil(t, model.UpdatedFrom)
		require.NotNil(t, model.UpdatedTo)

		assert.Equal(t, "2024-01-01", model.SaleDateFrom.Format("2006-01-02"))
		assert.Equal(t, "2024-01-31", model.SaleDateTo.Format("2006-01-02"))
		assert.Equal(t, "2024-01-01", model.CreatedFrom.Format("2006-01-02"))
		assert.Equal(t, "2024-12-31", model.CreatedTo.Format("2006-01-02"))
		assert.Equal(t, "2024-06-01", model.UpdatedFrom.Format("2006-01-02"))
		assert.Equal(t, "2024-06-30", model.UpdatedTo.Format("2006-01-02"))
	})

	t.Run("Retorna nil para campos vazios ou não preenchidos", func(t *testing.T) {
		dto := SaleFilterDTO{
			Limit:  10, // Agora precisa de um limit válido (>0)
			Offset: 0,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		require.NotNil(t, model)

		// Verificar que ponteiros são nil
		assert.Nil(t, model.ClientID)
		assert.Nil(t, model.UserID)
		assert.Nil(t, model.MinTotalAmount)
		assert.Nil(t, model.MaxTotalAmount)
		assert.Nil(t, model.MinTotalItemsAmount)
		assert.Nil(t, model.MaxTotalItemsAmount)
		assert.Nil(t, model.MinTotalItemsDiscount)
		assert.Nil(t, model.MaxTotalItemsDiscount)
		assert.Nil(t, model.MinTotalSaleDiscount)
		assert.Nil(t, model.MaxTotalSaleDiscount)
		assert.Nil(t, model.SaleDateFrom)
		assert.Nil(t, model.SaleDateTo)
		assert.Nil(t, model.CreatedFrom)
		assert.Nil(t, model.CreatedTo)
		assert.Nil(t, model.UpdatedFrom)
		assert.Nil(t, model.UpdatedTo)

		// Verificar campos de string vazios
		assert.Equal(t, "", model.PaymentType)
		assert.Equal(t, "", model.Status)

		// Verificar BaseFilter com valores fornecidos
		assert.Equal(t, 10, model.Limit)
		assert.Equal(t, 0, model.Offset)
	})

	t.Run("Converte strings vazias para nil", func(t *testing.T) {
		emptyString := ""
		clientID := int64(1)

		dto := SaleFilterDTO{
			ClientID:       &clientID,
			MinTotalAmount: &emptyString,
			CreatedFrom:    &emptyString,
			Limit:          10,
			Offset:         0,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		require.NotNil(t, model)

		// ClientID deve ser mantido
		assert.Equal(t, &clientID, model.ClientID)

		// Strings vazias devem resultar em nil
		assert.Nil(t, model.MinTotalAmount)
		assert.Nil(t, model.CreatedFrom)
	})

	t.Run("Mantém valores zero válidos", func(t *testing.T) {
		zeroFloat := "0.0"
		validDate := "2024-01-01"
		clientID := int64(0) // Zero é um ID válido

		dto := SaleFilterDTO{
			ClientID:       &clientID,
			MinTotalAmount: &zeroFloat,
			CreatedFrom:    &validDate,
			Limit:          10,
			Offset:         0,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		require.NotNil(t, model)

		// Zero values devem ser mantidos
		assert.Equal(t, &clientID, model.ClientID)
		require.NotNil(t, model.MinTotalAmount)
		assert.Equal(t, 0.0, *model.MinTotalAmount)
		assert.NotNil(t, model.CreatedFrom)
		assert.Equal(t, "2024-01-01", model.CreatedFrom.Format("2006-01-02"))
	})

	t.Run("Converte apenas alguns campos preenchidos", func(t *testing.T) {
		clientID := int64(10)
		minTotalAmount := "100.00"
		saleDateFrom := "2024-03-01"

		dto := SaleFilterDTO{
			ClientID:       &clientID,
			PaymentType:    "cash",
			MinTotalAmount: &minTotalAmount,
			SaleDateFrom:   &saleDateFrom,
			Limit:          15,
			Offset:         5,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		require.NotNil(t, model)

		// Campos preenchidos
		assert.Equal(t, &clientID, model.ClientID)
		assert.Equal(t, "cash", model.PaymentType)
		require.NotNil(t, model.MinTotalAmount)
		assert.Equal(t, 100.00, *model.MinTotalAmount)
		require.NotNil(t, model.SaleDateFrom)
		assert.Equal(t, "2024-03-01", model.SaleDateFrom.Format("2006-01-02"))
		assert.Equal(t, modelFilter.BaseFilter{Limit: 15, Offset: 5}, model.BaseFilter)

		// Campos não preenchidos devem ser nil ou vazios
		assert.Nil(t, model.UserID)
		assert.Equal(t, "", model.Status)
		assert.Nil(t, model.MaxTotalAmount)
		assert.Nil(t, model.MinTotalItemsAmount)
		assert.Nil(t, model.MaxTotalItemsAmount)
		assert.Nil(t, model.MinTotalItemsDiscount)
		assert.Nil(t, model.MaxTotalItemsDiscount)
		assert.Nil(t, model.MinTotalSaleDiscount)
		assert.Nil(t, model.MaxTotalSaleDiscount)
		assert.Nil(t, model.SaleDateTo)
		assert.Nil(t, model.CreatedFrom)
		assert.Nil(t, model.CreatedTo)
		assert.Nil(t, model.UpdatedFrom)
		assert.Nil(t, model.UpdatedTo)
	})

	t.Run("Aceita status vazio como válido", func(t *testing.T) {
		dto := SaleFilterDTO{
			Status: "",
			Limit:  10,
			Offset: 0,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		assert.Equal(t, "", model.Status)
	})

	t.Run("Aceita payment_type vazio como válido", func(t *testing.T) {
		dto := SaleFilterDTO{
			PaymentType: "",
			Limit:       10,
			Offset:      0,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		assert.Equal(t, "", model.PaymentType)
	})

	// Testes de erro - Validação de paginação
	t.Run("Retorna erro quando limit é zero", func(t *testing.T) {
		dto := SaleFilterDTO{
			Limit:  0,
			Offset: 0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "'limit' deve ser maior que 0")
	})

	t.Run("Retorna erro quando limit é negativo", func(t *testing.T) {
		dto := SaleFilterDTO{
			Limit:  -1,
			Offset: 0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "'limit' deve ser maior que 0")
	})

	t.Run("Retorna erro quando limit é maior que máximo", func(t *testing.T) {
		dto := SaleFilterDTO{
			Limit:  101,
			Offset: 0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "'limit' máximo é 100")
	})

	t.Run("Retorna erro quando offset é negativo", func(t *testing.T) {
		dto := SaleFilterDTO{
			Limit:  10,
			Offset: -1,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "'offset' não pode ser negativo")
	})

	// Testes de erro - Validação de datas
	t.Run("Retorna erro para data inválida no formato", func(t *testing.T) {
		invalidDate := "31/12/2024"

		dto := SaleFilterDTO{
			SaleDateFrom: &invalidDate,
			Limit:        10,
			Offset:       0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "sale_date_from")
		assert.Contains(t, err.Error(), "formato esperado: YYYY-MM-DD")
	})

	t.Run("Retorna erro para datas incoerentes - sale_date", func(t *testing.T) {
		saleDateFrom := "2024-02-01"
		saleDateTo := "2024-01-01" // Data anterior

		dto := SaleFilterDTO{
			SaleDateFrom: &saleDateFrom,
			SaleDateTo:   &saleDateTo,
			Limit:        10,
			Offset:       0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "'sale_date_from' não pode ser depois de 'sale_date_to'")
	})

	t.Run("Retorna erro para datas incoerentes - created", func(t *testing.T) {
		createdFrom := "2024-12-01"
		createdTo := "2024-01-01" // Data anterior

		dto := SaleFilterDTO{
			CreatedFrom: &createdFrom,
			CreatedTo:   &createdTo,
			Limit:       10,
			Offset:      0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "'created_from' não pode ser depois de 'created_to'")
	})

	t.Run("Retorna erro para datas incoerentes - updated", func(t *testing.T) {
		updatedFrom := "2024-06-02"
		updatedTo := "2024-06-01" // Data anterior

		dto := SaleFilterDTO{
			UpdatedFrom: &updatedFrom,
			UpdatedTo:   &updatedTo,
			Limit:       10,
			Offset:      0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "'updated_from' não pode ser depois de 'updated_to'")
	})

	// Testes de erro - Validação de valores numéricos
	t.Run("Retorna erro para valor numérico inválido", func(t *testing.T) {
		invalidFloat := "abc"

		dto := SaleFilterDTO{
			MinTotalAmount: &invalidFloat,
			Limit:          10,
			Offset:         0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "min_total_amount")
		assert.Contains(t, err.Error(), "valor numérico esperado")
	})

	t.Run("Retorna erro para valores monetários negativos", func(t *testing.T) {
		negativeAmount := "-100.00"

		dto := SaleFilterDTO{
			MinTotalAmount: &negativeAmount,
			Limit:          10,
			Offset:         0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "não pode ser negativo")
		assert.Contains(t, err.Error(), "min_total_amount")
	})

	t.Run("Retorna erro para intervalo numérico incoerente - total_amount", func(t *testing.T) {
		minTotalAmount := "200.00"
		maxTotalAmount := "100.00" // Menor que min

		dto := SaleFilterDTO{
			MinTotalAmount: &minTotalAmount,
			MaxTotalAmount: &maxTotalAmount,
			Limit:          10,
			Offset:         0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "'min_total_amount' não pode ser maior que 'max_total_amount'")
	})

	t.Run("Retorna erro para intervalo numérico incoerente - items_amount", func(t *testing.T) {
		minItemsAmount := "300.00"
		maxItemsAmount := "100.00" // Menor que min

		dto := SaleFilterDTO{
			MinItemsAmount: &minItemsAmount,
			MaxItemsAmount: &maxItemsAmount,
			Limit:          10,
			Offset:         0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "'min_items_amount' não pode ser maior que 'max_items_amount'")
	})

	t.Run("Retorna erro para intervalo numérico incoerente - items_discount", func(t *testing.T) {
		minItemsDiscount := "20.0"
		maxItemsDiscount := "10.0" // Menor que min

		dto := SaleFilterDTO{
			MinItemsDiscount: &minItemsDiscount,
			MaxItemsDiscount: &maxItemsDiscount,
			Limit:            10,
			Offset:           0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "'min_items_discount' não pode ser maior que 'max_items_discount'")
	})

	t.Run("Retorna erro para intervalo numérico incoerente - sale_discount", func(t *testing.T) {
		minSaleDiscount := "30.0"
		maxSaleDiscount := "15.0" // Menor que min

		dto := SaleFilterDTO{
			MinSaleDiscount: &minSaleDiscount,
			MaxSaleDiscount: &maxSaleDiscount,
			Limit:           10,
			Offset:          0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "'min_sale_discount' não pode ser maior que 'max_sale_discount'")
	})

	// Testes de erro - Validação de enums
	t.Run("Retorna erro para payment_type inválido", func(t *testing.T) {
		dto := SaleFilterDTO{
			PaymentType: "INVALID_PAYMENT",
			Limit:       10,
			Offset:      0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "'payment_type' com valor inválido 'INVALID_PAYMENT'")
	})

	t.Run("Retorna erro para status inválido", func(t *testing.T) {
		dto := SaleFilterDTO{
			Status: "INVALID_STATUS",
			Limit:  10,
			Offset: 0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "'status' com valor inválido 'INVALID_STATUS'")
	})

	t.Run("Aceita payment_types válidos", func(t *testing.T) {
		validPaymentTypes := []string{"credit", "debit", "cash", "pix", "bank_slip", ""}

		for _, paymentType := range validPaymentTypes {
			t.Run(paymentType, func(t *testing.T) {
				dto := SaleFilterDTO{
					PaymentType: paymentType,
					Limit:       10,
					Offset:      0,
				}

				model, err := dto.ToModel()
				require.NoError(t, err)
				assert.Equal(t, paymentType, model.PaymentType)
			})
		}
	})

	t.Run("Aceita status válidos", func(t *testing.T) {
		validStatuses := []string{"pending", "completed", "cancelled", "refunded", ""}

		for _, status := range validStatuses {
			t.Run(status, func(t *testing.T) {
				dto := SaleFilterDTO{
					Status: status,
					Limit:  10,
					Offset: 0,
				}

				model, err := dto.ToModel()
				require.NoError(t, err)
				assert.Equal(t, status, model.Status)
			})
		}
	})

	// Testes de borda/canto
	t.Run("Valores no limite máximo permitido", func(t *testing.T) {
		dto := SaleFilterDTO{
			Limit:  100, // Valor máximo
			Offset: 0,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		assert.Equal(t, 100, model.Limit)
	})

	t.Run("Valores no limite mínimo permitido", func(t *testing.T) {
		dto := SaleFilterDTO{
			Limit:  1, // Valor mínimo
			Offset: 0,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		assert.Equal(t, 1, model.Limit)
	})

	t.Run("Permite offset zero", func(t *testing.T) {
		dto := SaleFilterDTO{
			Limit:  10,
			Offset: 0, // Zero é permitido
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		assert.Equal(t, 0, model.Offset)
	})

	t.Run("Campos nil são tratados corretamente", func(t *testing.T) {
		dto := SaleFilterDTO{
			Limit:  10,
			Offset: 0,
			// Todos os campos de ponteiro são nil por padrão
		}

		model, err := dto.ToModel()
		require.NoError(t, err)

		// Verificar que todos os ponteiros são nil
		assert.Nil(t, model.ClientID)
		assert.Nil(t, model.UserID)
		assert.Nil(t, model.MinTotalAmount)
		assert.Nil(t, model.MaxTotalAmount)
		assert.Nil(t, model.MinTotalItemsAmount)
		assert.Nil(t, model.MaxTotalItemsAmount)
		assert.Nil(t, model.MinTotalItemsDiscount)
		assert.Nil(t, model.MaxTotalItemsDiscount)
		assert.Nil(t, model.MinTotalSaleDiscount)
		assert.Nil(t, model.MaxTotalSaleDiscount)
		assert.Nil(t, model.SaleDateFrom)
		assert.Nil(t, model.SaleDateTo)
		assert.Nil(t, model.CreatedFrom)
		assert.Nil(t, model.CreatedTo)
		assert.Nil(t, model.UpdatedFrom)
		assert.Nil(t, model.UpdatedTo)
	})

	t.Run("Valores decimais com diferentes formatos", func(t *testing.T) {
		testCases := []struct {
			name     string
			value    string
			expected float64
		}{
			{"inteiro", "100", 100.0},
			{"um decimal", "100.5", 100.5},
			{"dois decimais", "100.50", 100.5},
			{"três decimais", "100.555", 100.555},
			{"zero", "0", 0.0},
			{"zero decimal", "0.0", 0.0},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				dto := SaleFilterDTO{
					MinTotalAmount: &tc.value,
					Limit:          10,
					Offset:         0,
				}

				model, err := dto.ToModel()
				require.NoError(t, err)
				require.NotNil(t, model.MinTotalAmount)
				assert.Equal(t, tc.expected, *model.MinTotalAmount)
			})
		}
	})

	t.Run("Primeiro erro encontrado é retornado", func(t *testing.T) {
		// Testa que apenas o primeiro erro é retornado
		// (a implementação pode parar no primeiro erro ou coletar todos)
		invalidDate := "invalid"
		negativeLimit := -1

		dto := SaleFilterDTO{
			SaleDateFrom: &invalidDate,
			Limit:        negativeLimit,
			Offset:       0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		// Pode retornar erro de data ou de limit, dependendo da ordem de validação
		assert.Contains(t, err.Error(), "filtro inválido")
	})

	t.Run("Campos opcionais nil não causam erro", func(t *testing.T) {
		// Testa que campos que são ponteiros e podem ser nil
		// não causam erro quando são nil
		var nilString *string = nil

		dto := SaleFilterDTO{
			MinTotalAmount: nilString, // Explicitamente nil
			SaleDateFrom:   nilString,
			Limit:          10,
			Offset:         0,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		assert.Nil(t, model.MinTotalAmount)
		assert.Nil(t, model.SaleDateFrom)
	})

	t.Run("Retorna erro para data inválida em sale_date_to", func(t *testing.T) {
		invalidDate := "31-12-2024" // Formato inválido

		dto := SaleFilterDTO{
			SaleDateTo: &invalidDate,
			Limit:      10,
			Offset:     0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "sale_date_to")
		assert.Contains(t, err.Error(), "formato esperado: YYYY-MM-DD")
	})

	t.Run("Testa validação específica de created_from", func(t *testing.T) {
		testCases := []struct {
			name          string
			date          string
			shouldFail    bool
			errorContains string
		}{
			{"formato válido", "2024-01-01", false, ""},
			{"formato inválido (DD-MM-YYYY)", "01-01-2024", true, "created_from"},
			{"formato com barras", "2024/01/01", true, "created_from"},
			{"string vazia", "", false, ""}, // Não deve gerar erro, deve retornar nil
			{"data com hora", "2024-01-01T10:30:00", true, "formato esperado"},
			{"mês inválido", "2024-13-01", true, "created_from"},
			{"dia inválido", "2024-01-32", true, "created_from"},
			{"ano bissexto válido", "2024-02-29", false, ""},
			{"ano não bissexto inválido", "2023-02-29", true, "created_from"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				dto := SaleFilterDTO{
					CreatedFrom: &tc.date,
					Limit:       10,
					Offset:      0,
				}

				model, err := dto.ToModel()

				if tc.shouldFail {
					assert.Error(t, err)
					assert.Nil(t, model)
					if tc.errorContains != "" {
						assert.Contains(t, err.Error(), tc.errorContains)
					}
				} else {
					require.NoError(t, err)
					if tc.date == "" {
						// String vazia deve resultar em nil
						assert.Nil(t, model.CreatedFrom)
					} else {
						require.NotNil(t, model)
						require.NotNil(t, model.CreatedFrom)
						assert.Equal(t, tc.date, model.CreatedFrom.Format("2006-01-02"))
					}
				}
			})
		}
	})

	t.Run("Retorna erro para data inválida em created_to", func(t *testing.T) {
		invalidDate := "invalid-date-format"

		dto := SaleFilterDTO{
			CreatedTo: &invalidDate,
			Limit:     10,
			Offset:    0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "created_to")
		assert.Contains(t, err.Error(), "formato esperado: YYYY-MM-DD")
	})

	t.Run("Retorna erro para data inválida em updated_from", func(t *testing.T) {
		invalidDate := "2024/06/01" // Formato inválido

		dto := SaleFilterDTO{
			UpdatedFrom: &invalidDate,
			Limit:       10,
			Offset:      0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "updated_from")
		assert.Contains(t, err.Error(), "formato esperado: YYYY-MM-DD")
	})

	t.Run("Retorna erro para data inválida em updated_to", func(t *testing.T) {
		invalidDate := "01-06-2024" // Formato DD-MM-YYYY inválido

		dto := SaleFilterDTO{
			UpdatedTo: &invalidDate,
			Limit:     10,
			Offset:    0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "updated_to")
		assert.Contains(t, err.Error(), "formato esperado: YYYY-MM-DD")
	})

	t.Run("Retorna erro para valor inválido em max_total_amount", func(t *testing.T) {
		invalidValue := "not-a-number"

		dto := SaleFilterDTO{
			MaxTotalAmount: &invalidValue,
			Limit:          10,
			Offset:         0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "max_total_amount")
		assert.Contains(t, err.Error(), "valor numérico esperado")
	})

	t.Run("Retorna erro para valor inválido em min_items_amount", func(t *testing.T) {
		invalidValue := "abc"

		dto := SaleFilterDTO{
			MinItemsAmount: &invalidValue,
			Limit:          10,
			Offset:         0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "min_items_amount")
		assert.Contains(t, err.Error(), "valor numérico esperado")
	})

	t.Run("Retorna erro para valor inválido em max_items_amount", func(t *testing.T) {
		invalidValue := "xyz"

		dto := SaleFilterDTO{
			MaxItemsAmount: &invalidValue,
			Limit:          10,
			Offset:         0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "max_items_amount")
		assert.Contains(t, err.Error(), "valor numérico esperado")
	})

	t.Run("Retorna erro para valor inválido em min_items_amount", func(t *testing.T) {
		invalidValue := "abc"

		dto := SaleFilterDTO{
			MinItemsAmount: &invalidValue,
			Limit:          10,
			Offset:         0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "min_items_amount")
		assert.Contains(t, err.Error(), "valor numérico esperado")
	})

	t.Run("Retorna erro para valor inválido em max_items_amount", func(t *testing.T) {
		invalidValue := "xyz"

		dto := SaleFilterDTO{
			MaxItemsAmount: &invalidValue,
			Limit:          10,
			Offset:         0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "max_items_amount")
		assert.Contains(t, err.Error(), "valor numérico esperado")
	})

	t.Run("Retorna erro para valor inválido em min_sale_discount", func(t *testing.T) {
		invalidValue := "discount-invalid"

		dto := SaleFilterDTO{
			MinSaleDiscount: &invalidValue,
			Limit:           10,
			Offset:          0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "min_sale_discount")
		assert.Contains(t, err.Error(), "valor numérico esperado")
	})

	t.Run("Retorna erro para valor inválido em max_sale_discount", func(t *testing.T) {
		invalidValue := "invalid-discount"

		dto := SaleFilterDTO{
			MaxSaleDiscount: &invalidValue,
			Limit:           10,
			Offset:          0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "max_sale_discount")
		assert.Contains(t, err.Error(), "valor numérico esperado")
	})

	t.Run("Retorna erro para valor inválido em min_items_discount", func(t *testing.T) {
		invalidValue := "not-a-float"

		dto := SaleFilterDTO{
			MinItemsDiscount: &invalidValue,
			Limit:            10,
			Offset:           0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "min_items_discount")
		assert.Contains(t, err.Error(), "valor numérico esperado")
	})

	t.Run("Retorna erro para valor inválido em max_items_discount", func(t *testing.T) {
		invalidValue := "abc123"

		dto := SaleFilterDTO{
			MaxItemsDiscount: &invalidValue,
			Limit:            10,
			Offset:           0,
		}

		model, err := dto.ToModel()
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "max_items_discount")
		assert.Contains(t, err.Error(), "valor numérico esperado")
	})

}
