package handler

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *clientCpfHandler) GetVersionByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[ClientHandler - GetVersionByID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, nil)
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	version, err := h.service.GetVersionByID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"client_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"client_id": id,
		"version":   version,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Versão do cliente recuperada com sucesso",
		Data: map[string]any{
			"client_id": id,
			"version":   version,
		},
	})
}
