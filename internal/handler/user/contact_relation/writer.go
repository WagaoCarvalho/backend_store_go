package handler

import (
	"errors"
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/contact_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *UserContactRelation) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserContactRelationHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{})

	var requestData dto.UserContactRelationDTO
	if err := utils.FromJSON(r.Body, &requestData); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, fmt.Errorf("erro ao decodificar JSON"), http.StatusBadRequest)
		return
	}

	modelRelation := dto.ToContactRelationModel(requestData)

	// Validação simples de IDs antes de chamar o service
	if modelRelation == nil || modelRelation.UserID <= 0 || modelRelation.ContactID <= 0 {
		h.logger.Warn(ctx, ref+"modelo nulo ou ID inválido", map[string]any{})
		utils.ErrorResponse(w, fmt.Errorf("modelo nulo ou ID inválido"), http.StatusBadRequest)
		return
	}

	created, err := h.service.Create(ctx, modelRelation)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrDBInvalidForeignKey):
			h.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"user_id":    modelRelation.UserID,
				"contact_id": modelRelation.ContactID,
				"erro":       err.Error(),
			})
			utils.ErrorResponse(w, fmt.Errorf("chave estrangeira inválida"), http.StatusBadRequest)
			return
		case errors.Is(err, errMsg.ErrRelationExists):
			h.logger.Info(ctx, ref+logger.LogAlreadyExists, map[string]any{
				"user_id":    modelRelation.UserID,
				"contact_id": modelRelation.ContactID,
			})
			utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
				Data:    dto.ToContactRelationDTO(created),
				Message: "Relação já existente",
				Status:  http.StatusOK,
			})
			return
		default:
			h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
				"user_id":    modelRelation.UserID,
				"contact_id": modelRelation.ContactID,
			})
			utils.ErrorResponse(w, fmt.Errorf("erro ao criar relação: %v", err), http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"user_id":    modelRelation.UserID,
		"contact_id": modelRelation.ContactID,
	})

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Data:    dto.ToContactRelationDTO(created),
		Message: "Relação criada com sucesso",
		Status:  http.StatusCreated,
	})
}

// --- DELETE ---
func (h *UserContactRelation) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserContactRelationHandler - Delete] "
	ctx := r.Context()

	userID, err := utils.GetIDParam(r, "user_id")
	if err != nil {
		h.logger.Warn(ctx, ref+"ID inválido", map[string]any{"user_id": userID})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	contactID, err := utils.GetIDParam(r, "contact_id")
	if err != nil {
		h.logger.Warn(ctx, ref+"ID inválido", map[string]any{"contact_id": contactID})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(ctx, userID, contactID); err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao deletar relação", map[string]any{
			"user_id":    userID,
			"contact_id": contactID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+"Relação deletada com sucesso", map[string]any{
		"user_id":    userID,
		"contact_id": contactID,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    nil,
		Message: "Relação deletada com sucesso",
		Status:  http.StatusOK,
	})
}

// --- DELETE ALL ---
func (h *UserContactRelation) DeleteAll(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserContactRelationHandler - DeleteAll] "
	ctx := r.Context()

	userID, err := utils.GetIDParam(r, "user_id")
	if err != nil {
		h.logger.Warn(ctx, ref+"ID inválido", map[string]any{"user_id": userID})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteAll(ctx, userID); err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao deletar relações", map[string]any{
			"user_id": userID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+"Relações deletadas com sucesso", map[string]any{
		"user_id": userID,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    nil,
		Message: "Relações deletadas com sucesso",
		Status:  http.StatusOK,
	})
}
