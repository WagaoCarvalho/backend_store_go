package handler

import (
	"errors"
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/product/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/product/category_relation"
)

type ProductCategoryRelation struct {
	service service.ProductCategoryRelation
	logger  *logger.LogAdapter
}

func NewProductCategoryRelation(service service.ProductCategoryRelation, logger *logger.LogAdapter) *ProductCategoryRelation {
	return &ProductCategoryRelation{
		service: service,
		logger:  logger,
	}
}

func (h *ProductCategoryRelation) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[ProductCategoryRelationHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{})

	var requestData dto.ProductCategoryRelationsDTO
	if err := utils.FromJSON(r.Body, &requestData); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	modelRelation := dto.ToProductCategoryRelationsModel(requestData)

	created, wasCreated, err := h.service.Create(ctx, modelRelation.ProductID, modelRelation.CategoryID)
	if err != nil {
		if errors.Is(err, errMsg.ErrDBInvalidForeignKey) {
			h.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"product_id":  modelRelation.ProductID,
				"category_id": modelRelation.CategoryID,
				"erro":        err.Error(),
			})
			utils.ErrorResponse(w, fmt.Errorf("chave estrangeira inválida: %w", err), http.StatusBadRequest)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"product_id":  modelRelation.ProductID,
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
		"product_id":  modelRelation.ProductID,
		"category_id": modelRelation.CategoryID,
	})

	createdDTO := dto.ToProductCategoryRelationsDTO(created)

	utils.ToJSON(w, status, utils.DefaultResponse{
		Data:    createdDTO,
		Message: message,
		Status:  status,
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

	relations, err := h.service.GetAllRelationsByProductID(ctx, id)
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

	exists, err := h.service.HasProductCategoryRelation(ctx, productID, categoryID)
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

	if err := h.service.Delete(ctx, productID, categoryID); err != nil {
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

	if err := h.service.DeleteAll(ctx, productID); err != nil {
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
