package handler

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/application/cita"
	"github.com/tunek/centro-caribel/internal/interfaces/http/dto"
	"github.com/tunek/centro-caribel/internal/interfaces/http/middleware"
	apperrors "github.com/tunek/centro-caribel/pkg/errors"
	"github.com/tunek/centro-caribel/pkg/response"
	"github.com/tunek/centro-caribel/pkg/validator"
)

type CitaHandler struct {
	service *cita.Service
}

func NewCitaHandler(service *cita.Service) *CitaHandler {
	return &CitaHandler{service: service}
}

func (h *CitaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCitaRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.Error(w, err)
		return
	}

	userID, err := uuid.Parse(middleware.GetUserID(r.Context()))
	if err != nil {
		response.Error(w, apperrors.NewUnauthorized("Usuario no identificado"))
		return
	}

	c, err := h.service.Create(r.Context(), req.PacienteID, req.Fecha, req.Hora, req.TipoTratamiento, req.Turno, req.Observaciones, userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, c)
}

func (h *CitaHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	citas, total, err := h.service.GetAll(r.Context(), page, perPage)
	if err != nil {
		response.Error(w, err)
		return
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	response.JSONWithMeta(w, http.StatusOK, citas, &response.Meta{
		Page:      page,
		PerPage:   perPage,
		Total:     total,
		TotalPage: totalPages,
	})
}

func (h *CitaHandler) UpdateEstado(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.Error(w, apperrors.NewBadRequest("ID inválido"))
		return
	}

	var req dto.UpdateEstadoCitaRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.Error(w, err)
		return
	}

	if err := h.service.UpdateEstado(r.Context(), id, req.Estado); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Estado actualizado"})
}
