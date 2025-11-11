package handler

import (
	"errors"
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *userCategoryRelationHandler) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserCategoryRelationHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{})

	var requestData dto.UserCategoryRelationsDTO
	if err := utils.FromJSON(r.Body, &requestData); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, fmt.Errorf("erro ao decodificar JSON"), http.StatusBadRequest)
		return
	}

	modelRelation := dto.ToUserCategoryRelationsModel(requestData)

	// Validação simples de IDs antes de chamar o service
	if modelRelation == nil || modelRelation.UserID <= 0 || modelRelation.CategoryID <= 0 {
		h.logger.Warn(ctx, ref+"modelo nulo ou ID inválido", map[string]any{})
		utils.ErrorResponse(w, fmt.Errorf("modelo nulo ou ID inválido"), http.StatusBadRequest)
		return
	}

	created, err := h.service.Create(ctx, modelRelation)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrDBInvalidForeignKey):
			h.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"user_id":     modelRelation.UserID,
				"category_id": modelRelation.CategoryID,
				"erro":        err.Error(),
			})
			utils.ErrorResponse(w, fmt.Errorf("chave estrangeira inválida"), http.StatusBadRequest)
			return
		case errors.Is(err, errMsg.ErrRelationExists):
			h.logger.Info(ctx, ref+logger.LogAlreadyExists, map[string]any{
				"user_id":     modelRelation.UserID,
				"category_id": modelRelation.CategoryID,
			})
			utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
				Data:    dto.ToUserCategoryRelationsDTO(created),
				Message: "Relação já existente",
				Status:  http.StatusOK,
			})
			return
		default:
			h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
				"user_id":     modelRelation.UserID,
				"category_id": modelRelation.CategoryID,
			})
			utils.ErrorResponse(w, fmt.Errorf("erro ao criar relação: %v", err), http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"user_id":     modelRelation.UserID,
		"category_id": modelRelation.CategoryID,
	})

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Data:    dto.ToUserCategoryRelationsDTO(created),
		Message: "Relação criada com sucesso",
		Status:  http.StatusCreated,
	})
}

func (h *userCategoryRelationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserCategoryRelationHandler - Delete] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{})

	userID, errUserID := utils.GetIDParam(r, "user_id")
	categoryID, errCategoryID := utils.GetIDParam(r, "category_id")

	if errUserID != nil || errCategoryID != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro_user_id":     errUserID,
			"erro_category_id": errCategoryID,
		})
		utils.ErrorResponse(w, fmt.Errorf("IDs inválidos"), http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(ctx, userID, categoryID); err != nil {
		h.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"user_id":     userID,
			"category_id": categoryID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"user_id":     userID,
		"category_id": categoryID,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (h *userCategoryRelationHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserCategoryRelationHandler - DeleteAll] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{})

	userID, err := utils.GetIDParam(r, "user_id")

	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID de usuário inválido"), http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteAll(ctx, userID); err != nil {
		h.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"user_id": userID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"user_id": userID,
	})

	w.WriteHeader(http.StatusNoContent)
}
