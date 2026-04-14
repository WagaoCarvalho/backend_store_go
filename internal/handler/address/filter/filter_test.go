package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	mockAddress "github.com/WagaoCarvalho/backend_store_go/infra/mock/address"
	filterDTO "github.com/WagaoCarvalho/backend_store_go/internal/dto/address/filter"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/address/address"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/address/filter"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddressHandler_Filter(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	createRequest := func(queryParams string) *http.Request {
		fullURL := "/addresses/filter"
		if queryParams != "" {
			fullURL += "?" + queryParams
		}
		return httptest.NewRequest(http.MethodGet, fullURL, nil)
	}

	t.Run("erro - parâmetro desconhecido", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		req := createRequest("invalid_param=value")
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "parâmetro desconhecido")
	})

	t.Run("erro - formato de data inválido", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		req := createRequest("created_from=invalid-date")
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "formato de data inválido")
	})

	t.Run("erro - is_active inválido", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		req := createRequest("is_active=invalid")
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "valor inválido para 'is_active'")
	})

	t.Run("erro - user_id inválido", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		req := createRequest("user_id=abc")
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "valor inválido para 'user_id'")
	})

	t.Run("erro - user_id menor ou igual a zero", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		req := createRequest("user_id=0")
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "deve ser maior que zero")
	})

	t.Run("erro - client_cpf_id inválido", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		req := createRequest("client_cpf_id=abc")
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "valor inválido para 'client_cpf_id'")
	})

	t.Run("erro - client_cpf_id menor ou igual a zero", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		req := createRequest("client_cpf_id=0")
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "valor inválido para 'client_cpf_id'")
		assert.Contains(t, resp.Message, "deve ser maior que zero")
	})

	t.Run("erro - client_cpf_id negativo", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		req := createRequest("client_cpf_id=-5")
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Contains(t, resp.Message, "deve ser maior que zero")
	})

	t.Run("erro - supplier_id inválido", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		req := createRequest("supplier_id=abc")
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "valor inválido para 'supplier_id'")
	})

	t.Run("erro - supplier_id menor ou igual a zero", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		req := createRequest("supplier_id=0")
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "valor inválido para 'supplier_id'")
		assert.Contains(t, resp.Message, "deve ser maior que zero")
	})

	t.Run("erro - supplier_id negativo", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		req := createRequest("supplier_id=-10")
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Contains(t, resp.Message, "deve ser maior que zero")
	})

	t.Run("erro - nenhum filtro de conteúdo fornecido", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		req := createRequest("page=1&limit=10&sort_by=city&sort_order=asc")
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "pelo menos um filtro de busca deve ser fornecido")
	})

	t.Run("erro - falha no serviço", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return(nil, errors.New("db error"))

		req := createRequest("city=" + url.QueryEscape("São Paulo"))
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - filtro inválido", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return(nil, errMsg.ErrInvalidFilter)

		req := createRequest("city=" + url.QueryEscape("São Paulo"))
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - retorna lista vazia", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return([]*model.Address{}, nil)

		req := createRequest("city=" + url.QueryEscape("CidadeInexistente"))
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Status)

		data, ok := resp.Data.(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, float64(0), data["total"])

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro com is_active true", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)
		now := time.Now()

		mockAddresses := []*model.Address{
			{
				ID:           1,
				UserID:       utils.Int64Ptr(100),
				Street:       "Rua Teste",
				StreetNumber: "123",
				Complement:   "Apto 101",
				City:         "São Paulo",
				State:        "SP",
				Country:      "Brasil",
				PostalCode:   "01234567",
				IsActive:     true,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		}

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.AddressFilter) bool {
				return f.IsActive != nil
			})).
			Return(mockAddresses, nil)

		req := createRequest("is_active=true")
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		data := resp.Data.(map[string]interface{})
		assert.Equal(t, float64(1), data["total"])
		items := data["items"].([]interface{})
		item := items[0].(map[string]interface{})
		assert.Equal(t, "Rua Teste", item["street"])
		assert.Equal(t, "São Paulo", item["city"])

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro com data válida", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		mockAddresses := []*model.Address{
			{
				ID:        1,
				Street:    "Rua Teste",
				City:      "São Paulo",
				State:     "SP",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.AddressFilter) bool {
				return f.CreatedFrom != nil
			})).
			Return(mockAddresses, nil)

		req := createRequest("created_from=2024-01-01")
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - múltiplos filtros", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		mockAddresses := []*model.Address{
			{
				ID:           1,
				Street:       "Rua Teste",
				StreetNumber: "123",
				City:         "São Paulo",
				State:        "SP",
				PostalCode:   "01234567",
				IsActive:     true,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
		}

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.AddressFilter) bool {
				return f.City == "São Paulo" &&
					f.State == "SP" &&
					f.PostalCode == "01234567" &&
					f.IsActive != nil
			})).
			Return(mockAddresses, nil)

		req := createRequest("city=" + url.QueryEscape("São Paulo") + "&state=SP&postal_code=01234567&is_active=true")
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro por user_id", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		mockAddresses := []*model.Address{
			{
				ID:        1,
				UserID:    utils.Int64Ptr(100),
				Street:    "Rua Teste",
				City:      "São Paulo",
				State:     "SP",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.AddressFilter) bool {
				return f.UserID != nil && *f.UserID == 100
			})).
			Return(mockAddresses, nil)

		req := createRequest("user_id=100")
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - com ordenação", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		mockAddresses := []*model.Address{
			{
				ID:        1,
				Street:    "Rua Teste",
				City:      "São Paulo",
				State:     "SP",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return(mockAddresses, nil)

		req := createRequest("city=" + url.QueryEscape("São Paulo") + "&sort_by=city&sort_order=desc")
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - range de datas", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		mockAddresses := []*model.Address{
			{
				ID:        1,
				Street:    "Rua Teste",
				City:      "São Paulo",
				State:     "SP",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return(mockAddresses, nil)

		req := createRequest("city=" + url.QueryEscape("São Paulo") + "&created_from=2024-01-01&created_to=2024-12-31&updated_from=2024-01-01&updated_to=2024-12-31")
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("erro - created_from inválido", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		req := createRequest("created_from=invalid-date")
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Contains(t, resp.Message, "formato de data inválido para 'created_from'")
	})

	t.Run("erro - created_to inválido", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		req := createRequest("created_to=31-12-2024")
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Contains(t, resp.Message, "formato de data inválido para 'created_to'")
	})

	t.Run("erro - updated_from inválido", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		req := createRequest("updated_from=2024/01/01")
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Contains(t, resp.Message, "formato de data inválido para 'updated_from'")
	})

	t.Run("erro - updated_to inválido", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		req := createRequest("updated_to=abc")
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Contains(t, resp.Message, "formato de data inválido para 'updated_to'")
	})

	t.Run("sucesso - verifica contagem completa de filtros aplicados", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		mockAddresses := []*model.Address{
			{
				ID:        1,
				Street:    "Rua Exemplo",
				City:      "Cidade Teste",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return(mockAddresses, nil)

		req := createRequest(
			"client_cpf_id=100&" +
				"supplier_id=200&" +
				"street=Rua+Principal&" +
				"street_number=123&" +
				"complement=Apto+45&" +
				"country=Brasil&" +
				"city=São+Paulo&" +
				"state=SP&" +
				"postal_code=01234567&" +
				"is_active=true&" +
				"created_from=2024-01-01&" +
				"created_to=2024-12-31&" +
				"updated_from=2024-01-01&" +
				"updated_to=2024-12-31&" +
				"user_id=50",
		)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		data := resp.Data.(map[string]interface{})
		filtersApplied := int(data["filters_applied"].(float64))
		assert.GreaterOrEqual(t, filtersApplied, 14, "Deveria contar todos os 14 filtros fornecidos")

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - contagem de filtros para campos de texto", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return([]*model.Address{}, nil)

		t.Run("apenas street", func(t *testing.T) {
			req := createRequest("street=Rua+Teste")
			rec := httptest.NewRecorder()
			handler.Filter(rec, req)

			var resp utils.DefaultResponse
			json.Unmarshal(rec.Body.Bytes(), &resp)
			data := resp.Data.(map[string]interface{})
			assert.Equal(t, float64(1), data["filters_applied"])
		})

		t.Run("street e street_number", func(t *testing.T) {
			req := createRequest("street=Rua+Teste&street_number=123")
			rec := httptest.NewRecorder()
			handler.Filter(rec, req)

			var resp utils.DefaultResponse
			json.Unmarshal(rec.Body.Bytes(), &resp)
			data := resp.Data.(map[string]interface{})
			assert.Equal(t, float64(2), data["filters_applied"])
		})

		t.Run("complement", func(t *testing.T) {
			req := createRequest("complement=Apto+101")
			rec := httptest.NewRecorder()
			handler.Filter(rec, req)

			var resp utils.DefaultResponse
			json.Unmarshal(rec.Body.Bytes(), &resp)
			data := resp.Data.(map[string]interface{})
			assert.Equal(t, float64(1), data["filters_applied"])
		})

		t.Run("country", func(t *testing.T) {
			req := createRequest("country=Brasil")
			rec := httptest.NewRecorder()
			handler.Filter(rec, req)

			var resp utils.DefaultResponse
			json.Unmarshal(rec.Body.Bytes(), &resp)
			data := resp.Data.(map[string]interface{})
			assert.Equal(t, float64(1), data["filters_applied"])
		})

		t.Run("todos campos de texto", func(t *testing.T) {
			req := createRequest(
				"street=Rua+X&" +
					"street_number=10&" +
					"complement=Casa&" +
					"country=Argentina",
			)
			rec := httptest.NewRecorder()
			handler.Filter(rec, req)

			var resp utils.DefaultResponse
			json.Unmarshal(rec.Body.Bytes(), &resp)
			data := resp.Data.(map[string]interface{})
			assert.Equal(t, float64(4), data["filters_applied"])
		})
	})

	t.Run("sucesso - contagem de filtros para IDs", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		handler := NewAddressFilterHandler(mockService, logger)

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return([]*model.Address{}, nil)

		t.Run("client_cpf_id", func(t *testing.T) {
			req := createRequest("client_cpf_id=123")
			rec := httptest.NewRecorder()
			handler.Filter(rec, req)

			var resp utils.DefaultResponse
			json.Unmarshal(rec.Body.Bytes(), &resp)
			data := resp.Data.(map[string]interface{})
			assert.Equal(t, float64(1), data["filters_applied"])
		})

		t.Run("supplier_id", func(t *testing.T) {
			req := createRequest("supplier_id=456")
			rec := httptest.NewRecorder()
			handler.Filter(rec, req)

			var resp utils.DefaultResponse
			json.Unmarshal(rec.Body.Bytes(), &resp)
			data := resp.Data.(map[string]interface{})
			assert.Equal(t, float64(1), data["filters_applied"])
		})

		t.Run("ambos os IDs", func(t *testing.T) {
			req := createRequest("client_cpf_id=123&supplier_id=456")
			rec := httptest.NewRecorder()
			handler.Filter(rec, req)

			var resp utils.DefaultResponse
			json.Unmarshal(rec.Body.Bytes(), &resp)
			data := resp.Data.(map[string]interface{})
			assert.Equal(t, float64(2), data["filters_applied"])
		})
	})
}

func TestAddressHandler_ParseTimeParam(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	handler := &addressFilterHandler{logger: logger}

	tests := []struct {
		name     string
		dateStr  string
		expected bool
	}{
		{
			name:     "data RFC3339",
			dateStr:  "2024-01-01T10:30:00Z",
			expected: true,
		},
		{
			name:     "data YYYY-MM-DD",
			dateStr:  "2024-01-01",
			expected: true,
		},
		{
			name:     "data com espaço",
			dateStr:  "2024-01-01 10:30:00",
			expected: true,
		},
		{
			name:     "data inválida",
			dateStr:  "invalid-date",
			expected: false,
		},
		{
			name:     "string vazia",
			dateStr:  "",
			expected: false,
		},
		{
			name:     "apenas espaços",
			dateStr:  "   ",
			expected: false,
		},
		{
			name:     "parâmetro não existe",
			dateStr:  "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := map[string][]string{}
			if tt.dateStr != "" {
				query["date"] = []string{tt.dateStr}
			}

			result := handler.parseTimeParam(query, "date")

			if tt.expected {
				assert.NotNil(t, result)
				assert.IsType(t, &time.Time{}, result)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestCountFiltersApplied(t *testing.T) {

	now := time.Now()

	dto := filterDTO.AddressFilterDTO{
		UserID:       utils.Int64Ptr(1),
		ClientCpfID:  utils.Int64Ptr(2),
		SupplierID:   utils.Int64Ptr(3),
		Street:       "Rua Teste",
		StreetNumber: "123",
		Complement:   "Apto 101",
		City:         "São Paulo",
		State:        "SP",
		Country:      "Brasil",
		PostalCode:   "01234567",
		IsActive:     utils.BoolPtr(true),
		CreatedFrom:  &now,
		CreatedTo:    &now,
		UpdatedFrom:  &now,
		UpdatedTo:    &now,
	}

	count := countFiltersApplied(dto)
	assert.Equal(t, 15, count)

	dto2 := filterDTO.AddressFilterDTO{}
	count2 := countFiltersApplied(dto2)
	assert.Equal(t, 0, count2)
}
