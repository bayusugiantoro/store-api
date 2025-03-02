package handler

import (
	"api-otto/internal/domain"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type VoucherHandler struct {
    service domain.VoucherService
}

func NewVoucherHandler(service domain.VoucherService) *VoucherHandler {
    return &VoucherHandler{service: service}
}

func (h *VoucherHandler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    var voucher domain.Voucher
    if err := json.NewDecoder(r.Body).Decode(&voucher); err != nil {
        writeError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    if err := h.service.Create(&voucher); err != nil {
        writeError(w, http.StatusInternalServerError, err.Error())
        return
    }

    resp := Response{
        Status:  http.StatusCreated,
        Message: "Voucher created successfully",
        Data:    voucher,
    }
    writeJSON(w, http.StatusCreated, resp)
}

func (h *VoucherHandler) GetByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    idStr := ps.ByName("id")
    if idStr == "" {
        writeError(w, http.StatusBadRequest, "ID is required")
        return
    }

    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        writeError(w, http.StatusBadRequest, "Invalid ID")
        return
    }

    voucher, err := h.service.GetByID(id)
    if err != nil {
        writeError(w, http.StatusInternalServerError, err.Error())
        return
    }
    if voucher == nil {
        writeError(w, http.StatusNotFound, "Voucher not found")
        return
    }

    resp := Response{
        Status:  http.StatusOK,
        Message: "Success",
        Data:    voucher,
    }
    writeJSON(w, http.StatusOK, resp)
}

func (h *VoucherHandler) GetByBrandID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    idStr := ps.ByName("id")
    if idStr == "" {
        writeError(w, http.StatusBadRequest, "Brand ID is required")
        return
    }

    brandID, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        writeError(w, http.StatusBadRequest, "Invalid brand ID")
        return
    }

    vouchers, err := h.service.GetByBrandID(brandID)
    if err != nil {
        if err.Error() == "brand not found" {
            writeError(w, http.StatusNotFound, "Brand not found")
            return
        }
        writeError(w, http.StatusInternalServerError, err.Error())
        return
    }

    if len(vouchers) == 0 {
        writeError(w, http.StatusNotFound, "Brand tidak memiliki voucher")
        return
    }

    resp := Response{
        Status:  http.StatusOK,
        Message: "Success",
        Data:    vouchers,
    }
    writeJSON(w, http.StatusOK, resp)
} 