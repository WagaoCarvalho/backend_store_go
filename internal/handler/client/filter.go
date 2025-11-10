package handler

import (
	"errors"
	"net/http"
	"strconv"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/client/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *Client) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const ref = "[clientHandler - GetAll] "

	var dtoFilter dto.ClientFilterDTO
	query := r.URL.Query()

	dtoFilter.Name = query.Get("name")
	dtoFilter.Email = query.Get("email")
	dtoFilter.CPF = query.Get("cpf")
	dtoFilter.CNPJ = query.Get("cnpj")

	if v := query.Get("status"); v != "" {
		if parsed, err := strconv.ParseBool(v); err == nil {
			dtoFilter.Status = &parsed
		}
	}

	if v := query.Get("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			dtoFilter.Limit = parsed
		}
	}

	if v := query.Get("offset"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			dtoFilter.Offset = parsed
		}
	}

	filter, _ := dtoFilter.ToModel() // conversão do DTO para o modelo

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"filtro": dtoFilter,
	})

	clients, err := h.service.GetAll(ctx, filter)
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

	clientDTOs := dto.ToClientDTOs(clients)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"total_encontrados": len(clientDTOs),
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Clientes listados com sucesso",
		Data:    clientDTOs,
	})
}
