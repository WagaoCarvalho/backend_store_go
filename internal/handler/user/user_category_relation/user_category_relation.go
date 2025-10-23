package handler

import (
	"errors"
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/user_category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/user_category_relation"
)

type UserCategoryRelation struct {
	service service.UserCategoryRelation
	logger  *logger.LogAdapter
}

func NewUserCategoryRelation(service service.UserCategoryRelation, logger *logger.LogAdapter) *UserCategoryRelation {
	return &UserCategoryRelation{
		service: service,
		logger:  logger,
	}
}

func (h *UserCategoryRelation) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserCategoryRelationHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{})

	var requestData dto.UserCategoryRelationsDTO
	if err := utils.FromJSON(r.Body, &requestData); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	modelRelation := dto.ToUserCategoryRelationsModel(requestData)

	created, wasCreated, err := h.service.Create(ctx, modelRelation.UserID, modelRelation.CategoryID)
	if err != nil {
		if errors.Is(err, errMsg.ErrDBInvalidForeignKey) {
			h.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"user_id":     modelRelation.UserID,
				"category_id": modelRelation.CategoryID,
				"erro":        err.Error(),
			})
			utils.ErrorResponse(w, fmt.Errorf("chave estrangeira inválida: %w", err), http.StatusBadRequest)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"user_id":     modelRelation.UserID,
			"category_id": modelRelation.CategoryID,
		})
		utils.ErrorResponse(w, fmt.Errorf("erro ao criar relação: %w", err), http.StatusInternalServerError)
		return
	}

	status := http.StatusOK
	message := "Relação já existente"
	logMsg := logger.LogAlreadyExists
	if wasCreated {
		status = http.StatusCreated
		message = "Relação criada com sucesso"
		logMsg = logger.LogCreateSuccess
	}

	h.logger.Info(ctx, ref+logMsg, map[string]any{
		"user_id":     modelRelation.UserID,
		"category_id": modelRelation.CategoryID,
	})

	createdDTO := dto.ToUserCategoryRelationsDTO(created)

	utils.ToJSON(w, status, utils.DefaultResponse{
		Data:    createdDTO,
		Message: message,
		Status:  status,
	})
}

func (h *UserCategoryRelation) GetAllRelationsByUserID(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserCategoryRelationHandler - GetAllRelationsByUserID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{})

	id, err := utils.GetIDParam(r, "user_id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID de usuário inválido"), http.StatusBadRequest)
		return
	}

	relations, err := h.service.GetAllRelationsByUserID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"user_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// 🔹 Conversão para DTO
	var relationsDTO []dto.UserCategoryRelationsDTO
	for _, rel := range relations {
		relationsDTO = append(relationsDTO, dto.ToUserCategoryRelationsDTO(rel))
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"user_id": id,
		"total":   len(relationsDTO),
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    relationsDTO,
		Message: "Relações recuperadas com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *UserCategoryRelation) HasUserCategoryRelation(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserCategoryRelationHandler - HasUserCategoryRelation] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogVerificationInit, map[string]any{})

	userID, err := utils.GetIDParam(r, "user_id")

	if err != nil || userID <= 0 {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"campo": "user_id",
			"erro":  err,
		})
		utils.ErrorResponse(w, fmt.Errorf("ID de usuário inválido"), http.StatusBadRequest)
		return
	}

	categoryID, err := utils.GetIDParam(r, "category_id")

	if err != nil || categoryID <= 0 {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"campo": "category_id",
			"erro":  err,
		})
		utils.ErrorResponse(w, fmt.Errorf("ID de categoria inválido"), http.StatusBadRequest)
		return
	}

	exists, err := h.service.HasUserCategoryRelation(ctx, userID, categoryID)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogVerificationError, map[string]any{
			"user_id":     userID,
			"category_id": categoryID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogVerificationSuccess, map[string]any{
		"user_id":     userID,
		"category_id": categoryID,
		"exists":      exists,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    map[string]bool{"exists": exists},
		Message: "Verificação concluída com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *UserCategoryRelation) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *UserCategoryRelation) DeleteAll(w http.ResponseWriter, r *http.Request) {
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
