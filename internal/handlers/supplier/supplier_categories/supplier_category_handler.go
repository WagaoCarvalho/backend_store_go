package handlers

import (
	"errors"
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_categories"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repositories/suppliers/supplier_categories"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/suppliers/supplier_categories"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
)

type SupplierCategoryHandler struct {
	service services.SupplierCategoryService
	logger  *logger.LoggerAdapter
}

func NewSupplierCategoryHandler(service services.SupplierCategoryService, logger *logger.LoggerAdapter) *SupplierCategoryHandler {
	return &SupplierCategoryHandler{
		service: service,
		logger:  logger,
	}
}

func (h *SupplierCategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	ref := "[SupplierCategoryHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, nil)

	var category *models.SupplierCategory
	if err := utils.FromJson(r.Body, &category); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJsonError, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	createdCategory, err := h.service.Create(ctx, category)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogCreateError, nil)
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{"category_id": createdCategory.ID})
	utils.ToJson(w, http.StatusCreated, utils.DefaultResponse{
		Data:    createdCategory,
		Message: "Categoria de fornecedor criada com sucesso",
		Status:  http.StatusCreated,
	})
}

func (h *SupplierCategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ref := "[SupplierCategoryHandler - GetByID] "
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
		if errors.Is(err, repo.ErrSupplierCategoryNotFound) {
			statusCode = http.StatusNotFound
		}

		utils.ErrorResponse(w, err, statusCode)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{"category_id": category.ID})
	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Data:    category,
		Message: "Categoria encontrada com sucesso",
		Status:  http.StatusOK,
	})
}
func (h *SupplierCategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ref := "[SupplierCategoryHandler - GetAll] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, nil)

	categories, err := h.service.GetAll(ctx)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, nil)
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{"total": len(categories)})
	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Data:    categories,
		Message: "Categorias encontradas com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *SupplierCategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	ref := "[SupplierCategoryHandler - Update] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogUpdateInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+"ID inválido no path", map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	var category *models.SupplierCategory
	if err := utils.FromJson(r.Body, &category); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJsonError, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	category.ID = id

	err = h.service.Update(ctx, category)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{"category_id": category.ID})

		statusCode := http.StatusInternalServerError
		if errors.Is(err, repo.ErrSupplierCategoryNotFound) {
			statusCode = http.StatusNotFound
		}

		utils.ErrorResponse(w, err, statusCode)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{"category_id": category.ID})
	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Message: "Categoria atualizada com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *SupplierCategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierCategoryHandler - Delete] "
	ctx := r.Context()

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{
		"category_id": id,
		"path":        r.URL.Path,
	})

	if err := h.service.Delete(ctx, id); err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, repo.ErrSupplierCategoryNotFound) {
			statusCode = http.StatusNotFound
		}

		h.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"category_id": id,
		})
		utils.ErrorResponse(w, err, statusCode)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"category_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}
