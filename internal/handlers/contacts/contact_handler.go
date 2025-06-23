package handlers

import (
	"encoding/json"
	"net/http"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/contacts"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
)

type ContactHandler struct {
	service services.ContactService
}

func NewContactHandler(service services.ContactService) *ContactHandler {
	return &ContactHandler{
		service: service,
	}
}

func (h *ContactHandler) Create(w http.ResponseWriter, r *http.Request) {
	var contact models.Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	createdContact, err := h.service.Create(r.Context(), &contact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdContact)
}

func (h *ContactHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	contact, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(contact)
}

func (h *ContactHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetIDParam(r, "userID")
	if err != nil {
		http.Error(w, "userID inválido", http.StatusBadRequest)
		return
	}

	contacts, err := h.service.GetByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(contacts)
}

func (h *ContactHandler) GetByClientID(w http.ResponseWriter, r *http.Request) {
	clientID, err := utils.GetIDParam(r, "clientID")
	if err != nil {
		http.Error(w, "clientID inválido", http.StatusBadRequest)
		return
	}

	contacts, err := h.service.GetByClientID(r.Context(), clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(contacts)
}

func (h *ContactHandler) GetBySupplierID(w http.ResponseWriter, r *http.Request) {
	supplierID, err := utils.GetIDParam(r, "supplierID")
	if err != nil {
		http.Error(w, "supplierID inválido", http.StatusBadRequest)
		return
	}

	contacts, err := h.service.GetBySupplierID(r.Context(), supplierID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(contacts)
}

func (h *ContactHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var contact models.Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	contact.ID = id

	if err := h.service.Update(r.Context(), &contact); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(contact)
}

func (h *ContactHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
