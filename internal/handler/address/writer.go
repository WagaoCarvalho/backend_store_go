package handler

import (
	"errors"
	"net/http"

	dtoAddress "github.com/WagaoCarvalho/backend_store_go/internal/dto/address"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

func (h *addressHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	h.logger.Info(ctx, ref+"[Create] "+logger.LogCreateInit, nil)

	var addressDTO dtoAddress.AddressDTO
	if err := utils.FromJSON(r.Body, &addressDTO); err != nil {
		h.logger.Warn(ctx, ref+"[Create] "+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, errors.New("JSON inválido"), http.StatusBadRequest)
		return
	}

	addressModel := dtoAddress.ToAddressModel(addressDTO)

	createdModel, err := h.service.Create(ctx, addressModel)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrDBInvalidForeignKey):
			h.logger.Warn(ctx, ref+"[Create] "+logger.LogForeignKeyViolation, map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, errors.New("chave estrangeira inválida"), http.StatusBadRequest)

		case errors.Is(err, errMsg.ErrDuplicate):
			h.logger.Warn(ctx, ref+"[Create] Endereço duplicado", map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, errors.New("endereço já existente"), http.StatusConflict)

		default:
			h.logger.Error(ctx, err, ref+"[Create] "+logger.LogCreateError, nil)
			utils.ErrorResponse(w, errors.New("erro interno"), http.StatusInternalServerError)
		}
		return
	}

	createdDTO := dtoAddress.ToAddressDTO(createdModel)

	h.logger.Info(ctx, ref+"[Create] "+logger.LogCreateSuccess, map[string]any{
		"address_id": createdDTO.ID,
	})

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Endereço criado com sucesso",
		Data:    createdDTO,
	})
}

func (h *addressHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+"[Update] "+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, errors.New("ID inválido"), http.StatusBadRequest)
		return
	}

	var dto dtoAddress.AddressDTO
	if err := utils.FromJSON(r.Body, &dto); err != nil {
		h.logger.Warn(ctx, ref+"[Update] "+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, errors.New("JSON inválido"), http.StatusBadRequest)
		return
	}

	dto.ID = &id
	addressModel := dtoAddress.ToAddressModel(dto)

	h.logger.Info(ctx, ref+"[Update] "+logger.LogUpdateInit, map[string]any{
		"address_id": id,
	})

	if err := h.service.Update(ctx, addressModel); err != nil {
		if ve, ok := err.(*validators.ValidationError); ok {
			h.logger.Warn(ctx, ref+"[Update] "+logger.LogValidateError, map[string]any{
				"erro": ve.Error(),
				"id":   id,
			})
			utils.ErrorResponse(w, errors.New("dados inválidos"), http.StatusBadRequest)
			return
		}

		switch {
		case errors.Is(err, errMsg.ErrZeroID):
			h.logger.Warn(ctx, ref+"[Update] Invalid ID", map[string]any{"id": id})
			utils.ErrorResponse(w, errors.New("ID inválido"), http.StatusBadRequest)

		case errors.Is(err, errMsg.ErrDBInvalidForeignKey):
			h.logger.Warn(ctx, ref+"[Update] Invalid Foreign Key", map[string]any{
				"erro": err.Error(),
				"id":   id,
			})
			utils.ErrorResponse(w, errors.New("chave estrangeira inválida"), http.StatusBadRequest)

		default:
			h.logger.Error(ctx, err, ref+"[Update] "+logger.LogUpdateError, map[string]any{"id": id})
			utils.ErrorResponse(w, errors.New("erro interno"), http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info(ctx, ref+"[Update] "+logger.LogUpdateSuccess, map[string]any{
		"id": id,
	})

	updatedDTO := dtoAddress.ToAddressDTO(addressModel)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereço atualizado com sucesso",
		Data:    updatedDTO,
	})
}

func (h *addressHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+"[Delete] "+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, errors.New("ID inválido"), http.StatusBadRequest)
		return
	}

	h.logger.Info(ctx, ref+"[Delete] "+logger.LogDeleteInit, map[string]any{
		"address_id": id,
		"path":       r.URL.Path,
	})

	if err := h.service.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+"[Delete] "+logger.LogNotFound, map[string]any{
				"address_id": id,
			})
			utils.ErrorResponse(w, errors.New("endereço não encontrado"), http.StatusNotFound)

		case errors.Is(err, errMsg.ErrZeroID):
			h.logger.Warn(ctx, ref+"[Delete] "+logger.LogValidateError, map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, errors.New("ID inválido"), http.StatusBadRequest)

		default:
			h.logger.Error(ctx, err, ref+"[Delete] "+logger.LogDeleteError, map[string]any{
				"address_id": id,
			})
			utils.ErrorResponse(w, errors.New("erro interno"), http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info(ctx, ref+"[Delete] "+logger.LogDeleteSuccess, map[string]any{
		"address_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}
