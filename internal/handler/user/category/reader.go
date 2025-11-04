package handler

import (
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/category"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *UserCategory) GetAll(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserCategoryHandler - GetAll] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{})

	categories, err := h.service.GetAll(ctx)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"total": len(categories),
	})

	userDTOs := dto.ToUserCategoryDTOs(categories)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    userDTOs,
		Message: "Categorias recuperadas com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *UserCategory) GetByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserCategoryHandler - GetById] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{})

	id, err := utils.GetIDParam(r, "id")

	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	category, err := h.service.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "categoria não encontrada" {
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"id": id,
			})
			utils.ErrorResponse(w, err, http.StatusNotFound)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"id": id,
	})

	userDTO := dto.ToUserCategoryDTO(category)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    userDTO,
		Message: "Categoria recuperada com sucesso",
		Status:  http.StatusOK,
	})
}
