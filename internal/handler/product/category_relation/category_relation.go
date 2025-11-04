package handler

import (
	"errors"
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/product/category_relation"
	iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

type ProductCategoryRelation struct {
	productCategoryRelation iface.ProductCategoryRelation
	logger                  *logger.LogAdapter
}

func NewProductCategoryRelation(productCategoryRelation iface.ProductCategoryRelation, logger *logger.LogAdapter) *ProductCategoryRelation {
	return &ProductCategoryRelation{
		productCategoryRelation: productCategoryRelation,
		logger:                  logger,
	}
}

func (h *ProductCategoryRelation) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[ProductCategoryRelationHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{})

	var requestData dto.ProductCategoryRelationsDTO
	if err := utils.FromJSON(r.Body, &requestData); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, fmt.Errorf("erro ao decodificar JSON"), http.StatusBadRequest)
		return
	}

	modelRelation := dto.ToProductCategoryRelationsModel(requestData)

	// Validação simples antes de chamar o service
	if modelRelation == nil || modelRelation.ProductID <= 0 || modelRelation.CategoryID <= 0 {
		h.logger.Warn(ctx, ref+"modelo nulo ou ID inválido", map[string]any{})
		utils.ErrorResponse(w, fmt.Errorf("modelo nulo ou ID inválido"), http.StatusBadRequest)
		return
	}

	created, err := h.productCategoryRelation.Create(ctx, modelRelation)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrDBInvalidForeignKey):
			h.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"product_id":  modelRelation.ProductID,
				"category_id": modelRelation.CategoryID,
				"erro":        err.Error(),
			})
			utils.ErrorResponse(w, fmt.Errorf("chave estrangeira inválida"), http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrRelationExists):
			h.logger.Info(ctx, ref+logger.LogAlreadyExists, map[string]any{
				"product_id":  modelRelation.ProductID,
				"category_id": modelRelation.CategoryID,
			})
			utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
				Data:    dto.ToProductCategoryRelationsDTO(created),
				Message: "Relação já existente",
				Status:  http.StatusOK,
			})
			return

		default:
			h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
				"product_id":  modelRelation.ProductID,
				"category_id": modelRelation.CategoryID,
			})
			utils.ErrorResponse(w, fmt.Errorf("erro ao criar relação: %v", err), http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"product_id":  modelRelation.ProductID,
		"category_id": modelRelation.CategoryID,
	})

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Data:    dto.ToProductCategoryRelationsDTO(created),
		Message: "Relação criada com sucesso",
		Status:  http.StatusCreated,
	})
}

func (h *ProductCategoryRelation) GetAllRelationsByProductID(w http.ResponseWriter, r *http.Request) {
	const ref = "[ProductCategoryRelationHandler - GetAllRelationsByProductID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{})

	id, err := utils.GetIDParam(r, "product_id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID de usuário inválido"), http.StatusBadRequest)
		return
	}

	relations, err := h.productCategoryRelation.GetAllRelationsByProductID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"product_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// 🔹 Conversão para DTO
	var relationsDTO []dto.ProductCategoryRelationsDTO
	for _, rel := range relations {
		relationsDTO = append(relationsDTO, dto.ToProductCategoryRelationsDTO(rel))
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"product_id": id,
		"total":      len(relationsDTO),
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    relationsDTO,
		Message: "Relações recuperadas com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *ProductCategoryRelation) HasProductCategoryRelation(w http.ResponseWriter, r *http.Request) {
	const ref = "[ProductCategoryRelationHandler - HasProductCategoryRelation] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogVerificationInit, map[string]any{})

	productID, err := utils.GetIDParam(r, "product_id")

	if err != nil || productID <= 0 {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"campo": "product_id",
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

	exists, err := h.productCategoryRelation.HasProductCategoryRelation(ctx, productID, categoryID)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogVerificationError, map[string]any{
			"product_id":  productID,
			"category_id": categoryID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogVerificationSuccess, map[string]any{
		"product_id":  productID,
		"category_id": categoryID,
		"exists":      exists,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    map[string]bool{"exists": exists},
		Message: "Verificação concluída com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *ProductCategoryRelation) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[ProductCategoryRelationHandler - Delete] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{})

	productID, errProductID := utils.GetIDParam(r, "product_id")
	categoryID, errCategoryID := utils.GetIDParam(r, "category_id")

	if errProductID != nil || errCategoryID != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro_product_id":  errProductID,
			"erro_category_id": errCategoryID,
		})
		utils.ErrorResponse(w, fmt.Errorf("IDs inválidos"), http.StatusBadRequest)
		return
	}

	if err := h.productCategoryRelation.Delete(ctx, productID, categoryID); err != nil {
		h.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"product_id":  productID,
			"category_id": categoryID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"product_id":  productID,
		"category_id": categoryID,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductCategoryRelation) DeleteAll(w http.ResponseWriter, r *http.Request) {
	const ref = "[ProductCategoryRelationHandler - DeleteAll] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{})

	productID, err := utils.GetIDParam(r, "product_id")

	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID de usuário inválido"), http.StatusBadRequest)
		return
	}

	if err := h.productCategoryRelation.DeleteAll(ctx, productID); err != nil {
		h.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"product_id": productID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"product_id": productID,
	})

	w.WriteHeader(http.StatusNoContent)
}
