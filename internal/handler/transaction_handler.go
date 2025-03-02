package handler

import (
	"api-otto/internal/domain"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type TransactionHandler struct {
    service domain.TransactionService
}

func NewTransactionHandler(service domain.TransactionService) *TransactionHandler {
    return &TransactionHandler{service: service}
}

func (h *TransactionHandler) CreateRedemption(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    var transaction domain.Transaction
    if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
        writeError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    if err := h.service.CreateRedemption(&transaction); err != nil {
        writeError(w, http.StatusInternalServerError, err.Error())
        return
    }

    resp := Response{
        Status:  http.StatusCreated,
        Message: "Redemption created successfully",
        Data:    transaction,
    }
    writeJSON(w, http.StatusCreated, resp)
}

func (h *TransactionHandler) GetTransactionByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    idStr := ps.ByName("id")
    if idStr == "" {
        writeError(w, http.StatusBadRequest, "Transaction ID is required")
        return
    }

    transactionID, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        writeError(w, http.StatusBadRequest, "Invalid transaction ID")
        return
    }

    transaction, err := h.service.GetTransactionByID(transactionID)
    if err != nil {
        if err.Error() == "transaction not found" {
            writeError(w, http.StatusNotFound, "Transaction not found")
            return
        }
        writeError(w, http.StatusInternalServerError, err.Error())
        return
    }

    if transaction == nil {
        writeError(w, http.StatusNotFound, "Transaction not found")
        return
    }

    resp := Response{
        Status:  http.StatusOK,
        Message: "Success",
        Data:    transaction,
    }
    writeJSON(w, http.StatusOK, resp)
} 