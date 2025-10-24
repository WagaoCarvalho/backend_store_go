package handler

import (
	"errors"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/contact_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/contact_relation"
)

type UserContactRelation struct {
	service service.UserContactRelation
	logger  *logger.LogAdapter
}

func NewUserContactRelation(service service.UserContactRelation, logger *logger.LogAdapter) *UserContactRelation {
	return &UserContactRelation{
		service: service,
		logger:  logger,
	}
}

func (h *UserContactRelation) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserContactRelationHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, nil)

	var requestData dto.UserContactRelationDTO
	if err := utils.FromJSON(r.Body, &requestData); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	modelRelation := dto.ToContactRelationModel(requestData)

	created, wasCreated, err := h.service.Create(ctx, modelRelation.UserID, modelRelation.ContactID)
	if err != nil {
		if errors.Is(err, errMsg.ErrDBInvalidForeignKey) {
			h.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"user_id":    modelRelation.UserID,
				"contact_id": modelRelation.ContactID,
				"erro":       err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"user_id":    modelRelation.UserID,
			"contact_id": modelRelation.ContactID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	status := http.StatusOK
	message := "Relação já existente"
	logMsg := logger.LogAlreadyExists

	if wasCreated {
		status = http.StatusCreated
		message = "Relação criada com sucesso"
		logMsg = logger.LogCreateSuccess
	}

	h.logger.Info(ctx, ref+logMsg, map[string]any{
		"user_id":    modelRelation.UserID,
		"contact_id": modelRelation.ContactID,
	})

	createdDTO := dto.ToContactRelationDTO(created)

	utils.ToJSON(w, status, utils.DefaultResponse{
		Data:    createdDTO,
		Message: message,
		Status:  status,
	})
}

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

// --- DELETE ---
func (h *UserContactRelation) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserContactRelationHandler - Delete] "
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

	if err := h.service.Delete(ctx, userID, contactID); err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao deletar relação", map[string]any{
			"user_id":    userID,
			"contact_id": contactID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+"Relação deletada com sucesso", map[string]any{
		"user_id":    userID,
		"contact_id": contactID,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    nil,
		Message: "Relação deletada com sucesso",
		Status:  http.StatusOK,
	})
}

// --- DELETE ALL ---
func (h *UserContactRelation) DeleteAll(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserContactRelationHandler - DeleteAll] "
	ctx := r.Context()

	userID, err := utils.GetIDParam(r, "user_id")
	if err != nil {
		h.logger.Warn(ctx, ref+"ID inválido", map[string]any{"user_id": userID})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteAll(ctx, userID); err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao deletar relações", map[string]any{
			"user_id": userID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+"Relações deletadas com sucesso", map[string]any{
		"user_id": userID,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    nil,
		Message: "Relações deletadas com sucesso",
		Status:  http.StatusOK,
	})
}
