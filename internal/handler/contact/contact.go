package handler

import (
	"errors"
	"net/http"

	dto_contact "github.com/WagaoCarvalho/backend_store_go/internal/dto/contact"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/contact"
)

type ContactHandler struct {
	service service.ContactService
	logger  *logger.LogAdapter
}

func NewContactHandler(service service.ContactService, logger *logger.LogAdapter) *ContactHandler {
	return &ContactHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ContactHandler) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[ContactHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, nil)

	var contactDTO dto_contact.ContactDTO
	if err := utils.FromJSON(r.Body, &contactDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	contactModel := dto_contact.ToContactModel(contactDTO)

	createdContact, err := h.service.Create(ctx, contactModel)
	if err != nil {
		if errors.Is(err, errMsg.ErrInvalidForeignKey) {
			h.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogCreateError, nil)
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	createdDTO := dto_contact.ToContactDTO(createdContact)

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"contact_id": createdDTO.ID,
	})

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Contato criado com sucesso",
		Data:    createdDTO,
	})
}

func (h *ContactHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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

	contactDTO := dto_contact.ToContactDTO(contactModel)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"contact_id": contactDTO.ID,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Contato encontrado",
		Data:    contactDTO,
	})
}

func (h *ContactHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	const ref = "[ContactHandler - GetByUserID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, nil)

	userID, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	contactModels, err := h.service.GetByUserID(ctx, userID)
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
			"user_id": userID,
			"erro":    err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	contactDTOs := make([]dto_contact.ContactDTO, len(contactModels))
	for i, c := range contactModels {
		contactDTOs[i] = dto_contact.ToContactDTO(c)
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"user_id": userID,
		"count":   len(contactDTOs),
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Contatos do usuário encontrados",
		Data:    contactDTOs,
	})
}

func (h *ContactHandler) GetByClientID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	contactModels, err := h.service.GetByClientID(r.Context(), id)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	contactDTOs := make([]dto_contact.ContactDTO, len(contactModels))
	for i, c := range contactModels {
		contactDTOs[i] = dto_contact.ToContactDTO(c)
	}

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Contatos do cliente encontrados",
		Data:    contactDTOs,
	})
}

func (h *ContactHandler) GetBySupplierID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	contactModels, err := h.service.GetBySupplierID(r.Context(), id)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	contactDTOs := make([]dto_contact.ContactDTO, len(contactModels))
	for i, c := range contactModels {
		contactDTOs[i] = dto_contact.ToContactDTO(c)
	}

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Contatos do fornecedor encontrados",
		Data:    contactDTOs,
	})
}

func (h *ContactHandler) Update(w http.ResponseWriter, r *http.Request) {
	const ref = "[ContactHandler - Update] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogUpdateInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	var dto dto_contact.ContactDTO
	if err := utils.FromJSON(r.Body, &dto); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	contactModel := dto_contact.ToContactModel(dto)
	contactModel.ID = id

	if err := h.service.Update(ctx, contactModel); err != nil {
		h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"contact_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	updatedDTO := dto_contact.ToContactDTO(contactModel)
	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"contact_id": updatedDTO.ID,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Contato atualizado com sucesso",
		Data:    updatedDTO,
	})
}

func (h *ContactHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[ContactHandler - Delete] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogDeleteInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		h.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"contact_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"contact_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}
