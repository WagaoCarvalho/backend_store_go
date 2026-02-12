package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	dtoFilter "github.com/WagaoCarvalho/backend_store_go/internal/dto/product/filter"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *productFilterHandler) Filter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const ref = "[productFilterHandler - Filter] "

	// Validação de método HTTP
	if r.Method != http.MethodGet {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	// VALIDAÇÃO DE PAGINAÇÃO - PRIMEIRO, antes de qualquer outra coisa
	limit, offset := utils.GetPaginationParams(r)

	// Validar limit e offset ANTES de qualquer processamento
	if limit < 0 || offset < 0 {
		h.logger.Warn(ctx, ref+"paginação inválida", map[string]any{
			"limit":  limit,
			"offset": offset,
		})
		utils.ErrorResponse(w, fmt.Errorf("parâmetros de paginação inválidos"), http.StatusBadRequest)
		return
	}

	var dtoFilter dtoFilter.ProductFilterDTO
	query := r.URL.Query()

	dtoFilter.ProductName = query.Get("product_name")
	dtoFilter.Manufacturer = query.Get("manufacturer")
	dtoFilter.Barcode = query.Get("barcode")
	dtoFilter.Limit = limit
	dtoFilter.Offset = offset

	if v := query.Get("status"); v != "" {
		if parsed, err := strconv.ParseBool(v); err == nil {
			dtoFilter.Status = &parsed
		} else {
			h.logger.Warn(ctx, ref+"status inválido", map[string]any{
				"valor": v,
			})
		}
	}

	if v := query.Get("supplier_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			dtoFilter.SupplierID = &parsed
		} else {
			h.logger.Warn(ctx, ref+"supplier_id inválido", map[string]any{
				"valor": v,
			})
		}
	}

	if v := query.Get("allow_discount"); v != "" {
		if parsed, err := strconv.ParseBool(v); err == nil {
			dtoFilter.AllowDiscount = &parsed
		} else {
			h.logger.Warn(ctx, ref+"allow_discount inválido", map[string]any{
				"valor": v,
			})
		}
	}

	filter, err := dtoFilter.ToModel()
	if err != nil {
		h.logger.Warn(ctx, ref+"filtro inválido", map[string]any{
			"erro":   err.Error(),
			"filtro": dtoFilter,
		})
		utils.ErrorResponse(w, fmt.Errorf("filtro inválido"), http.StatusBadRequest)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"filtro": dtoFilter,
	})

	products, err := h.service.Filter(ctx, filter)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrInvalidFilter):
			h.logger.Warn(ctx, ref+"filtro inválido", map[string]any{
				"erro":   err.Error(),
				"filtro": dtoFilter,
			})
			utils.ErrorResponse(w, fmt.Errorf("filtro inválido"), http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrZeroID):
			h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
			return

		default:
			h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
				"filtro": dtoFilter,
			})
			utils.ErrorResponse(w, fmt.Errorf("erro ao filtrar produtos"), http.StatusInternalServerError)
			return
		}
	}

	productDTOs := dto.ToProductDTOs(products)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"total_encontrados": len(productDTOs),
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Produtos listados com sucesso",
		Data: map[string]any{
			"total": len(productDTOs),
			"items": productDTOs,
		},
	})
}
