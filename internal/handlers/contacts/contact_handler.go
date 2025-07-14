package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/contacts"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
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
		if errors.Is(err, repositories.ErrInvalidForeignKey) {
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
	ref := "[ContactHandler - GetByID] "

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), ref+"ID inválido na busca por ID", map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), ref+"Iniciando busca de contato por ID", map[string]any{
		"contact_id": id,
		"path":       r.URL.Path,
	})

	contact, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error(r.Context(), err, ref+"Erro ao buscar contato por ID", map[string]any{
			"contact_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	h.logger.Info(r.Context(), ref+"Contato recuperado com sucesso", map[string]any{
		"contact_id": id,
	})

	utils.ToJson(w, http.StatusOK, contact)
}

func (h *ContactHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	ref := "[ContactHandler - GetByUserID] "

	userID, err := utils.GetIDParam(r, "userID")
	if err != nil {
		h.logger.Warn(r.Context(), ref+"userID inválido na requisição", map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), ref+"Iniciando busca de contatos por usuário", map[string]any{
		"user_id": userID,
		"path":    r.URL.Path,
	})

	contacts, err := h.service.GetByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Error(r.Context(), err, ref+"Erro ao buscar contatos por usuário", map[string]any{
			"user_id": userID,
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), ref+"Contatos recuperados com sucesso", map[string]any{
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
	ref := "[ContactHandler - GetByClientID] "

	clientID, err := utils.GetIDParam(r, "clientID")
	if err != nil {
		h.logger.Warn(r.Context(), ref+"clientID inválido na requisição", map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), ref+"Iniciando busca de contatos por clientID", map[string]any{
		"client_id": clientID,
		"path":      r.URL.Path,
	})

	contacts, err := h.service.GetByClientID(r.Context(), clientID)
	if err != nil {
		h.logger.Error(r.Context(), err, ref+"Erro ao buscar contatos por clientID", map[string]any{
			"client_id": clientID,
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), ref+"Contatos recuperados com sucesso", map[string]any{
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
	ref := "[ContactHandler - GetBySupplierID] "

	supplierID, err := utils.GetIDParam(r, "supplierID")
	if err != nil {
		h.logger.Warn(r.Context(), ref+"supplierID inválido na requisição", map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), ref+"Iniciando busca de contatos por supplierID", map[string]any{
		"supplier_id": supplierID,
		"path":        r.URL.Path,
	})

	contacts, err := h.service.GetBySupplierID(r.Context(), supplierID)
	if err != nil {
		h.logger.Error(r.Context(), err, ref+"Erro ao buscar contatos por supplierID", map[string]any{
			"supplier_id": supplierID,
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), ref+"Contatos recuperados com sucesso", map[string]any{
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
	ref := "[ContactHandler - Update] "

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), ref+"ID inválido para atualização", map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	var contact models.Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		h.logger.Warn(r.Context(), ref+"JSON inválido no update", map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	contact.ID = id

	h.logger.Info(r.Context(), ref+"Iniciando atualização de contato", map[string]any{
		"contact_id": id,
		"path":       r.URL.Path,
	})

	if err := h.service.Update(r.Context(), &contact); err != nil {
		h.logger.Error(r.Context(), err, ref+"Erro ao atualizar contato", map[string]any{
			"contact_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), ref+"Contato atualizado com sucesso", map[string]any{
		"contact_id": id,
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Contato atualizado com sucesso",
		Data:    contact,
	})
}

func (h *ContactHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ref := "[ContactHandler - Delete] "

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), ref+"ID inválido para exclusão", map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), ref+"Iniciando exclusão de contato", map[string]any{
		"contact_id": id,
		"path":       r.URL.Path,
	})

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.logger.Error(r.Context(), err, ref+"Erro ao excluir contato", map[string]any{
			"contact_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	h.logger.Info(r.Context(), ref+"Contato excluído com sucesso", map[string]any{
		"contact_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}
