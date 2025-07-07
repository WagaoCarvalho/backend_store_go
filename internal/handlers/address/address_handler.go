package handlers

import (
	"errors"
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/addresses"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
)

type AddressHandler struct {
	service services.AddressService
	logger  *logger.LoggerAdapter
}

func NewAddressHandler(service services.AddressService, logger *logger.LoggerAdapter) *AddressHandler {
	return &AddressHandler{
		service: service,
		logger:  logger,
	}
}

func (h *AddressHandler) Create(w http.ResponseWriter, r *http.Request) {
	var address models.Address

	h.logger.Info(r.Context(), "[AddressHandler] - "+logger.LogCreateInit, map[string]interface{}{})

	if err := utils.FromJson(r.Body, &address); err != nil {
		h.logger.Warn(r.Context(), "[AddressHandler] - "+logger.LogParseJsonError, map[string]interface{}{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	createdAddress, err := h.service.Create(r.Context(), &address)
	if err != nil {
		if errors.Is(err, repositories.ErrInvalidForeignKey) {
			h.logger.Warn(r.Context(), "[AddressHandler] - "+logger.LogForeignKeyViolation, map[string]interface{}{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		h.logger.Error(r.Context(), err, "[AddressHandler] - "+logger.LogCreateError, map[string]interface{}{})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(r.Context(), "[AddressHandler] - "+logger.LogCreateSuccess, map[string]interface{}{
		"address_id": createdAddress.ID,
	})

	utils.ToJson(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Endereço criado com sucesso",
		Data:    createdAddress,
	})
}

func (h *AddressHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	h.logger.Info(r.Context(), "[AddressHandler] - Iniciando busca de endereço por ID", map[string]interface{}{
		"path": r.URL.Path,
	})

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), "[AddressHandler] - ID inválido recebido", map[string]interface{}{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	address, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error(r.Context(), err, "[AddressHandler] - Erro ao buscar endereço", map[string]interface{}{
			"address_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	h.logger.Info(r.Context(), "[AddressHandler] - Endereço encontrado com sucesso", map[string]interface{}{
		"address_id": address.ID,
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereço encontrado",
		Data:    address,
	})
}

func (h *AddressHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	h.logger.Info(r.Context(), "[AddressHandler] - Iniciando busca de endereços por UserID", map[string]interface{}{
		"path": r.URL.Path,
	})

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), "[AddressHandler] - ID inválido para busca por UserID", map[string]interface{}{
			"error": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), "[AddressHandler] - Iniciando busca de endereços por UserID", map[string]interface{}{
		"user_id": id,
	})

	addresses, err := h.service.GetByUserID(r.Context(), id)
	if err != nil {
		h.logger.Error(r.Context(), err, "[AddressHandler] - Erro ao buscar endereços por UserID", map[string]interface{}{
			"user_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	h.logger.Info(r.Context(), "[AddressHandler] - Endereços do usuário encontrados com sucesso", map[string]interface{}{
		"user_id": id,
		"count":   len(addresses),
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereços do usuário encontrados",
		Data:    addresses,
	})
}

func (h *AddressHandler) GetByClientID(w http.ResponseWriter, r *http.Request) {
	h.logger.Info(r.Context(), "[AddressHandler] - Iniciando busca de endereços por ClientID", map[string]interface{}{
		"path": r.URL.Path,
	})

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), "[AddressHandler] - ID do cliente inválido", map[string]interface{}{
			"error": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), "[AddressHandler] - Iniciando busca de endereços por ClientID", map[string]interface{}{
		"client_id": id,
	})

	addresses, err := h.service.GetByClientID(r.Context(), id)
	if err != nil {
		h.logger.Error(r.Context(), err, "[AddressHandler] - Erro ao buscar endereços por ClientID", map[string]interface{}{
			"client_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	h.logger.Info(r.Context(), "[AddressHandler] - Endereços do cliente encontrados com sucesso", map[string]interface{}{
		"client_id": id,
		"count":     len(addresses),
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereços do cliente encontrados",
		Data:    addresses,
	})
}

func (h *AddressHandler) GetBySupplierID(w http.ResponseWriter, r *http.Request) {
	h.logger.Info(r.Context(), "[AddressHandler] - Iniciando busca de endereços por SupplierID", map[string]interface{}{
		"path": r.URL.Path,
	})

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), "[AddressHandler] - ID inválido ao buscar endereço por SupplierID", map[string]interface{}{
			"error": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), "[AddressHandler] - Buscando endereços por SupplierID", map[string]interface{}{
		"supplier_id": id,
	})

	addresses, err := h.service.GetBySupplierID(r.Context(), id)
	if err != nil {
		h.logger.Warn(r.Context(), "[AddressHandler] - Erro ao buscar endereços por SupplierID", map[string]interface{}{
			"error":       err.Error(),
			"supplier_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	h.logger.Info(r.Context(), "[AddressHandler] - Endereços do fornecedor encontrados", map[string]interface{}{
		"supplier_id": id,
		"count":       len(addresses),
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereços do fornecedor encontrados",
		Data:    addresses,
	})
}

func (h *AddressHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), "[AddressHandler] - ID inválido para update", map[string]interface{}{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	var address models.Address
	if err := utils.FromJson(r.Body, &address); err != nil {
		h.logger.Warn(r.Context(), "[AddressHandler] - JSON inválido no update", map[string]interface{}{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	address.ID = id

	h.logger.Info(r.Context(), "[AddressHandler] - Iniciando atualização de endereço", map[string]interface{}{
		"path":       r.URL.Path,
		"address_id": address.ID,
	})

	if err := h.service.Update(r.Context(), &address); err != nil {
		if ve, ok := err.(*utils.ValidationError); ok {
			h.logger.Warn(r.Context(), "[AddressHandler] - Falha na validação ao atualizar endereço", map[string]interface{}{
				"erro": ve.Error(),
				"id":   id,
			})
			utils.ErrorResponse(w, ve, http.StatusBadRequest)
			return
		}

		h.logger.Error(r.Context(), err, "[AddressHandler] - Erro ao atualizar endereço", map[string]interface{}{
			"id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(r.Context(), "[AddressHandler] - Endereço atualizado com sucesso", map[string]interface{}{
		"id": id,
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereço atualizado com sucesso",
		Data:    nil,
	})
}

func (h *AddressHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), "[AddressHandler] - ID inválido para exclusão", map[string]interface{}{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, errors.New("ID inválido (esperado número inteiro)"), http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), "[AddressHandler] - Iniciando exclusão de endereço", map[string]interface{}{
		"address_id": id,
	})

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrNotFound):
			h.logger.Warn(r.Context(), "[AddressHandler] - Endereço não encontrado para exclusão", map[string]interface{}{
				"address_id": id,
				"erro":       err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusNotFound)
		case errors.Is(err, services.ErrAddressIDRequired):
			h.logger.Warn(r.Context(), "[AddressHandler] - ID do endereço obrigatório para exclusão", nil)
			utils.ErrorResponse(w, errors.New("endereço ID é obrigatório"), http.StatusBadRequest)
		default:
			h.logger.Error(r.Context(), err, "[AddressHandler] - Erro ao deletar endereço", map[string]interface{}{
				"address_id": id,
			})
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info(r.Context(), "[AddressHandler] - Endereço excluído com sucesso", map[string]interface{}{
		"address_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}
