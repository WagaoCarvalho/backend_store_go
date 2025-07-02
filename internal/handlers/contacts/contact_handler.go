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
	var contact models.Contact

	h.logger.Info(r.Context(), "[ContactHandler] - Iniciando criação de contato", map[string]interface{}{})

	if err := utils.FromJson(r.Body, &contact); err != nil {
		h.logger.Warn(r.Context(), "[ContactHandler] - Falha ao fazer parse do JSON", map[string]interface{}{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	createdContact, err := h.service.Create(r.Context(), &contact)
	if err != nil {
		if errors.Is(err, repositories.ErrInvalidForeignKey) {
			h.logger.Warn(r.Context(), "[ContactHandler] - Foreign key inválida", map[string]interface{}{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		if err.Error() == "nome obrigatório" || err.Error() == "erro interno" {
			h.logger.Warn(r.Context(), "[ContactHandler] - Erro de validação ou negócio", map[string]interface{}{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		h.logger.Error(r.Context(), err, "[ContactHandler] - Erro ao criar contato", map[string]interface{}{})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(r.Context(), "[ContactHandler] - Contato criado com sucesso", map[string]interface{}{
		"contact_id": createdContact.ID,
	})

	utils.ToJson(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Contato criado com sucesso",
		Data:    createdContact,
	})
}

func (h *ContactHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), "[ContactHandler] - ID inválido na busca por ID", map[string]interface{}{
			"erro": err.Error(),
		})
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), "[ContactHandler] - Iniciando busca de contato por ID", map[string]interface{}{
		"contact_id": id,
		"path":       r.URL.Path,
	})

	contact, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error(r.Context(), err, "[ContactHandler] - Erro ao buscar contato por ID", map[string]interface{}{
			"contact_id": id,
		})
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.logger.Info(r.Context(), "[ContactHandler] - Contato recuperado com sucesso", map[string]interface{}{
		"contact_id": id,
	})

	json.NewEncoder(w).Encode(contact)
}

func (h *ContactHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetIDParam(r, "userID")
	if err != nil {
		h.logger.Warn(r.Context(), "[ContactHandler] - userID inválido na busca por usuário", map[string]interface{}{
			"erro": err.Error(),
		})
		http.Error(w, "userID inválido", http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), "[ContactHandler] - Iniciando busca de contatos por usuário", map[string]interface{}{
		"user_id": userID,
		"path":    r.URL.Path,
	})

	contacts, err := h.service.GetByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Error(r.Context(), err, "[ContactHandler] - Erro ao buscar contatos por usuário", map[string]interface{}{
			"user_id": userID,
		})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), "[ContactHandler] - Contatos recuperados com sucesso", map[string]interface{}{
		"user_id": userID,
		"count":   len(contacts),
	})

	json.NewEncoder(w).Encode(contacts)
}

func (h *ContactHandler) GetByClientID(w http.ResponseWriter, r *http.Request) {
	clientID, err := utils.GetIDParam(r, "clientID")
	if err != nil {
		h.logger.Warn(r.Context(), "[ContactHandler] - clientID inválido na busca por clientID", map[string]interface{}{
			"erro": err.Error(),
		})
		http.Error(w, "clientID inválido", http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), "[ContactHandler] - Iniciando busca de contatos por clientID", map[string]interface{}{
		"client_id": clientID,
		"path":      r.URL.Path,
	})

	contacts, err := h.service.GetByClientID(r.Context(), clientID)
	if err != nil {
		h.logger.Error(r.Context(), err, "[ContactHandler] - Erro ao buscar contatos por clientID", map[string]interface{}{
			"client_id": clientID,
		})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), "[ContactHandler] - Contatos recuperados com sucesso", map[string]interface{}{
		"client_id": clientID,
		"count":     len(contacts),
	})

	json.NewEncoder(w).Encode(contacts)
}

func (h *ContactHandler) GetBySupplierID(w http.ResponseWriter, r *http.Request) {
	supplierID, err := utils.GetIDParam(r, "supplierID")
	if err != nil {
		h.logger.Warn(r.Context(), "[ContactHandler] - supplierID inválido na busca por supplierID", map[string]interface{}{
			"erro": err.Error(),
		})
		http.Error(w, "supplierID inválido", http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), "[ContactHandler] - Iniciando busca de contatos por supplierID", map[string]interface{}{
		"supplier_id": supplierID,
		"path":        r.URL.Path,
	})

	contacts, err := h.service.GetBySupplierID(r.Context(), supplierID)
	if err != nil {
		h.logger.Error(r.Context(), err, "[ContactHandler] - Erro ao buscar contatos por supplierID", map[string]interface{}{
			"supplier_id": supplierID,
		})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), "[ContactHandler] - Contatos recuperados com sucesso", map[string]interface{}{
		"supplier_id": supplierID,
		"count":       len(contacts),
	})

	json.NewEncoder(w).Encode(contacts)
}

func (h *ContactHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), "[ContactHandler] - ID inválido para atualização", map[string]interface{}{
			"erro": err.Error(),
		})
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var contact models.Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		h.logger.Warn(r.Context(), "[ContactHandler] - JSON inválido no update", map[string]interface{}{
			"erro": err.Error(),
		})
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	contact.ID = id

	h.logger.Info(r.Context(), "[ContactHandler] - Iniciando atualização de contato", map[string]interface{}{
		"contact_id": id,
		"path":       r.URL.Path,
	})

	if err := h.service.Update(r.Context(), &contact); err != nil {
		h.logger.Error(r.Context(), err, "[ContactHandler] - Erro ao atualizar contato", map[string]interface{}{
			"contact_id": id,
		})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), "[ContactHandler] - Contato atualizado com sucesso", map[string]interface{}{
		"contact_id": id,
	})

	json.NewEncoder(w).Encode(contact)
}

func (h *ContactHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), "[ContactHandler] - ID inválido para exclusão", map[string]interface{}{
			"erro": err.Error(),
		})
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), "[ContactHandler] - Iniciando exclusão de contato", map[string]interface{}{
		"contact_id": id,
		"path":       r.URL.Path,
	})

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.logger.Error(r.Context(), err, "[ContactHandler] - Erro ao excluir contato", map[string]interface{}{
			"contact_id": id,
		})
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.logger.Info(r.Context(), "[ContactHandler] - Contato excluído com sucesso", map[string]interface{}{
		"contact_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}
