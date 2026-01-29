package handler

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/application/paciente"
	"github.com/tunek/centro-caribel/internal/interfaces/http/dto"
	"github.com/tunek/centro-caribel/internal/interfaces/http/middleware"
	apperrors "github.com/tunek/centro-caribel/pkg/errors"
	"github.com/tunek/centro-caribel/pkg/response"
	"github.com/tunek/centro-caribel/pkg/validator"
)

type PacienteHandler struct {
	service *paciente.Service
}

func NewPacienteHandler(service *paciente.Service) *PacienteHandler {
	return &PacienteHandler{service: service}
}

func (h *PacienteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePacienteRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.Error(w, err)
		return
	}

	userID, err := uuid.Parse(middleware.GetUserID(r.Context()))
	if err != nil {
		response.Error(w, apperrors.NewUnauthorized("Usuario no identificado"))
		return
	}

	pac, err := h.service.Create(r.Context(), req.NombreCompleto, req.CI, req.FechaNacimiento, req.Celular, req.Direccion, userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, pac)
}

func (h *PacienteHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	query := r.URL.Query().Get("q")
	pacientes, total, err := h.service.GetAll(r.Context(), page, perPage, query)
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

	response.JSONWithMeta(w, http.StatusOK, pacientes, &response.Meta{
		Page:      page,
		PerPage:   perPage,
		Total:     total,
		TotalPage: totalPages,
	})
}

func (h *PacienteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.Error(w, apperrors.NewBadRequest("ID inválido"))
		return
	}

	pac, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, pac)
}
