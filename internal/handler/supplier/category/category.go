package handler

import (
	"errors"
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/category"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/category"
)

type SupplierCategory struct {
	service service.SupplierCategory
	logger  *logger.LogAdapter
}

func NewSupplierCategory(service service.SupplierCategory, logger *logger.LogAdapter) *SupplierCategory {
	return &SupplierCategory{
		service: service,
		logger:  logger,
	}
}

func (h *SupplierCategory) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierCategoryHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, nil)

	var requestDTO dto.SupplierCategoryDTO

	if err := utils.FromJSON(r.Body, &requestDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	modelCategory := dto.ToSupplierCategoryModel(requestDTO)

	createdCategory, err := h.service.Create(ctx, modelCategory)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"name": modelCategory.Name,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{"category_id": createdCategory.ID})

	createdDTO := dto.ToSupplierCategoryDTO(createdCategory)

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Categoria de fornecedor criada com sucesso",
		Data:    createdDTO,
	})
}

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

func (h *SupplierCategory) Update(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierCategoryHandler - Update] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogUpdateInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	var requestDTO dto.SupplierCategoryDTO

	if err := utils.FromJSON(r.Body, &requestDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	modelCategory := dto.ToSupplierCategoryModel(requestDTO)
	modelCategory.ID = id

	if err := h.service.Update(ctx, modelCategory); err != nil {
		h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{"category_id": modelCategory.ID})

		statusCode := http.StatusInternalServerError
		if errors.Is(err, errMsg.ErrNotFound) {
			statusCode = http.StatusNotFound
		}

		utils.ErrorResponse(w, err, statusCode)
		return
	}

	createdDTO := dto.ToSupplierCategoryDTO(modelCategory)

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{"category_id": modelCategory.ID})
	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    createdDTO,
		Message: "Categoria atualizada com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *SupplierCategory) Delete(w http.ResponseWriter, r *http.Request) {
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
		if errors.Is(err, errMsg.ErrNotFound) {
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
