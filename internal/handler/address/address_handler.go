package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/address"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/address"
)

type AddressHandler struct {
	service service.AddressService
	logger  *logger.LoggerAdapter
}

func NewAddressHandler(service service.AddressService, logger *logger.LoggerAdapter) *AddressHandler {
	return &AddressHandler{
		service: service,
		logger:  logger,
	}
}

func (h *AddressHandler) Create(w http.ResponseWriter, r *http.Request) {
	ref := "[addressHandler - Create] "
	var address models.Address

	h.logger.Info(r.Context(), ref+logger.LogCreateInit, map[string]any{})

	if err := utils.FromJson(r.Body, &address); err != nil {
		h.logger.Warn(r.Context(), ref+logger.LogParseJsonError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	createdAddress, err := h.service.Create(r.Context(), &address)
	if err != nil {
		if errors.Is(err, repo.ErrInvalidForeignKey) {
			h.logger.Warn(r.Context(), ref+logger.LogForeignKeyViolation, map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		h.logger.Error(r.Context(), err, ref+logger.LogCreateError, map[string]any{})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(r.Context(), ref+logger.LogCreateSuccess, map[string]any{
		"address_id": createdAddress.ID,
	})

	utils.ToJson(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Endereço criado com sucesso",
		Data:    createdAddress,
	})
}

func (h *AddressHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ref := "[addressHandler - GetByID] "
	h.logger.Info(r.Context(), ref+logger.LogGetInit, map[string]any{})

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	address, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error(r.Context(), err, ref+logger.LogGetError, map[string]any{
			"address_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	h.logger.Info(r.Context(), ref+logger.LogGetSuccess, map[string]any{
		"address_id": address.ID,
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereço encontrado",
		Data:    address,
	})
}

func (h *AddressHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	ref := "[addressHandler - GetByUserID] "
	h.logger.Info(r.Context(), ref+logger.LogGetInit, map[string]any{})

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	addresses, err := h.service.GetByUserID(r.Context(), id)
	if err != nil {
		h.logger.Error(r.Context(), err, ref+logger.LogGetError, map[string]any{
			"user_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	h.logger.Info(r.Context(), ref+logger.LogGetSuccess, map[string]any{
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
	ref := "[addressHandler - GetByClientID] "
	h.logger.Info(r.Context(), ref+logger.LogGetInit, map[string]any{})

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	addresses, err := h.service.GetByClientID(r.Context(), id)
	if err != nil {
		h.logger.Error(r.Context(), err, ref+logger.LogGetError, map[string]any{
			"client_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	h.logger.Info(r.Context(), ref+logger.LogGetSuccess, map[string]any{
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
	ref := "[addressHandler - GetBySupplierID] "
	h.logger.Info(r.Context(), ref+logger.LogGetInit, map[string]any{})

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	addresses, err := h.service.GetBySupplierID(r.Context(), id)
	if err != nil {
		h.logger.Error(r.Context(), err, ref+logger.LogGetError, map[string]any{
			"supplier_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	h.logger.Info(r.Context(), ref+logger.LogGetSuccess, map[string]any{
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
	const ref = "[AddressHandler - Update] "

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	var address models.Address
	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		h.logger.Warn(r.Context(), ref+logger.LogParseJsonError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	address.ID = id

	h.logger.Info(r.Context(), ref+logger.LogUpdateInit, map[string]any{
		"address_id": address.ID,
	})

	if err := h.service.Update(r.Context(), &address); err != nil {
		if ve, ok := err.(*validators.ValidationError); ok {
			h.logger.Warn(r.Context(), ref+logger.LogValidateError, map[string]any{
				"erro": ve.Error(),
				"id":   id,
			})
			utils.ErrorResponse(w, ve, http.StatusBadRequest)
			return
		}

		h.logger.Error(r.Context(), err, ref+logger.LogUpdateError, map[string]any{
			"id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(r.Context(), ref+logger.LogUpdateSuccess, map[string]any{
		"id": id,
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereço atualizado com sucesso",
		Data:    address,
	})
}

func (h *AddressHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[AddressHandler - Delete] "

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), ref+logger.LogDeleteInit, map[string]any{
		"address_id": id,
		"path":       r.URL.Path,
	})

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, validators.ErrNotFound):
			h.logger.Warn(r.Context(), ref+logger.LogNotFound, map[string]any{
				"address_id": id,
				"erro":       err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusNotFound)
		case errors.Is(err, service.ErrAddressIDRequired):
			h.logger.Warn(r.Context(), ref+logger.LogValidateError, map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, errors.New("endereço ID é obrigatório"), http.StatusBadRequest)
		default:
			h.logger.Error(r.Context(), err, ref+logger.LogDeleteError, map[string]any{
				"address_id": id,
			})
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info(r.Context(), ref+logger.LogDeleteSuccess, map[string]any{
		"address_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}
