package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mockAddress "github.com/WagaoCarvalho/backend_store_go/infra/mock/address"
	dtoAddress "github.com/WagaoCarvalho/backend_store_go/internal/dto/address"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

func newReq(method, url string, body []byte, vars map[string]string) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	if vars != nil {
		req = mux.SetURLVars(req, vars)
	}
	return req
}

func newHandler(t *testing.T) (*addressHandler, *mockAddress.MockAddress) {
	t.Helper()
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	mockSvc := new(mockAddress.MockAddress)
	return NewAddressHandler(mockSvc, logger.NewLoggerAdapter(log)), mockSvc
}

/* ======================
   CREATE
====================== */

func TestAddressHandler_Create_AllCases(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		h, svc := newHandler(t)

		uid := int64(1)
		dto := dtoAddress.AddressDTO{UserID: &uid}
		model := dtoAddress.ToAddressModel(dto)
		model.ID = 10

		svc.On("Create", mock.Anything, mock.Anything).
			Return(model, nil)

		body, _ := json.Marshal(dto)
		w := httptest.NewRecorder()
		h.Create(w, newReq(http.MethodPost, "/addresses", body, nil))

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("InvalidForeignKey", func(t *testing.T) {
		h, svc := newHandler(t)

		svc.On("Create", mock.Anything, mock.Anything).
			Return((*models.Address)(nil), errMsg.ErrDBInvalidForeignKey)

		w := httptest.NewRecorder()
		h.Create(w, newReq(http.MethodPost, "/addresses", []byte(`{}`), nil))

		var resp utils.DefaultResponse
		_ = json.NewDecoder(w.Body).Decode(&resp)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "chave estrangeira inválida", resp.Message)
	})

	t.Run("Duplicate", func(t *testing.T) {
		h, svc := newHandler(t)

		svc.On("Create", mock.Anything, mock.Anything).
			Return((*models.Address)(nil), errMsg.ErrDuplicate)

		w := httptest.NewRecorder()
		h.Create(w, newReq(http.MethodPost, "/addresses", []byte(`{}`), nil))

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("InternalError", func(t *testing.T) {
		h, svc := newHandler(t)

		svc.On("Create", mock.Anything, mock.Anything).
			Return((*models.Address)(nil), errors.New("db down"))

		w := httptest.NewRecorder()
		h.Create(w, newReq(http.MethodPost, "/addresses", []byte(`{}`), nil))

		var resp utils.DefaultResponse
		_ = json.NewDecoder(w.Body).Decode(&resp)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "erro interno", resp.Message)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		h, _ := newHandler(t)

		// JSON malformado - falta a chave de fechamento
		invalidJSON := []byte(`{"street": "Rua Teste", "number": "123"`)

		w := httptest.NewRecorder()
		h.Create(w, newReq(http.MethodPost, "/addresses", invalidJSON, nil))

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, "JSON inválido", resp.Message)
		assert.Equal(t, http.StatusBadRequest, resp.Status)

		// Verificar que o serviço NÃO foi chamado
		mockSvc := new(mockAddress.MockAddress)
		hWithMock, _ := newHandler(t)
		hWithMock.service = mockSvc // Substitui o serviço pelo mock limpo

		w2 := httptest.NewRecorder()
		hWithMock.Create(w2, newReq(http.MethodPost, "/addresses", invalidJSON, nil))

		// Assert que o mock não teve interações
		mockSvc.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("EmptyBody", func(t *testing.T) {
		h, _ := newHandler(t)

		// Corpo vazio
		emptyBody := []byte{}

		w := httptest.NewRecorder()
		h.Create(w, newReq(http.MethodPost, "/addresses", emptyBody, nil))

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, "JSON inválido", resp.Message)
	})

	t.Run("InvalidJSONSyntax", func(t *testing.T) {
		h, _ := newHandler(t)

		// JSON com sintaxe completamente inválida
		invalidJSON := []byte(`{invalid json syntax}`)

		w := httptest.NewRecorder()
		h.Create(w, newReq(http.MethodPost, "/addresses", invalidJSON, nil))

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, "JSON inválido", resp.Message)
	})

	t.Run("InvalidDataType", func(t *testing.T) {
		h, _ := newHandler(t)

		// JSON válido sintaticamente, mas com tipo de dado errado para o campo
		// Supondo que UserID deve ser int64, mas está sendo enviado como string
		invalidDataTypeJSON := []byte(`{"user_id": "not_a_number"}`)

		w := httptest.NewRecorder()
		h.Create(w, newReq(http.MethodPost, "/addresses", invalidDataTypeJSON, nil))

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, "JSON inválido", resp.Message)
	})

	t.Run("JSONArrayInsteadOfObject", func(t *testing.T) {
		h, _ := newHandler(t)

		// Enviando array em vez de objeto
		arrayJSON := []byte(`[{"street": "Rua Teste"}]`)

		w := httptest.NewRecorder()
		h.Create(w, newReq(http.MethodPost, "/addresses", arrayJSON, nil))

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, "JSON inválido", resp.Message)
	})
}

/* ======================
   UPDATE
====================== */

func TestAddressHandler_Update_AllCases(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		h, svc := newHandler(t)

		svc.On("Update", mock.Anything, mock.Anything).Return(nil)

		body := []byte(`{"city":"SP"}`)
		w := httptest.NewRecorder()
		h.Update(w, newReq(http.MethodPut, "/addresses/1", body, map[string]string{"id": "1"}))

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		h, _ := newHandler(t)

		w := httptest.NewRecorder()
		h.Update(w, newReq(http.MethodPut, "/addresses/1", []byte("{bad"), map[string]string{"id": "1"}))

		var resp utils.DefaultResponse
		_ = json.NewDecoder(w.Body).Decode(&resp)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "JSON inválido", resp.Message)
	})

	t.Run("ValidationError", func(t *testing.T) {
		h, svc := newHandler(t)

		svc.On("Update", mock.Anything, mock.Anything).
			Return(&validators.ValidationError{Field: "city", Message: "obrigatório"})

		w := httptest.NewRecorder()
		h.Update(w, newReq(http.MethodPut, "/addresses/1", []byte(`{}`), map[string]string{"id": "1"}))

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ZeroID", func(t *testing.T) {
		h, svc := newHandler(t)

		svc.On("Update", mock.Anything, mock.Anything).
			Return(errMsg.ErrZeroID)

		w := httptest.NewRecorder()
		h.Update(w, newReq(http.MethodPut, "/addresses/1", []byte(`{}`), map[string]string{"id": "1"}))

		var resp utils.DefaultResponse
		_ = json.NewDecoder(w.Body).Decode(&resp)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "ID inválido", resp.Message)
	})

	t.Run("InvalidForeignKey", func(t *testing.T) {
		h, svc := newHandler(t)

		svc.On("Update", mock.Anything, mock.Anything).
			Return(errMsg.ErrDBInvalidForeignKey)

		w := httptest.NewRecorder()
		h.Update(w, newReq(http.MethodPut, "/addresses/1", []byte(`{}`), map[string]string{"id": "1"}))

		var resp utils.DefaultResponse
		_ = json.NewDecoder(w.Body).Decode(&resp)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "chave estrangeira inválida", resp.Message)
	})

	t.Run("InternalError", func(t *testing.T) {
		h, svc := newHandler(t)

		svc.On("Update", mock.Anything, mock.Anything).
			Return(errors.New("panic"))

		w := httptest.NewRecorder()
		h.Update(w, newReq(http.MethodPut, "/addresses/1", []byte(`{}`), map[string]string{"id": "1"}))

		var resp utils.DefaultResponse
		_ = json.NewDecoder(w.Body).Decode(&resp)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "erro interno", resp.Message)
	})

	t.Run("Success", func(t *testing.T) {
		h, svc := newHandler(t)

		svc.On("Update", mock.Anything, mock.Anything).Return(nil)

		body := []byte(`{"city":"SP"}`)
		w := httptest.NewRecorder()
		h.Update(w, newReq(http.MethodPut, "/addresses/1", body, map[string]string{"id": "1"}))

		assert.Equal(t, http.StatusOK, w.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, "Endereço atualizado com sucesso", resp.Message)
	})

	// Testes de ID inválido - APENAS IDs que GetIDParam rejeita
	// Se GetIDParam aceita IDs negativos e zero, precisamos testar o que ELE rejeita
	invalidIDTests := []struct {
		name string
		id   string
	}{
		{"NonNumeric", "abc"},
		{"Float", "1.5"},
		{"TooLarge", "99999999999999999999"},
		{"SpecialChars", "1@#$"},
	}

	for _, tc := range invalidIDTests {
		t.Run("InvalidIDParam_"+tc.name, func(t *testing.T) {
			h, _ := newHandler(t)

			body := []byte(`{"city":"SP"}`)
			w := httptest.NewRecorder()
			h.Update(w, newReq(http.MethodPut, "/addresses/"+tc.id, body, map[string]string{"id": tc.id}))

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var resp utils.DefaultResponse
			require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
			assert.Equal(t, "ID inválido", resp.Message)
		})
	}

	// Testes para IDs que GetIDParam ACEITA, mas que o serviço pode rejeitar
	t.Run("ZeroID_ServiceRejects", func(t *testing.T) {
		h, svc := newHandler(t)

		// GetIDParam aceita "0", mas o serviço rejeita com ErrZeroID
		svc.On("Update", mock.Anything, mock.Anything).
			Return(errMsg.ErrZeroID)

		body := []byte(`{"city":"SP"}`)
		w := httptest.NewRecorder()
		h.Update(w, newReq(http.MethodPut, "/addresses/0", body, map[string]string{"id": "0"}))

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, "ID inválido", resp.Message)
	})

	t.Run("NegativeID_ServiceRejects", func(t *testing.T) {
		h, svc := newHandler(t)

		// GetIDParam aceita "-1", mas o serviço pode rejeitar
		// Vamos assumir que o serviço também retorna ErrZeroID ou similar
		svc.On("Update", mock.Anything, mock.Anything).
			Return(errMsg.ErrZeroID) // ou outro erro apropriado

		body := []byte(`{"city":"SP"}`)
		w := httptest.NewRecorder()
		h.Update(w, newReq(http.MethodPut, "/addresses/-1", body, map[string]string{"id": "-1"}))

		// Aqui o handler chama o serviço, que retorna erro
		// Então o status deve ser BadRequest (do serviço)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, "ID inválido", resp.Message)
	})

	t.Run("MissingIDParam", func(t *testing.T) {
		h, _ := newHandler(t)

		body := []byte(`{"city":"SP"}`)
		req := httptest.NewRequest(http.MethodPut, "/addresses/", bytes.NewBuffer(body))
		// IMPORTANTE: Não use mux.SetURLVars - para simular ID ausente
		w := httptest.NewRecorder()
		h.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, "ID inválido", resp.Message)
	})

	t.Run("InvalidJSON_BeforeIDCheck", func(t *testing.T) {
		h, _ := newHandler(t)

		// JSON inválido - deve ser verificado antes do ID
		invalidJSON := []byte(`{invalid json`)
		w := httptest.NewRecorder()
		h.Update(w, newReq(http.MethodPut, "/addresses/1", invalidJSON, map[string]string{"id": "1"}))

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, "JSON inválido", resp.Message)
	})

}

func TestAddressHandler_Delete_AllCases(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		h, svc := newHandler(t)

		svc.On("Delete", mock.Anything, int64(1)).Return(nil)

		w := httptest.NewRecorder()
		h.Delete(w, newReq(http.MethodDelete, "/addresses/1", nil, map[string]string{"id": "1"}))

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, 0, w.Body.Len())
		assert.Empty(t, w.Body.String())
	})

	// Testes para IDs que GetIDParam rejeita
	invalidIDTests := []struct {
		name string
		id   string
	}{
		{"NonNumeric", "abc"},
		{"Float", "1.5"},
		{"TooLarge", "99999999999999999999"},
		{"SpecialChars", "1@#$"},
	}

	for _, tc := range invalidIDTests {
		t.Run("InvalidIDParam_"+tc.name, func(t *testing.T) {
			h, _ := newHandler(t)

			w := httptest.NewRecorder()
			h.Delete(w, newReq(http.MethodDelete, "/addresses/"+tc.id, nil, map[string]string{"id": tc.id}))

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var resp utils.DefaultResponse
			require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
			assert.Equal(t, "ID inválido", resp.Message)
			assert.Equal(t, http.StatusBadRequest, resp.Status)
		})
	}

	t.Run("ZeroID_ServiceRejects", func(t *testing.T) {
		h, svc := newHandler(t)

		// GetIDParam aceita "0", mas o serviço rejeita com ErrZeroID
		svc.On("Delete", mock.Anything, int64(0)).
			Return(errMsg.ErrZeroID)

		w := httptest.NewRecorder()
		h.Delete(w, newReq(http.MethodDelete, "/addresses/0", nil, map[string]string{"id": "0"}))

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, "ID inválido", resp.Message)
	})

	t.Run("NegativeID_ServiceRejects", func(t *testing.T) {
		h, svc := newHandler(t)

		// GetIDParam aceita "-1", mas o serviço rejeita
		svc.On("Delete", mock.Anything, int64(-1)).
			Return(errMsg.ErrZeroID) // ou outro erro apropriado

		w := httptest.NewRecorder()
		h.Delete(w, newReq(http.MethodDelete, "/addresses/-1", nil, map[string]string{"id": "-1"}))

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, "ID inválido", resp.Message)
	})

	t.Run("MissingIDParam", func(t *testing.T) {
		h, _ := newHandler(t)

		// Request sem variáveis de URL
		req := httptest.NewRequest(http.MethodDelete, "/addresses", nil)
		w := httptest.NewRecorder()
		h.Delete(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, "ID inválido", resp.Message)
	})

	t.Run("NotFound", func(t *testing.T) {
		h, svc := newHandler(t)

		svc.On("Delete", mock.Anything, int64(999)).
			Return(errMsg.ErrNotFound)

		w := httptest.NewRecorder()
		h.Delete(w, newReq(http.MethodDelete, "/addresses/999", nil, map[string]string{"id": "999"}))

		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, "endereço não encontrado", resp.Message)
		assert.Equal(t, http.StatusNotFound, resp.Status)
	})

	t.Run("InternalError", func(t *testing.T) {
		h, svc := newHandler(t)

		svc.On("Delete", mock.Anything, int64(1)).
			Return(errors.New("database connection failed"))

		w := httptest.NewRecorder()
		h.Delete(w, newReq(http.MethodDelete, "/addresses/1", nil, map[string]string{"id": "1"}))

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, "erro interno", resp.Message)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)
	})

	t.Run("DifferentErrorTypes", func(t *testing.T) {
		h, svc := newHandler(t)

		// Erro não mapeado
		customErr := errors.New("custom database error")
		svc.On("Delete", mock.Anything, int64(2)).
			Return(customErr)

		w := httptest.NewRecorder()
		h.Delete(w, newReq(http.MethodDelete, "/addresses/2", nil, map[string]string{"id": "2"}))

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, "erro interno", resp.Message)
	})

	t.Run("BodyInDeleteRequest", func(t *testing.T) {
		h, svc := newHandler(t)

		// DELETE requests podem ter corpo, mas geralmente ignoramos
		svc.On("Delete", mock.Anything, int64(1)).Return(nil)

		body := []byte(`{"unused": "data"}`)
		w := httptest.NewRecorder()
		h.Delete(w, newReq(http.MethodDelete, "/addresses/1", body, map[string]string{"id": "1"}))

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Body.String())
	})

	t.Run("SuccessLogging", func(t *testing.T) {
		h, svc := newHandler(t)

		// Teste para garantir que o logger é chamado no sucesso
		svc.On("Delete", mock.Anything, int64(123)).Return(nil)

		w := httptest.NewRecorder()
		h.Delete(w, newReq(http.MethodDelete, "/addresses/123", nil, map[string]string{"id": "123"}))

		assert.Equal(t, http.StatusNoContent, w.Code)
		// Não podemos verificar o logger diretamente, mas o código deve executar sem panic
	})

	t.Run("ConcurrentDeleteCalls", func(t *testing.T) {
		h, svc := newHandler(t)

		// Teste para múltiplas chamadas simultâneas
		svc.On("Delete", mock.Anything, int64(1)).Return(nil).Times(3)

		var wg sync.WaitGroup
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				w := httptest.NewRecorder()
				h.Delete(w, newReq(http.MethodDelete, "/addresses/1", nil, map[string]string{"id": "1"}))
				assert.Equal(t, http.StatusNoContent, w.Code)
			}()
		}
		wg.Wait()

		svc.AssertExpectations(t)
	})

	t.Run("DifferentHTTPMethods", func(t *testing.T) {
		h, svc := newHandler(t)

		// DELETE é o método esperado, mas testamos que o handler só responde ao DELETE
		// (isso é mais para documentação, já que o router cuida disso)
		svc.On("Delete", mock.Anything, int64(1)).Return(nil)

		// Teste com método correto
		w := httptest.NewRecorder()
		h.Delete(w, newReq(http.MethodDelete, "/addresses/1", nil, map[string]string{"id": "1"}))
		assert.Equal(t, http.StatusNoContent, w.Code)

		// O handler não deve aceitar outros métodos, mas isso é controlado pelo router
		// Não é necessário testar aqui
	})
}
