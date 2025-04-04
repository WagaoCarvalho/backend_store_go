package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/contacts"
	"github.com/WagaoCarvalho/backend_store_go/utils"
)

type ContactService interface {
	CreateContact(ctx context.Context, contact *models.Contact) error
	GetContactByID(ctx context.Context, id int64) (*models.Contact, error)
	ListContactsByUser(ctx context.Context, userID int64) ([]*models.Contact, error)
	ListContactsByClient(ctx context.Context, clientID int64) ([]*models.Contact, error)
	ListContactsBySupplier(ctx context.Context, supplierID int64) ([]*models.Contact, error)
	UpdateContact(ctx context.Context, contact *models.Contact) error
	DeleteContact(ctx context.Context, id int64) error
}

type ContactHandler struct {
	service services.ContactService
}

func NewContactHandler(service services.ContactService) *ContactHandler {
	return &ContactHandler{
		service: service,
	}
}

// POST /contacts
func (h *ContactHandler) CreateContact(w http.ResponseWriter, r *http.Request) {
	var contact models.Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateContact(r.Context(), &contact); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(contact)
}

// GET /contacts/{id:[0-9]+}
func (h *ContactHandler) GetContactByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	contact, err := h.service.GetContactByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(contact)
}

// GET /contacts/user/{userID:[0-9]+}
func (h *ContactHandler) ListContactsByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetIDParam(r, "userID")
	if err != nil {
		http.Error(w, "userID inválido", http.StatusBadRequest)
		return
	}

	contacts, err := h.service.ListContactsByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(contacts)
}

// GET /contacts/client/{clientID:[0-9]+}
func (h *ContactHandler) ListContactsByClient(w http.ResponseWriter, r *http.Request) {
	clientID, err := utils.GetIDParam(r, "clientID")
	if err != nil {
		http.Error(w, "clientID inválido", http.StatusBadRequest)
		return
	}

	contacts, err := h.service.ListContactsByClient(r.Context(), clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(contacts)
}

// GET /contacts/supplier/{supplierID:[0-9]+}
func (h *ContactHandler) ListContactsBySupplier(w http.ResponseWriter, r *http.Request) {
	supplierID, err := utils.GetIDParam(r, "supplierID")
	if err != nil {
		http.Error(w, "supplierID inválido", http.StatusBadRequest)
		return
	}

	contacts, err := h.service.ListContactsBySupplier(r.Context(), supplierID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(contacts)
}

// PUT /contacts/{id:[0-9]+}
func (h *ContactHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
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

	if err := h.service.UpdateContact(r.Context(), &contact); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(contact)
}

// DELETE /contacts/{id:[0-9]+}
func (h *ContactHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteContact(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
