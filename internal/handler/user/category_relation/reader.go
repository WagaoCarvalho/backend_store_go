package handler

import (
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/category_relation"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

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
