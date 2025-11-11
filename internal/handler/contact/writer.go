package handler

import (
	"errors"
	"net/http"

	dtoContact "github.com/WagaoCarvalho/backend_store_go/internal/dto/contact"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *contactHandler) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[ContactHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, nil)

	var contactDTO dtoContact.ContactDTO
	if err := utils.FromJSON(r.Body, &contactDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	contactModel := dtoContact.ToContactModel(contactDTO)

	createdContact, err := h.service.Create(ctx, contactModel)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrDBInvalidForeignKey):
			h.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{"erro": err.Error()})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		case errors.Is(err, errMsg.ErrInvalidData):
			h.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{"erro": err.Error()})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		case errors.Is(err, errMsg.ErrDuplicate):
			h.logger.Warn(ctx, ref+logger.LogErrDuplicate, map[string]any{"erro": err.Error()})
			utils.ErrorResponse(w, err, http.StatusConflict)
			return
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{"erro": err.Error()})
			utils.ErrorResponse(w, err, http.StatusNotFound)
			return
		default:
			h.logger.Error(ctx, err, ref+logger.LogCreateError, nil)
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	createdDTO := dtoContact.ToContactDTO(createdContact)

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"contact_id": createdDTO.ID,
	})

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Contato criado com sucesso",
		Data:    createdDTO,
	})
}

func (h *contactHandler) Update(w http.ResponseWriter, r *http.Request) {
	const ref = "[ContactHandler - Update] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogUpdateInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	var dto dtoContact.ContactDTO
	if err := utils.FromJSON(r.Body, &dto); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	contactModel := dtoContact.ToContactModel(dto)
	contactModel.ID = id

	if err := h.service.Update(ctx, contactModel); err != nil {
		var status int
		switch {
		case errors.Is(err, errMsg.ErrZeroID), errors.Is(err, errMsg.ErrInvalidData):
			status = http.StatusBadRequest
			h.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{"erro": err.Error()})
		case errors.Is(err, errMsg.ErrNotFound):
			status = http.StatusNotFound
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{"erro": err.Error()})
		case errors.Is(err, errMsg.ErrDuplicate):
			status = http.StatusConflict
			h.logger.Warn(ctx, ref+logger.LogErrDuplicate, map[string]any{"erro": err.Error()})
		default:
			status = http.StatusInternalServerError
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{"contact_id": id})
		}

		utils.ErrorResponse(w, err, status)
		return
	}

	updatedDTO := dtoContact.ToContactDTO(contactModel)
	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"contact_id": updatedDTO.ID,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Contato atualizado com sucesso",
		Data:    updatedDTO,
	})
}

func (h *contactHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[ContactHandler - Delete] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogDeleteInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		h.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"contact_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"contact_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}
