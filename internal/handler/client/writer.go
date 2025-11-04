package handler

import (
	"errors"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/client/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *Client) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[ClientHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, nil)

	var clientDTO dto.ClientDTO
	if err := utils.FromJSON(r.Body, &clientDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	clientModel := dto.ToClientModel(clientDTO)

	createdModel, err := h.service.Create(ctx, clientModel)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrDBInvalidForeignKey):
			h.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, errMsg.ErrDBInvalidForeignKey, http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrDuplicate):
			h.logger.Warn(ctx, ref+"Cliente duplicado", map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, errMsg.ErrDuplicate, http.StatusConflict)
			return

		default:
			h.logger.Error(ctx, err, ref+logger.LogCreateError, nil)
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	createdDTO := dto.ToClientDTO(createdModel)

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"client_id": createdDTO.ID,
	})

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Cliente criado com sucesso",
		Data:    createdDTO,
	})
}

func (h *Client) Update(w http.ResponseWriter, r *http.Request) {
	const ref = "[ClientHandler - Update] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogUpdateInit, nil)

	uid, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	var clientDTO dto.ClientDTO
	if err := utils.FromJSON(r.Body, &clientDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	clientModel := dto.ToClientModel(clientDTO)
	clientModel.ID = uid

	err = h.service.Update(ctx, clientModel)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrInvalidData),
			errors.Is(err, errMsg.ErrZeroID):
			h.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
				"client_id": uid,
				"erro":      err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrDuplicate):
			h.logger.Warn(ctx, ref+"Cliente duplicado", map[string]any{
				"client_id": uid,
				"erro":      err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusConflict)
			return

		case errors.Is(err, errMsg.ErrVersionConflict):
			h.logger.Warn(ctx, ref+logger.LogUpdateVersionConflict, map[string]any{
				"client_id": uid,
			})
			utils.ErrorResponse(w, err, http.StatusConflict)
			return

		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"client_id": uid,
			})
			utils.ErrorResponse(w, err, http.StatusNotFound)
			return

		default:
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"client_id": uid,
			})
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	updatedDTO := dto.ToClientDTO(clientModel)

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"client_id": uid,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Cliente atualizado com sucesso",
		Data:    updatedDTO,
	})
}

func (h *Client) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[clientHandler - Delete] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{})

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	err = h.service.Delete(ctx, id)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"client_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"client_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}
