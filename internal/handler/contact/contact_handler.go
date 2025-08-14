package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repositories/contacts"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/contacts"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	"github.com/WagaoCarvalho/backend_store_go/logger"
)

type ContactHandler struct {
	service services.ContactService
	logger  *logger.LoggerAdapter
}

func NewContactHandler(service services.ContactService, logger *logger.LoggerAdapter) *ContactHandler {
	return &ContactHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ContactHandler) Create(w http.ResponseWriter, r *http.Request) {
	ref := "[ContactHandler - Create] "
	var contact models.Contact

	h.logger.Info(r.Context(), ref+logger.LogCreateInit, map[string]any{})

	if err := utils.FromJson(r.Body, &contact); err != nil {
		h.logger.Warn(r.Context(), ref+logger.LogParseJsonError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	createdContact, err := h.service.Create(r.Context(), &contact)
	if err != nil {
		if errors.Is(err, repo.ErrInvalidForeignKey) {
			h.logger.Warn(r.Context(), ref+logger.LogForeignKeyViolation, map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		if err.Error() == "nome obrigatório" || err.Error() == "erro interno" {
			h.logger.Warn(r.Context(), ref+logger.LogValidateError, map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		h.logger.Error(r.Context(), err, ref+logger.LogCreateError, map[string]any{})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(r.Context(), ref+logger.LogCreateSuccess, map[string]any{
		"contact_id": createdContact.ID,
	})

	utils.ToJson(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Contato criado com sucesso",
		Data:    createdContact,
	})
}

func (h *ContactHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[ContactHandler - GetByID] "

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), ref+logger.LogGetInit, map[string]any{
		"contact_id": id,
		"path":       r.URL.Path,
	})

	contact, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error(r.Context(), err, ref+logger.LogGetError, map[string]any{
			"contact_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	h.logger.Info(r.Context(), ref+logger.LogGetSuccess, map[string]any{
		"contact_id": id,
	})

	utils.ToJson(w, http.StatusOK, contact)
}

func (h *ContactHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	const ref = "[ContactHandler - GetByUserID] "

	userID, err := utils.GetIDParam(r, "userID")
	if err != nil {
		h.logger.Warn(r.Context(), ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), ref+logger.LogGetInit, map[string]any{
		"user_id": userID,
		"path":    r.URL.Path,
	})

	contacts, err := h.service.GetByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Error(r.Context(), err, ref+logger.LogGetError, map[string]any{
			"user_id": userID,
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), ref+logger.LogGetSuccess, map[string]any{
		"user_id": userID,
		"count":   len(contacts),
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Contatos encontrados",
		Data:    contacts,
	})
}

func (h *ContactHandler) GetByClientID(w http.ResponseWriter, r *http.Request) {
	const ref = "[ContactHandler - GetByClientID] "

	clientID, err := utils.GetIDParam(r, "clientID")
	if err != nil {
		h.logger.Warn(r.Context(), ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), ref+logger.LogGetInit, map[string]any{
		"client_id": clientID,
		"path":      r.URL.Path,
	})

	contacts, err := h.service.GetByClientID(r.Context(), clientID)
	if err != nil {
		h.logger.Error(r.Context(), err, ref+logger.LogGetError, map[string]any{
			"client_id": clientID,
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), ref+logger.LogGetSuccess, map[string]any{
		"client_id": clientID,
		"count":     len(contacts),
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Contatos encontrados",
		Data:    contacts,
	})
}

func (h *ContactHandler) GetBySupplierID(w http.ResponseWriter, r *http.Request) {
	const ref = "[ContactHandler - GetBySupplierID] "

	supplierID, err := utils.GetIDParam(r, "supplierID")
	if err != nil {
		h.logger.Warn(r.Context(), ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), ref+logger.LogGetInit, map[string]any{
		"supplier_id": supplierID,
		"path":        r.URL.Path,
	})

	contacts, err := h.service.GetBySupplierID(r.Context(), supplierID)
	if err != nil {
		h.logger.Error(r.Context(), err, ref+logger.LogGetError, map[string]any{
			"supplier_id": supplierID,
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), ref+logger.LogGetSuccess, map[string]any{
		"supplier_id": supplierID,
		"count":       len(contacts),
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Contatos encontrados",
		Data:    contacts,
	})
}

func (h *ContactHandler) Update(w http.ResponseWriter, r *http.Request) {
	const ref = "[ContactHandler - Update] "

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	var contact models.Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		h.logger.Warn(r.Context(), ref+logger.LogParseJsonError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	contact.ID = id

	h.logger.Info(r.Context(), ref+logger.LogUpdateInit, map[string]any{
		"contact_id": id,
		"path":       r.URL.Path,
	})

	if err := h.service.Update(r.Context(), &contact); err != nil {
		h.logger.Error(r.Context(), err, ref+logger.LogUpdateError, map[string]any{
			"contact_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), ref+logger.LogUpdateSuccess, map[string]any{
		"contact_id": id,
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Contato atualizado com sucesso",
		Data:    contact,
	})
}

func (h *ContactHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[ContactHandler - Delete] "

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), ref+logger.LogDeleteInit, map[string]any{
		"contact_id": id,
		"path":       r.URL.Path,
	})

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.logger.Error(r.Context(), err, ref+logger.LogDeleteError, map[string]any{
			"contact_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	h.logger.Info(r.Context(), ref+logger.LogDeleteSuccess, map[string]any{
		"contact_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}
