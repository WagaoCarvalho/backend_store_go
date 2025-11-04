package handler

import (
	"errors"
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *SupplierHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierHandler - GetAll] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{})

	suppliers, err := h.service.GetAll(ctx)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{})
		utils.ErrorResponse(w, fmt.Errorf("erro ao buscar fornecedores: %w", err), http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"quantidade": len(suppliers),
	})

	supplierDTO := dto.ToSupplierDTOs(suppliers)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Fornecedores encontrados",
		Data:    supplierDTO,
	})
}

func (h *SupplierHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierHandler - GetByID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	supplier, err := h.service.GetByID(ctx, id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "fornecedor não encontrado" {
			status = http.StatusNotFound
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"supplier_id": id,
			})
		} else {
			h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
				"supplier_id": id,
				"status":      status,
			})
		}

		utils.ErrorResponse(w, err, status)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"supplier_id": supplier.ID,
		"name":        supplier.Name,
	})

	createdDTO := dto.ToSupplierDTO(supplier)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Fornecedor encontrado",
		Data:    createdDTO,
	})
}

func (h *SupplierHandler) GetVersionByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierHandler - GetVersionByID] "
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

	version, err := h.service.GetVersionByID(ctx, id)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, errMsg.ErrNotFound) {
			status = http.StatusNotFound
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"supplier_id": id,
			})
		} else {
			h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
				"supplier_id": id,
				"status":      status,
			})
		}

		utils.ErrorResponse(w, err, status)
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

func (h *SupplierHandler) GetByName(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierHandler - GetByName] "
	ctx := r.Context()

	name, err := utils.GetStringParam(r, "name")

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{"name": name})

	suppliers, err := h.service.GetByName(ctx, name)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, errMsg.ErrNotFound) {
			status = http.StatusNotFound
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{"name": name})
		} else {
			h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{"name": name, "status": status})
		}
		utils.ErrorResponse(w, err, status)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"count": len(suppliers),
	})

	supplierDTO := dto.ToSupplierDTOs(suppliers)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Fornecedores encontrados com sucesso",
		Data:    supplierDTO,
	})
}
