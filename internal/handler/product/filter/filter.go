package handler

import (
	"errors"
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
	const ref = "[productHandler - Filter] "

	var dtoFilter dtoFilter.ProductFilterDTO
	query := r.URL.Query()

	dtoFilter.ProductName = query.Get("product_name")
	dtoFilter.Manufacturer = query.Get("manufacturer")
	dtoFilter.Barcode = query.Get("barcode")

	if v := query.Get("status"); v != "" {
		if parsed, err := strconv.ParseBool(v); err == nil {
			dtoFilter.Status = &parsed
		}
	}

	if v := query.Get("supplier_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			dtoFilter.SupplierID = &parsed
		}
	}

	if v := query.Get("allow_discount"); v != "" {
		if parsed, err := strconv.ParseBool(v); err == nil {
			dtoFilter.AllowDiscount = &parsed
		}
	}

	dtoFilter.Limit, dtoFilter.Offset = utils.GetPaginationParams(r)

	filter, _ := dtoFilter.ToModel()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"filtro": dtoFilter,
	})

	products, err := h.service.Filter(ctx, filter)
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
