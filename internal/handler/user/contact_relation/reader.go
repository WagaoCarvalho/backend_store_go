package handler

import (
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/contact_relation"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *UserContactRelation) GetAllByUserID(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserContactRelationHandler - GetAllByUserID] "
	ctx := r.Context()

	userID, err := utils.GetIDParam(r, "user_id")
	if err != nil {
		h.logger.Warn(ctx, ref+"ID inválido", map[string]any{"user_id": userID})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	relations, err := h.service.GetAllRelationsByUserID(ctx, userID)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao buscar relações", map[string]any{"user_id": userID})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+"Relações retornadas com sucesso", map[string]any{"user_id": userID})

	relationsDTO := dto.ToUserContactRelationsDTOs(relations)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    relationsDTO,
		Message: "Relações encontradas",
		Status:  http.StatusOK,
	})
}

// --- HAS RELATION ---
func (h *UserContactRelation) HasRelation(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserContactRelationHandler - HasRelation] "
	ctx := r.Context()

	userID, err := utils.GetIDParam(r, "user_id")
	if err != nil {
		h.logger.Warn(ctx, ref+"ID inválido", map[string]any{"user_id": userID})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	contactID, err := utils.GetIDParam(r, "contact_id")
	if err != nil {
		h.logger.Warn(ctx, ref+"ID inválido", map[string]any{"contact_id": contactID})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	exists, err := h.service.HasUserContactRelation(ctx, userID, contactID)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao verificar relação", map[string]any{
			"user_id":    userID,
			"contact_id": contactID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+"Verificação concluída", map[string]any{
		"user_id":    userID,
		"contact_id": contactID,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    map[string]bool{"exists": exists},
		Message: "Verificação concluída",
		Status:  http.StatusOK,
	})
}
