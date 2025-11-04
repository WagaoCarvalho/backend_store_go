package handler

import (
	"errors"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/category"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *SupplierCategory) GetByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierCategoryHandler - GetByID] "
	ctx := r.Context()

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+"ID inválido no path", map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{"category_id": id})

	category, err := h.service.GetByID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{"category_id": id})

		statusCode := http.StatusInternalServerError
		if errors.Is(err, errMsg.ErrNotFound) {
			statusCode = http.StatusNotFound
		}

		utils.ErrorResponse(w, err, statusCode)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{"category_id": category.ID})

	createdDTO := dto.ToSupplierCategoryDTO(category)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    createdDTO,
		Message: "Categoria encontrada com sucesso",
		Status:  http.StatusOK,
	})
}
func (h *SupplierCategory) GetAll(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierCategoryHandler - GetAll] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, nil)

	categories, err := h.service.GetAll(ctx)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, nil)
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{"total": len(categories)})

	categoryDTO := dto.ToSupplierCategoryDTOs(categories)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    categoryDTO,
		Message: "Categorias encontradas com sucesso",
		Status:  http.StatusOK,
	})
}
