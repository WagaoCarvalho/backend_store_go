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

	var dto dto_contact.ContactDTO
	if err := utils.FromJSON(r.Body, &dto); err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	contactModel := dto_contact.ToContactModel(dto)

	createdContact, err := h.service.Create(r.Context(), contactModel)
	if err != nil {
		if errors.Is(err, errMsg.ErrInvalidForeignKey) {
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	createdDTO := dto_contact.ToContactDTO(createdContact)

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Contato criado com sucesso",
		Data:    createdDTO,
	})
}

func (h *ContactHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[ContactHandler - GetByID] "

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	contactModel, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	contactDTO := dto_contact.ToContactDTO(contactModel)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Contato encontrado",
		Data:    contactDTO,
	})
}

func (h *ContactHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	contactModels, err := h.service.GetByUserID(r.Context(), id)
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
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	var dto dto_contact.ContactDTO
	if err := utils.FromJSON(r.Body, &dto); err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	contactModel := dto_contact.ToContactModel(dto)
	contactModel.ID = id

	if err := h.service.Update(r.Context(), contactModel); err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	updatedDTO := dto_contact.ToContactDTO(contactModel)
	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Contato atualizado com sucesso",
		Data:    updatedDTO,
	})
}

func (h *ContactHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
