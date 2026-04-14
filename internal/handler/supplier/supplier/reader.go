package handler

import (
	"errors"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *supplierHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[supplierHandler - GetByID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	supplierModel, err := h.service.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			h.logger.Warn(ctx, ref+logger.LogGetError, map[string]any{
				"supplier_id": id,
			})
			utils.ErrorResponse(w, errMsg.ErrNotFound, http.StatusNotFound)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"supplier_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	supplierDTO := dto.ToSupplierDTO(supplierModel)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"supplier_id": supplierDTO.ID,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Fornecedor encontrado",
		Data:    supplierDTO,
	})
}

func (h *supplierHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	const ref = "[supplierHandler - GetAll] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, nil)

	suppliersModel, err := h.service.GetAll(ctx)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, nil)
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	suppliersDTO := dto.ToSupplierDTOs(suppliersModel)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"quantidade": len(suppliersDTO),
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Fornecedores encontrados",
		Data:    suppliersDTO,
	})
}

func (h *supplierHandler) GetByName(w http.ResponseWriter, r *http.Request) {
	const ref = "[supplierHandler - GetByName] "
	ctx := r.Context()

	name, err := utils.GetStringParam(r, "name")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidParam, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"name": name,
	})

	suppliersModel, err := h.service.GetByName(ctx, name)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			h.logger.Warn(ctx, ref+logger.LogGetError, map[string]any{
				"name": name,
			})
			utils.ErrorResponse(w, errMsg.ErrNotFound, http.StatusNotFound)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"name": name,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	suppliersDTO := dto.ToSupplierDTOs(suppliersModel)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"quantidade": len(suppliersDTO),
		"name":       name,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Fornecedores encontrados",
		Data:    suppliersDTO,
	})
}

func (h *supplierHandler) GetVersionByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[supplierHandler - GetVersionByID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	version, err := h.service.GetVersionByID(ctx, id)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			h.logger.Warn(ctx, ref+logger.LogGetError, map[string]any{
				"supplier_id": id,
			})
			utils.ErrorResponse(w, errMsg.ErrNotFound, http.StatusNotFound)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"supplier_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"supplier_id": id,
		"version":     version,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Versão do fornecedor obtida com sucesso",
		Data: map[string]int64{
			"version": version,
		},
	})
}
