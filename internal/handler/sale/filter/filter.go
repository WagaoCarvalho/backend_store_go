package handler

import (
	"errors"
	"net/http"
	"strconv"

	dtoFilter "github.com/WagaoCarvalho/backend_store_go/internal/dto/sale/filter"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *saleFilterHandler) Filter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const ref = "[saleHandler - Filter] "

	var dtoFilter dtoFilter.SaleFilterDTO
	query := r.URL.Query()

	if v := query.Get("client_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			dtoFilter.ClientID = &parsed
		}
	}

	if v := query.Get("user_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			dtoFilter.UserID = &parsed
		}
	}

	dtoFilter.Status = query.Get("status")
	dtoFilter.PaymentType = query.Get("payment_type")

	// ✅ CORRETO: datas via utilitário
	utils.ParseTimeRange(
		query,
		"sale_date_from",
		"sale_date_to",
		&dtoFilter.SaleDateFrom,
		&dtoFilter.SaleDateTo,
	)

	dtoFilter.Limit, dtoFilter.Offset = utils.GetPaginationParams(r)

	filter, err := dtoFilter.ToModel()
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{"filtro": dtoFilter})

	sales, err := h.service.Filter(ctx, filter)
	if err != nil {
		if errors.Is(err, errMsg.ErrInvalidFilter) {
			h.logger.Warn(ctx, ref+"filtro inválido", map[string]any{
				"erro":   err.Error(),
				"filtro": dtoFilter,
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{"filtro": dtoFilter})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	saleDTOs := dto.ToSaleDTOs(sales)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"total_encontrados": len(saleDTOs),
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Vendas listadas com sucesso",
		Data: map[string]any{
			"total": len(saleDTOs),
			"items": saleDTOs,
		},
	})
}
