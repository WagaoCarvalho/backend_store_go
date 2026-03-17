package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	dtoFilter "github.com/WagaoCarvalho/backend_store_go/internal/dto/sale/filter"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

// Lista de parâmetros válidos para validação
var validSaleFilterParams = map[string]bool{
	"client_id":      true,
	"user_id":        true,
	"status":         true,
	"payment_type":   true,
	"sale_date_from": true,
	"sale_date_to":   true,
	"limit":          true,
	"offset":         true,
}

func (h *saleFilterHandler) Filter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const ref = "[saleHandler - Filter] "

	query := r.URL.Query()

	// VALIDAÇÃO 1: Verificar parâmetros desconhecidos na query
	for param := range query {
		if !validSaleFilterParams[param] {
			h.logger.Warn(ctx, ref+"parâmetro desconhecido", map[string]any{
				"parametro": param,
				"valor":     query.Get(param),
			})
			utils.ErrorResponse(w, fmt.Errorf("parâmetro de consulta inválido: %s", param), http.StatusBadRequest)
			return
		}
	}

	var dtoFilter dtoFilter.SaleFilterDTO

	// VALIDAÇÃO 2: client_id com valor inválido deve retornar erro
	if v := query.Get("client_id"); v != "" {
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			h.logger.Warn(ctx, ref+"client_id inválido", map[string]any{
				"valor": v,
			})
			utils.ErrorResponse(w, fmt.Errorf("client_id deve ser um número inteiro"), http.StatusBadRequest)
			return
		}
		dtoFilter.ClientID = &parsed
	}

	// VALIDAÇÃO 3: user_id com valor inválido deve retornar erro
	if v := query.Get("user_id"); v != "" {
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			h.logger.Warn(ctx, ref+"user_id inválido", map[string]any{
				"valor": v,
			})
			utils.ErrorResponse(w, fmt.Errorf("user_id deve ser um número inteiro"), http.StatusBadRequest)
			return
		}
		dtoFilter.UserID = &parsed
	}

	// Campos de string (não precisam de validação de formato)
	dtoFilter.Status = query.Get("status")
	dtoFilter.PaymentType = query.Get("payment_type")

	// Datas via utilitário (assumindo que ParseTimeRange já valida formato)
	utils.ParseTimeRange(
		query,
		"sale_date_from",
		"sale_date_to",
		&dtoFilter.SaleDateFrom,
		&dtoFilter.SaleDateTo,
	)

	// VALIDAÇÃO 4: Paginação (limit/offset não negativos)
	limit, offset := utils.GetPaginationParams(r)
	if limit < 0 || offset < 0 {
		h.logger.Warn(ctx, ref+"paginação inválida", map[string]any{
			"limit":  limit,
			"offset": offset,
		})
		utils.ErrorResponse(w, fmt.Errorf("parâmetros de paginação inválidos"), http.StatusBadRequest)
		return
	}
	dtoFilter.Limit = limit
	dtoFilter.Offset = offset

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
