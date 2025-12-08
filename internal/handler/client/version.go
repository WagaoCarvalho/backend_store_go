package handler

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *clientHandler) GetVersionByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[clientHandler - GetVersionByID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{})

	uid, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	version, err := h.service.GetVersionByID(ctx, uid)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"client_id": uid,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"client_id": uid,
		"version":   version,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Versão do cliente recuperada com sucesso",
		Data: map[string]any{
			"client_id": uid,
			"version":   version,
		},
	})
}
