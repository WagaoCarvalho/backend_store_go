package handler

import (
	"errors"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/client_cpf/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *clientCpfHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[ClientHandler - GetByID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	clientModel, err := h.service.GetByID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"client_id": id,
		})

		status := http.StatusInternalServerError
		if errors.Is(err, errMsg.ErrNotFound) {
			status = http.StatusNotFound
		}

		utils.ErrorResponse(w, err, status)
		return
	}

	clientDTO := dto.ToClientCpfDTO(clientModel)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"client_id": clientDTO.ID,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Cliente encontrado",
		Data:    clientDTO,
	})
}
