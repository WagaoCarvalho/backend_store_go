package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_category_relations"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/users/user_category_relations"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	"github.com/gorilla/mux"
)

type UserCategoryRelationHandler struct {
	service services.UserCategoryRelationServices
	logger  *logger.LoggerAdapter
}

func NewUserCategoryRelationHandler(service services.UserCategoryRelationServices, logger *logger.LoggerAdapter) *UserCategoryRelationHandler {
	return &UserCategoryRelationHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UserCategoryRelationHandler) Create(w http.ResponseWriter, r *http.Request) {
	ref := "[UserCategoryRelationHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, map[string]interface{}{})

	var relation *models.UserCategoryRelations
	if err := utils.FromJson(r.Body, &relation); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJsonError, map[string]interface{}{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	created, wasCreated, err := h.service.Create(ctx, relation.UserID, relation.CategoryID)
	if err != nil {
		if errors.Is(err, repositories.ErrInvalidForeignKey) {
			h.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]interface{}{
				"user_id":     relation.UserID,
				"category_id": relation.CategoryID,
				"erro":        err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]interface{}{
			"user_id":     relation.UserID,
			"category_id": relation.CategoryID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	var status int
	var message string
	var logMsg string

	if wasCreated {
		status = http.StatusCreated
		message = "Relação criada com sucesso"
		logMsg = logger.LogCreateSuccess
	} else {
		status = http.StatusOK
		message = "Relação já existente"
		logMsg = logger.LogAlreadyExists
	}

	h.logger.Info(ctx, ref+logMsg, map[string]interface{}{
		"user_id":     relation.UserID,
		"category_id": relation.CategoryID,
	})

	utils.ToJson(w, status, utils.DefaultResponse{
		Data:    created,
		Message: message,
		Status:  status,
	})
}

func (h *UserCategoryRelationHandler) GetAllRelationsByUserID(w http.ResponseWriter, r *http.Request) {
	ref := "[UserCategoryRelationHandler - GetAllRelationsByUserID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]interface{}{})

	id, err := strconv.ParseInt(mux.Vars(r)["user_id"], 10, 64)
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]interface{}{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID de usuário inválido"), http.StatusBadRequest)
		return
	}

	relations, err := h.service.GetAllRelationsByUserID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]interface{}{
			"user_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]interface{}{
		"user_id": id,
		"total":   len(relations),
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Data:    relations,
		Message: "Relações recuperadas com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *UserCategoryRelationHandler) HasUserCategoryRelation(w http.ResponseWriter, r *http.Request) {
	ref := "[UserCategoryRelationHandler - HasUserCategoryRelation] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogVerificationInit, map[string]interface{}{})

	userID, err := strconv.ParseInt(mux.Vars(r)["user_id"], 10, 64)
	if err != nil || userID <= 0 {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]interface{}{
			"campo": "user_id",
			"erro":  err,
		})
		utils.ErrorResponse(w, fmt.Errorf("ID de usuário inválido"), http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.ParseInt(mux.Vars(r)["category_id"], 10, 64)
	if err != nil || categoryID <= 0 {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]interface{}{
			"campo": "category_id",
			"erro":  err,
		})
		utils.ErrorResponse(w, fmt.Errorf("ID de categoria inválido"), http.StatusBadRequest)
		return
	}

	exists, err := h.service.HasUserCategoryRelation(ctx, userID, categoryID)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogVerificationError, map[string]interface{}{
			"user_id":     userID,
			"category_id": categoryID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogVerificationSuccess, map[string]interface{}{
		"user_id":     userID,
		"category_id": categoryID,
		"exists":      exists,
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Data:    map[string]bool{"exists": exists},
		Message: "Verificação concluída com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *UserCategoryRelationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ref := "[UserCategoryRelationHandler - Delete] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]interface{}{})

	userID, errUserID := strconv.ParseInt(mux.Vars(r)["user_id"], 10, 64)
	categoryID, errCategoryID := strconv.ParseInt(mux.Vars(r)["category_id"], 10, 64)

	if errUserID != nil || errCategoryID != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]interface{}{
			"erro_user_id":     errUserID,
			"erro_category_id": errCategoryID,
		})
		utils.ErrorResponse(w, fmt.Errorf("IDs inválidos"), http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(ctx, userID, categoryID); err != nil {
		h.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]interface{}{
			"user_id":     userID,
			"category_id": categoryID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]interface{}{
		"user_id":     userID,
		"category_id": categoryID,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserCategoryRelationHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	ref := "[UserCategoryRelationHandler - DeleteAll] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]interface{}{})

	userID, err := strconv.ParseInt(mux.Vars(r)["user_id"], 10, 64)
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]interface{}{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID de usuário inválido"), http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteAll(ctx, userID); err != nil {
		h.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]interface{}{
			"user_id": userID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]interface{}{
		"user_id": userID,
	})

	w.WriteHeader(http.StatusNoContent)
}
