package handler

import (
	"errors"
	"net/http"
	"strconv"

	dtoFilter "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/filter"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

// Lista de parâmetros válidos para validação
var validFilterParams = map[string]bool{
	"username":     true,
	"email":        true,
	"status":       true,
	"created_from": true,
	"created_to":   true,
	"updated_from": true,
	"updated_to":   true,
	"limit":        true,
	"offset":       true,
}

func (h *userFilterHandler) Filter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const ref = "[userHandler - Filter] "

	query := r.URL.Query()

	// Validação 1: Verificar parâmetros desconhecidos
	for param := range query {
		if !validFilterParams[param] {
			utils.ErrorResponse(w,
				errors.New("parâmetro de consulta inválido: "+param),
				http.StatusBadRequest)
			return
		}
	}

	var dtoFilter dtoFilter.UserFilterDTO

	dtoFilter.Username = query.Get("username")
	dtoFilter.Email = query.Get("email")

	// Validação 2: Status com valor inválido deve retornar erro
	if v := query.Get("status"); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			utils.ErrorResponse(w,
				errors.New("status deve ser true ou false"),
				http.StatusBadRequest)
			return
		}
		dtoFilter.Status = &parsed
	}

	// Datas via utilitário (assumindo que ParseTimeRange já valida formato)
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

	users, err := h.service.Filter(ctx, filter)
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

	userDTOs := dto.ToUserDTOs(users)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"total_encontrados": len(userDTOs),
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuários listados com sucesso",
		Data: map[string]any{
			"total": len(userDTOs),
			"items": userDTOs,
		},
	})
}
