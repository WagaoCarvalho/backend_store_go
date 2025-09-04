package handler

import (
	"errors"
	"net/http"

	dtoAddress "github.com/WagaoCarvalho/backend_store_go/internal/dto/address"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/address"
)

type AddressHandler struct {
	service service.AddressService
	logger  *logger.LogAdapter
}

func NewAddressHandler(service service.AddressService, logger *logger.LogAdapter) *AddressHandler {
	return &AddressHandler{
		service: service,
		logger:  logger,
	}
}

func (h *AddressHandler) Create(w http.ResponseWriter, r *http.Request) {
	ref := "[AddressHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, nil)

	var addressDTO dtoAddress.AddressDTO
	if err := utils.FromJSON(r.Body, &addressDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// DTO → Model
	addressModel := dtoAddress.ToAddressModel(addressDTO)

	createdModel, err := h.service.Create(ctx, addressModel)
	if err != nil {
		if errors.Is(err, errMsg.ErrInvalidForeignKey) {
			h.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogCreateError, nil)
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Model → DTO
	createdDTO := dtoAddress.ToAddressDTO(createdModel)

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"address_id": createdDTO.ID,
	})

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Endereço criado com sucesso",
		Data:    createdDTO,
	})
}

func (h *AddressHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ref := "[addressHandler - GetByID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// service retorna model
	addressModel, err := h.service.GetByID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"address_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	// converter model -> DTO
	addressDTO := dtoAddress.ToAddressDTO(addressModel)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"address_id": addressDTO.ID,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereço encontrado",
		Data:    addressDTO,
	})
}

func (h *AddressHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	ref := "[addressHandler - GetByUserID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	addressModels, err := h.service.GetByUserID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{"user_id": id})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	// conversão model -> DTO
	addressDTOs := make([]*dtoAddress.AddressDTO, len(addressModels))
	for i, m := range addressModels {
		dto := dtoAddress.ToAddressDTO(m)
		addressDTOs[i] = &dto
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"user_id": id,
		"count":   len(addressDTOs),
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereços do usuário encontrados",
		Data:    addressDTOs,
	})
}

func (h *AddressHandler) GetByClientID(w http.ResponseWriter, r *http.Request) {
	ref := "[addressHandler - GetByClientID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	addressModels, err := h.service.GetByClientID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{"client_id": id})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	// conversão model -> DTO
	addressDTOs := make([]*dtoAddress.AddressDTO, len(addressModels))
	for i, m := range addressModels {
		dto := dtoAddress.ToAddressDTO(m)
		addressDTOs[i] = &dto
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"client_id": id,
		"count":     len(addressDTOs),
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereços do cliente encontrados",
		Data:    addressDTOs,
	})
}

func (h *AddressHandler) GetBySupplierID(w http.ResponseWriter, r *http.Request) {
	ref := "[addressHandler - GetBySupplierID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	addressModels, err := h.service.GetBySupplierID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{"supplier_id": id})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	// conversão model -> DTO
	addressDTOs := make([]*dtoAddress.AddressDTO, len(addressModels))
	for i, m := range addressModels {
		dto := dtoAddress.ToAddressDTO(m)
		addressDTOs[i] = &dto
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"supplier_id": id,
		"count":       len(addressDTOs),
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereços do fornecedor encontrados",
		Data:    addressDTOs,
	})
}

func (h *AddressHandler) Update(w http.ResponseWriter, r *http.Request) {
	const ref = "[AddressHandler - Update] "
	ctx := r.Context()

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	var dto dtoAddress.AddressDTO
	if err := utils.FromJSON(r.Body, &dto); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	dto.ID = &id
	addressModel := dtoAddress.ToAddressModel(dto)

	h.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"address_id": id,
	})

	// service agora recebe model
	if err := h.service.Update(ctx, addressModel); err != nil {
		if ve, ok := err.(*validators.ValidationError); ok {
			h.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
				"erro": ve.Error(),
				"id":   id,
			})
			utils.ErrorResponse(w, ve, http.StatusBadRequest)
			return
		}

		if errors.Is(err, errMsg.ErrID) {
			h.logger.Warn(ctx, ref+"Invalid ID", map[string]any{"id": id})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{"id": id})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{"id": id})

	// converter de volta para DTO na resposta
	updatedDTO := dtoAddress.ToAddressDTO(addressModel)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereço atualizado com sucesso",
		Data:    updatedDTO,
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
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(r.Context(), ref+logger.LogNotFound, map[string]any{
				"address_id": id,
				"erro":       err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusNotFound)
		case errors.Is(err, errMsg.ErrID):
			h.logger.Warn(r.Context(), ref+logger.LogValidateError, map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, errMsg.ErrID, http.StatusBadRequest)
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
