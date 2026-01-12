package handler

import (
	"errors"
	"net/http"
	"strconv"

	dtoFilter "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/filter"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *supplierFilterHandler) Filter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const ref = "[supplierHandler - Filter] "

	var dtoFilter dtoFilter.SupplierFilterDTO
	query := r.URL.Query()

	dtoFilter.Name = query.Get("name")
	dtoFilter.CPF = query.Get("cpf")
	dtoFilter.CNPJ = query.Get("cnpj")

	if v := query.Get("status"); v != "" {
		if parsed, err := strconv.ParseBool(v); err == nil {
			dtoFilter.Status = &parsed
		}
	}

	// Datas via utilitário
	utils.ParseTimeRange(
		query,
		"created_from",
		"created_to",
		&dtoFilter.CreatedFrom,
		&dtoFilter.CreatedTo,
	)

	utils.ParseTimeRange(
		query,
		"updated_from",
		"updated_to",
		&dtoFilter.UpdatedFrom,
		&dtoFilter.UpdatedTo,
	)

	dtoFilter.Limit, dtoFilter.Offset = utils.GetPaginationParams(r)

	filter, err := dtoFilter.ToModel()
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"filtro": dtoFilter,
	})

	suppliers, err := h.service.Filter(ctx, filter)
	if err != nil {
		if errors.Is(err, errMsg.ErrInvalidFilter) {
			h.logger.Warn(ctx, ref+"filtro inválido", map[string]any{
				"erro":   err.Error(),
				"filtro": dtoFilter,
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"filtro": dtoFilter,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	supplierDTOs := dto.ToSupplierDTOs(suppliers)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"total_encontrados": len(supplierDTOs),
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Fornecedores listados com sucesso",
		Data: map[string]any{
			"total": len(supplierDTOs),
			"items": supplierDTOs,
		},
	})
}
