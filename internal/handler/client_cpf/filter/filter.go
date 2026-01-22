package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/client_cpf/client"
	dtoclientCpfFilter "github.com/WagaoCarvalho/backend_store_go/internal/dto/client_cpf/filter"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *clientCpfFilterHandler) Filter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const ref = "[clientHandler - Filter] "

	// 1. Validar query parameters
	if err := h.validateQueryParams(r); err != nil {
		h.logger.Warn(ctx, ref+"validação de parâmetros falhou", map[string]any{
			"erro":  err.Error(),
			"query": r.URL.Query(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// 2. Parsear filtros
	dtoFilter, err := h.parseFilterDTO(r)
	if err != nil {
		h.logger.Warn(ctx, ref+"erro ao parsear filtros", map[string]any{
			"erro":  err.Error(),
			"query": r.URL.Query(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// 3. Validar DTO
	if err := dtoFilter.Validate(); err != nil {
		h.logger.Warn(ctx, ref+"validação de DTO falhou", map[string]any{
			"erro": err.Error(),
			"dto":  dtoFilter,
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// 4. Converter para modelo
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

	// 5. Executar filtro
	clients, err := h.service.Filter(ctx, filter)
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

	// 6. Preparar resposta
	clientDTOs := dto.ToClientCpfDTOs(clients)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"total_encontrados": len(clientDTOs),
		"filtros_aplicados": countFiltersApplied(dtoFilter),
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Clientes listados com sucesso",
		Data: map[string]any{
			"total":           len(clientDTOs),
			"items":           clientDTOs,
			"filters_applied": countFiltersApplied(dtoFilter),
			"has_more":        len(clientDTOs) == dtoFilter.Limit,
		},
	})
}

func (h *clientCpfFilterHandler) validateQueryParams(r *http.Request) error {
	allowedParams := map[string]bool{
		"name":         true,
		"email":        true,
		"cpf":          true,
		"cnpj":         true,
		"status":       true,
		"version":      true,
		"created_from": true,
		"created_to":   true,
		"updated_from": true,
		"updated_to":   true,
		"page":         true,
		"limit":        true,
		"sort_by":      true,
		"sort_order":   true,
	}

	query := r.URL.Query()

	// Verificar parâmetros desconhecidos
	for param := range query {
		paramLower := strings.ToLower(param)
		if !allowedParams[paramLower] {
			return errors.New("parâmetro desconhecido: '" + param + "'.")
		}
	}

	return nil
}

func (h *clientCpfFilterHandler) parseFilterDTO(r *http.Request) (dtoclientCpfFilter.ClientFilterDTO, error) {
	var dto dtoclientCpfFilter.ClientFilterDTO
	query := r.URL.Query()

	// Campos de texto
	dto.Name = strings.TrimSpace(query.Get("name"))
	dto.Email = strings.TrimSpace(query.Get("email"))
	dto.CPF = strings.TrimSpace(query.Get("cpf"))
	dto.CNPJ = strings.TrimSpace(query.Get("cnpj"))

	// Status (boolean)
	if v := strings.TrimSpace(query.Get("status")); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			return dto, errors.New("valor inválido para 'status': deve ser 'true' ou 'false'")
		}
		dto.Status = &parsed
	}

	// Version (integer)
	if v := strings.TrimSpace(query.Get("version")); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return dto, errors.New("valor inválido para 'version': deve ser um número inteiro")
		}
		if parsed <= 0 {
			return dto, errors.New("valor inválido para 'version': deve ser maior que zero")
		}
		dto.Version = &parsed
	}

	// Datas
	dto.CreatedFrom = h.parseTimeParam(query, "created_from")
	dto.CreatedTo = h.parseTimeParam(query, "created_to")
	dto.UpdatedFrom = h.parseTimeParam(query, "updated_from")
	dto.UpdatedTo = h.parseTimeParam(query, "updated_to")

	// Validar formatos de data (se algum foi fornecido)
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

	// Paginação e ordenação
	dto.Limit, dto.Offset = utils.GetPaginationParams(r)
	dto.SortBy = strings.TrimSpace(query.Get("sort_by"))
	dto.SortOrder = strings.TrimSpace(query.Get("sort_order"))

	return dto, nil
}

func (h *clientCpfFilterHandler) parseTimeParam(query map[string][]string, param string) *time.Time {
	values, exists := query[param]
	if !exists || len(values) == 0 {
		return nil
	}

	v := strings.TrimSpace(values[0])
	if v == "" {
		return nil
	}

	// Tentar diferentes formatos de data
	formats := []string{
		time.RFC3339,          // "2006-01-02T15:04:05Z07:00"
		"2006-01-02T15:04:05", // Sem timezone
		"2006-01-02 15:04:05", // MySQL datetime format
		"2006-01-02",          // Date only
	}

	for _, format := range formats {
		parsed, err := time.Parse(format, v)
		if err == nil {
			return &parsed
		}
	}

	return nil
}

// Função auxiliar para contar filtros aplicados
func countFiltersApplied(dto dtoclientCpfFilter.ClientFilterDTO) int {
	count := 0

	if dto.Name != "" {
		count++
	}
	if dto.Email != "" {
		count++
	}
	if dto.CPF != "" {
		count++
	}
	if dto.CNPJ != "" {
		count++
	}
	if dto.Status != nil {
		count++
	}
	if dto.Version != nil {
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
