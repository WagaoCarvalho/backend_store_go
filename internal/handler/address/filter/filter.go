package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/address/address"
	filterDTO "github.com/WagaoCarvalho/backend_store_go/internal/dto/address/filter"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *addressFilterHandler) Filter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const ref = "[addressHandler - Filter] "

	if err := h.validateQueryParams(r); err != nil {
		h.logger.Warn(ctx, ref+"validação de parâmetros falhou", map[string]any{
			"erro":  err.Error(),
			"query": r.URL.Query(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	dtoFilter, err := h.parseFilterDTO(r)
	if err != nil {
		h.logger.Warn(ctx, ref+"erro ao parsear filtros", map[string]any{
			"erro":  err.Error(),
			"query": r.URL.Query(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := dtoFilter.Validate(); err != nil {
		h.logger.Warn(ctx, ref+"validação de DTO falhou", map[string]any{
			"erro": err.Error(),
			"dto":  dtoFilter,
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	filter, err := dtoFilter.ToModel()
	if err != nil {
		h.logger.Warn(ctx, ref+"erro ao converter filtro", map[string]any{
			"erro": err.Error(),
			"dto":  dtoFilter,
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{"filtro": dtoFilter})

	addresses, err := h.service.Filter(ctx, filter)
	if err != nil {
		if errors.Is(err, errMsg.ErrInvalidFilter) {
			h.logger.Warn(ctx, ref+"filtro inválido", map[string]any{"erro": err.Error(), "filtro": dtoFilter})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{"filtro": dtoFilter})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	addressDTOs := dto.ToAddressDTOs(addresses)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"total_encontrados": len(addressDTOs),
		"filtros_aplicados": countFiltersApplied(dtoFilter),
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereços listados com sucesso",
		Data: map[string]any{
			"total":           len(addressDTOs),
			"items":           addressDTOs,
			"filters_applied": countFiltersApplied(dtoFilter),
			"has_more":        len(addressDTOs) == dtoFilter.Limit,
		},
	})
}

func (h *addressFilterHandler) validateQueryParams(r *http.Request) error {
	allowedParams := map[string]bool{
		"user_id":       true,
		"client_cpf_id": true,
		"supplier_id":   true,
		"street":        true,
		"street_number": true,
		"complement":    true,
		"city":          true,
		"state":         true,
		"country":       true,
		"postal_code":   true,
		"is_active":     true,
		"created_from":  true,
		"created_to":    true,
		"updated_from":  true,
		"updated_to":    true,
		"page":          true,
		"limit":         true,
		"sort_by":       true,
		"sort_order":    true,
	}

	query := r.URL.Query()

	for param := range query {
		paramLower := strings.ToLower(param)
		if !allowedParams[paramLower] {
			return errors.New("parâmetro desconhecido: '" + param + "'.")
		}
	}

	return nil
}

func (h *addressFilterHandler) parseFilterDTO(r *http.Request) (filterDTO.AddressFilterDTO, error) {
	var dto filterDTO.AddressFilterDTO
	query := r.URL.Query()

	if v := strings.TrimSpace(query.Get("user_id")); v != "" {
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return dto, errors.New("valor inválido para 'user_id': deve ser um número inteiro")
		}
		if parsed <= 0 {
			return dto, errors.New("valor inválido para 'user_id': deve ser maior que zero")
		}
		dto.UserID = &parsed
	}

	if v := strings.TrimSpace(query.Get("client_cpf_id")); v != "" {
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return dto, errors.New("valor inválido para 'client_cpf_id': deve ser um número inteiro")
		}
		if parsed <= 0 {
			return dto, errors.New("valor inválido para 'client_cpf_id': deve ser maior que zero")
		}
		dto.ClientCpfID = &parsed
	}

	if v := strings.TrimSpace(query.Get("supplier_id")); v != "" {
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return dto, errors.New("valor inválido para 'supplier_id': deve ser um número inteiro")
		}
		if parsed <= 0 {
			return dto, errors.New("valor inválido para 'supplier_id': deve ser maior que zero")
		}
		dto.SupplierID = &parsed
	}

	dto.Street = strings.TrimSpace(query.Get("street"))
	dto.StreetNumber = strings.TrimSpace(query.Get("street_number"))
	dto.Complement = strings.TrimSpace(query.Get("complement"))
	dto.City = strings.TrimSpace(query.Get("city"))
	dto.State = strings.TrimSpace(query.Get("state"))
	dto.Country = strings.TrimSpace(query.Get("country"))
	dto.PostalCode = strings.TrimSpace(query.Get("postal_code"))

	if v := strings.TrimSpace(query.Get("is_active")); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			return dto, errors.New("valor inválido para 'is_active': deve ser 'true' ou 'false'")
		}
		dto.IsActive = &parsed
	}

	dto.CreatedFrom = h.parseTimeParam(query, "created_from")
	dto.CreatedTo = h.parseTimeParam(query, "created_to")
	dto.UpdatedFrom = h.parseTimeParam(query, "updated_from")
	dto.UpdatedTo = h.parseTimeParam(query, "updated_to")

	if v := query.Get("created_from"); v != "" && dto.CreatedFrom == nil {
		return dto, errors.New("formato de data inválido para 'created_from': use formato RFC3339 (ex: 2024-01-01T00:00:00Z) ou YYYY-MM-DD")
	}
	if v := query.Get("created_to"); v != "" && dto.CreatedTo == nil {
		return dto, errors.New("formato de data inválido para 'created_to': use formato RFC3339 (ex: 2024-12-31T23:59:59Z) ou YYYY-MM-DD")
	}
	if v := query.Get("updated_from"); v != "" && dto.UpdatedFrom == nil {
		return dto, errors.New("formato de data inválido para 'updated_from': use formato RFC3339 ou YYYY-MM-DD")
	}
	if v := query.Get("updated_to"); v != "" && dto.UpdatedTo == nil {
		return dto, errors.New("formato de data inválido para 'updated_to': use formato RFC3339 ou YYYY-MM-DD")
	}

	dto.Limit, dto.Offset = utils.GetPaginationParams(r)
	dto.SortBy = strings.TrimSpace(query.Get("sort_by"))
	dto.SortOrder = strings.TrimSpace(query.Get("sort_order"))

	return dto, nil
}

func (h *addressFilterHandler) parseTimeParam(query map[string][]string, param string) *time.Time {
	values, exists := query[param]
	if !exists || len(values) == 0 {
		return nil
	}

	v := strings.TrimSpace(values[0])
	if v == "" {
		return nil
	}

	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		parsed, err := time.Parse(format, v)
		if err == nil {
			return &parsed
		}
	}

	return nil
}

func countFiltersApplied(dto filterDTO.AddressFilterDTO) int {
	count := 0

	if dto.UserID != nil {
		count++
	}
	if dto.ClientCpfID != nil {
		count++
	}
	if dto.SupplierID != nil {
		count++
	}
	if dto.Street != "" {
		count++
	}
	if dto.StreetNumber != "" {
		count++
	}
	if dto.Complement != "" {
		count++
	}
	if dto.City != "" {
		count++
	}
	if dto.State != "" {
		count++
	}
	if dto.Country != "" {
		count++
	}
	if dto.PostalCode != "" {
		count++
	}
	if dto.IsActive != nil {
		count++
	}
	if dto.CreatedFrom != nil {
		count++
	}
	if dto.CreatedTo != nil {
		count++
	}
	if dto.UpdatedFrom != nil {
		count++
	}
	if dto.UpdatedTo != nil {
		count++
	}

	return count
}
