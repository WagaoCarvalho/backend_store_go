package handler

import (
	"net/http"

	dtoContact "github.com/WagaoCarvalho/backend_store_go/internal/dto/contact"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *contactHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[ContactHandler - GetByID] "
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

	contactModel, err := h.service.GetByID(ctx, id)
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
			"contact_id": id,
			"erro":       err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	contactDTO := dtoContact.ToContactDTO(contactModel)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"contact_id": contactDTO.ID,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Contato encontrado",
		Data:    contactDTO,
	})
}
