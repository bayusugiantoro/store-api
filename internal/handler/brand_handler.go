package handler

import (
	"api-otto/internal/domain"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"

	"github.com/julienschmidt/httprouter"
)

type BrandHandler struct {
	service   domain.BrandService
	validator *validator.Validate
}

func NewBrandHandler(service domain.BrandService) *BrandHandler {
	return &BrandHandler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *BrandHandler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var brand domain.Brand
	if err := json.NewDecoder(r.Body).Decode(&brand); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validasi brand
	if err := h.validator.Struct(brand); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			writeError(w, http.StatusInternalServerError, "Invalid validation error")
			return
		}

		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				errorMessages = append(errorMessages, "Nama brand tidak boleh kosong")
			case "min":
				errorMessages = append(errorMessages, "Nama brand minimal harus 3 karakter")
            case"max":
                errorMessages = append(errorMessages, "Nama brand maksimal harus 200 karakter")
			}
		}

		writeError(w, http.StatusBadRequest, errorMessages[0])
		return
	}

	if err := h.service.Create(&brand); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := Response{
		Status:  http.StatusCreated,
		Message: "Brand created successfully",
		Data:    brand,
	}
	writeJSON(w, http.StatusCreated, resp)
}

func (h *BrandHandler) GetByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	brand, err := h.service.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if brand == nil {
		writeError(w, http.StatusNotFound, "Brand not found")
		return
	}

	resp := Response{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    brand,
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *BrandHandler) GetAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	brands, err := h.service.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := Response{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    brands,
	}
	writeJSON(w, http.StatusOK, resp)
} 