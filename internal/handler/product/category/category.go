package handler

import (
	"errors"
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/product/category"
	iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

type ProductCategory struct {
	service iface.ProductCategory
	logger  *logger.LogAdapter
}

func NewProductCategory(service iface.ProductCategory, logger *logger.LogAdapter) *ProductCategory {
	return &ProductCategory{
		service: service,
		logger:  logger,
	}
}

func (h *ProductCategory) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[ProductCategoryHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{})

	var requestDTO dto.ProductCategoryDTO
	if err := utils.FromJSON(r.Body, &requestDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	categoryModel := dto.ToProductCategoryModel(requestDTO)

	createdCategory, err := h.service.Create(ctx, categoryModel)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"category_id": createdCategory.ID,
	})

	createdDTO := dto.ToProductCategoryDTO(createdCategory)

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Data:    createdDTO,
		Message: "Categoria criada com sucesso",
		Status:  http.StatusCreated,
	})
}

func (h *ProductCategory) GetByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[ProductCategoryHandler - GetById] "
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

	productDTO := dto.ToProductCategoryDTO(category)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    productDTO,
		Message: "Categoria recuperada com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *ProductCategory) GetAll(w http.ResponseWriter, r *http.Request) {
	const ref = "[ProductCategoryHandler - GetAll] "
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

	productDTOs := dto.ToProductCategoryDTOs(categories)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    productDTOs,
		Message: "Categorias recuperadas com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *ProductCategory) Update(w http.ResponseWriter, r *http.Request) {
	const ref = "[ProductCategoryHandler - Update] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{})

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	var requestDTO dto.ProductCategoryDTO

	if err := utils.FromJSON(r.Body, &requestDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	modelCategory := dto.ToProductCategoryModel(requestDTO)
	modelCategory.ID = uint(id)

	err = h.service.Update(ctx, modelCategory)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("categoria não encontrada"), http.StatusNotFound)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"id": id,
		})
		utils.ErrorResponse(w, fmt.Errorf("erro ao atualizar categoria: %v", err), http.StatusInternalServerError)
		return
	}

	updatedCategory, err := h.service.GetByID(ctx, int64(id))
	if err != nil {
		h.logger.Error(ctx, err, ref+"erro ao buscar categoria atualizada", map[string]any{
			"id": id,
		})
		utils.ErrorResponse(w, fmt.Errorf("erro ao buscar categoria atualizada"), http.StatusInternalServerError)
		return
	}

	updatedDTO := dto.ToProductCategoryDTO(updatedCategory)

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"id": id,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    updatedDTO,
		Message: "Categoria atualizada com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *ProductCategory) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[ProductCategoryHandler - Delete] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{})

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		h.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}
