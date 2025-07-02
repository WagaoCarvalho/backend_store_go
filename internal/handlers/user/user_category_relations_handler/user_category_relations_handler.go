package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_category_relations"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/user/user_category_relations"
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
	ctx := r.Context()
	h.logger.Info(ctx, "[UserCategoryRelationHandler] - Iniciando criação de relação usuário-categoria", nil)

	var relation *models.UserCategoryRelations
	if err := utils.FromJson(r.Body, &relation); err != nil {
		h.logger.Warn(ctx, "[UserCategoryRelationHandler] - JSON inválido ao criar relação", map[string]interface{}{
			"error": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	created, wasCreated, err := h.service.Create(ctx, relation.UserID, relation.CategoryID)
	if err != nil {
		if errors.Is(err, repositories.ErrInvalidForeignKey) {
			h.logger.Warn(ctx, "[UserCategoryRelationHandler] - Chave estrangeira inválida", map[string]interface{}{
				"user_id":     relation.UserID,
				"category_id": relation.CategoryID,
				"error":       err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		h.logger.Error(ctx, err, "[UserCategoryRelationHandler] - Erro ao criar relação", map[string]interface{}{
			"user_id":     relation.UserID,
			"category_id": relation.CategoryID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	var status int
	var message string

	if wasCreated {
		status = http.StatusCreated
		message = "Relação criada com sucesso"
		h.logger.Info(ctx, "[UserCategoryRelationHandler] - Relação criada com sucesso", map[string]interface{}{
			"user_id":     relation.UserID,
			"category_id": relation.CategoryID,
		})
	} else {
		status = http.StatusOK
		message = "Relação já existente"
		h.logger.Info(ctx, "[UserCategoryRelationHandler] - Relação já existente", map[string]interface{}{
			"user_id":     relation.UserID,
			"category_id": relation.CategoryID,
		})
	}

	utils.ToJson(w, status, utils.DefaultResponse{
		Data:    created,
		Message: message,
		Status:  status,
	})
}

func (h *UserCategoryRelationHandler) GetAllRelationsByUserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	h.logger.Info(ctx, "[UserCategoryRelationHandler] - Iniciando busca de relações por ID de usuário", nil)

	id, err := strconv.ParseInt(mux.Vars(r)["user_id"], 10, 64)
	if err != nil {
		h.logger.Warn(ctx, "[UserCategoryRelationHandler] - ID de usuário inválido recebido para busca", map[string]interface{}{
			"error": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID de usuário inválido"), http.StatusBadRequest)
		return
	}

	relations, err := h.service.GetAllRelationsByUserID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, err, "[UserCategoryRelationHandler] - Erro ao recuperar relações", map[string]interface{}{
			"user_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, "[UserCategoryRelationHandler] - Relações recuperadas com sucesso", map[string]interface{}{
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
	ctx := r.Context()

	h.logger.Info(ctx, "[UserCategoryRelationHandler] - Iniciando verificação de relação usuário-categoria", nil)

	userID, err := strconv.ParseInt(mux.Vars(r)["user_id"], 10, 64)
	if err != nil || userID <= 0 {
		h.logger.Warn(ctx, "[UserCategoryRelationHandler] - ID de usuário inválido", map[string]interface{}{
			"error": err,
		})
		utils.ErrorResponse(w, fmt.Errorf("ID de usuário inválido"), http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.ParseInt(mux.Vars(r)["category_id"], 10, 64)
	if err != nil || categoryID <= 0 {
		h.logger.Warn(ctx, "[UserCategoryRelationHandler] - ID de categoria inválido", map[string]interface{}{
			"error": err,
		})
		utils.ErrorResponse(w, fmt.Errorf("ID de categoria inválido"), http.StatusBadRequest)
		return
	}

	exists, err := h.service.HasUserCategoryRelation(ctx, userID, categoryID)
	if err != nil {
		h.logger.Error(ctx, err, "[UserCategoryRelationHandler] - Erro ao verificar relação", map[string]interface{}{
			"user_id":     userID,
			"category_id": categoryID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, "[UserCategoryRelationHandler] - Verificação concluída", map[string]interface{}{
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
	ctx := r.Context()

	h.logger.Info(ctx, "[UserCategoryRelationHandler] - Iniciando exclusão de relação usuário-categoria", nil)

	userID, err1 := strconv.ParseInt(mux.Vars(r)["user_id"], 10, 64)
	categoryID, err2 := strconv.ParseInt(mux.Vars(r)["category_id"], 10, 64)

	if err1 != nil || err2 != nil {
		h.logger.Warn(ctx, "[UserCategoryRelationHandler] - IDs inválidos para exclusão de relação", map[string]interface{}{
			"error_user_id":     err1,
			"error_category_id": err2,
		})
		utils.ErrorResponse(w, fmt.Errorf("IDs inválidos"), http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(ctx, userID, categoryID); err != nil {
		h.logger.Error(ctx, err, "[UserCategoryRelationHandler] - Erro ao excluir relação", map[string]interface{}{
			"user_id":     userID,
			"category_id": categoryID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, "[UserCategoryRelationHandler] - Relação excluída com sucesso", map[string]interface{}{
		"user_id":     userID,
		"category_id": categoryID,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserCategoryRelationHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	h.logger.Info(ctx, "[UserCategoryRelationHandler] - Iniciando exclusão de todas as relações do usuário", nil)

	userID, err := strconv.ParseInt(mux.Vars(r)["user_id"], 10, 64)
	if err != nil {
		h.logger.Warn(ctx, "[UserCategoryRelationHandler] - ID de usuário inválido para exclusão em massa", map[string]interface{}{
			"error": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID de usuário inválido"), http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteAll(ctx, userID); err != nil {
		h.logger.Error(ctx, err, "[UserCategoryRelationHandler] - Erro ao excluir todas as relações", map[string]interface{}{
			"user_id": userID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, "[UserCategoryRelationHandler] - Relações do usuário excluídas com sucesso", map[string]interface{}{
		"user_id": userID,
	})

	w.WriteHeader(http.StatusNoContent)
}
