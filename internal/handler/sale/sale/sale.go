package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	dtoSale "github.com/WagaoCarvalho/backend_store_go/internal/dto/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/sale/sale"
	"github.com/gorilla/mux"
)

type SaleHandler struct {
	service service.SaleService
	logger  *logger.LogAdapter
}

func NewSaleHandler(service service.SaleService, logger *logger.LogAdapter) *SaleHandler {
	return &SaleHandler{
		service: service,
		logger:  logger,
	}
}

func (h *SaleHandler) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, nil)

	var saleDTO dtoSale.SaleDTO
	if err := utils.FromJSON(r.Body, &saleDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	saleModel := dtoSale.ToSaleModel(saleDTO)

	createdModel, err := h.service.Create(ctx, saleModel)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogCreateError, nil)
		if errors.Is(err, errMsg.ErrInvalidForeignKey) {
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	createdDTO := dtoSale.ToSaleDTO(createdModel)

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{"sale_id": createdDTO.ID})
	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Venda criada com sucesso",
		Data:    createdDTO,
	})
}

func (h *SaleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleHandler - GetByID] "
	ctx := r.Context()

	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		h.logger.Warn(ctx, ref+"ID inválido", map[string]any{"id": idStr})
		utils.ErrorResponse(w, errMsg.ErrID, http.StatusBadRequest)
		return
	}

	sale, err := h.service.GetByID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao buscar venda", map[string]any{"id": id})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Venda encontrada",
		Data:    dtoSale.ToSaleDTO(sale),
	})
}

func (h *SaleHandler) GetByClientID(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleHandler - GetByClientID] "
	ctx := r.Context()

	clientIDStr := r.URL.Query().Get("client_id")
	clientID, err := strconv.ParseInt(clientIDStr, 10, 64)
	if err != nil || clientID <= 0 {
		h.logger.Warn(ctx, ref+"clientID inválido", map[string]any{"client_id": clientIDStr})
		utils.ErrorResponse(w, errMsg.ErrID, http.StatusBadRequest)
		return
	}

	limit, offset := utils.ParseLimitOffset(r)
	orderBy, orderDir := utils.ParseOrder(r)

	sales, err := h.service.GetByClientID(ctx, clientID, limit, offset, orderBy, orderDir)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao buscar vendas por clientID", map[string]any{"client_id": clientID})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	salesDTO := dtoSale.ToSaleDTOList(sales)
	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Vendas do cliente recuperadas",
		Data:    salesDTO,
	})
}

func (h *SaleHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleHandler - GetByUserID] "
	ctx := r.Context()

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		h.logger.Warn(ctx, ref+"userID inválido", map[string]any{"user_id": userIDStr})
		utils.ErrorResponse(w, errMsg.ErrID, http.StatusBadRequest)
		return
	}

	limit, offset := utils.ParseLimitOffset(r)
	orderBy, orderDir := utils.ParseOrder(r)

	sales, err := h.service.GetByUserID(ctx, userID, limit, offset, orderBy, orderDir)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao buscar vendas por userID", map[string]any{"user_id": userID})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	salesDTO := dtoSale.ToSaleDTOList(sales)
	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Vendas do usuário recuperadas",
		Data:    salesDTO,
	})
}

func (h *SaleHandler) GetByStatus(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleHandler - GetByStatus] "
	ctx := r.Context()

	status := r.URL.Query().Get("status")
	if status == "" {
		h.logger.Warn(ctx, ref+"status vazio", map[string]any{"status": status})
		utils.ErrorResponse(w, errMsg.ErrInvalidData, http.StatusBadRequest)
		return
	}

	limit, offset := utils.ParseLimitOffset(r)
	orderBy, orderDir := utils.ParseOrder(r)

	sales, err := h.service.GetByStatus(ctx, status, limit, offset, orderBy, orderDir)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao buscar vendas por status", map[string]any{"status": status})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	salesDTO := dtoSale.ToSaleDTOList(sales)
	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Vendas por status recuperadas",
		Data:    salesDTO,
	})
}

func (h *SaleHandler) GetByDateRange(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleHandler - GetByDateRange] "
	ctx := r.Context()

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	start, errStart := time.Parse(time.RFC3339, startStr)
	end, errEnd := time.Parse(time.RFC3339, endStr)
	if errStart != nil || errEnd != nil {
		h.logger.Warn(ctx, ref+"datas inválidas", map[string]any{"start": startStr, "end": endStr})
		utils.ErrorResponse(w, errMsg.ErrInvalidData, http.StatusBadRequest)
		return
	}

	limit, offset := utils.ParseLimitOffset(r)
	orderBy, orderDir := utils.ParseOrder(r)

	sales, err := h.service.GetByDateRange(ctx, start, end, limit, offset, orderBy, orderDir)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao buscar vendas por período", map[string]any{"start": start, "end": end})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	salesDTO := dtoSale.ToSaleDTOList(sales)
	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Vendas por período recuperadas",
		Data:    salesDTO,
	})
}

func (h *SaleHandler) Update(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleHandler - Update] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+"Iniciando atualização da venda", nil)

	var saleDTO dtoSale.SaleDTO
	if err := utils.FromJSON(r.Body, &saleDTO); err != nil {
		h.logger.Warn(ctx, ref+"Erro ao parsear JSON", map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		utils.ErrorResponse(w, errMsg.ErrID, http.StatusBadRequest)
		return
	}
	saleDTO.ID = &id
	saleModel := dtoSale.ToSaleModel(saleDTO)

	if err := h.service.Update(ctx, saleModel); err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao atualizar venda", nil)
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+"Venda atualizada com sucesso", map[string]any{"sale_id": id})
	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Venda atualizada com sucesso",
	})
}

func (h *SaleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleHandler - Delete] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+"Iniciando exclusão da venda", nil)

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		utils.ErrorResponse(w, errMsg.ErrID, http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao excluir venda", nil)
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+"Venda excluída com sucesso", map[string]any{"sale_id": id})
	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Venda excluída com sucesso",
	})
}
