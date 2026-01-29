package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/application/historia"
	"github.com/tunek/centro-caribel/internal/interfaces/http/dto"
	"github.com/tunek/centro-caribel/internal/interfaces/http/middleware"
	apperrors "github.com/tunek/centro-caribel/pkg/errors"
	"github.com/tunek/centro-caribel/pkg/response"
	"github.com/tunek/centro-caribel/pkg/validator"
)

type HistoriaHandler struct {
	service *historia.Service
}

func NewHistoriaHandler(service *historia.Service) *HistoriaHandler {
	return &HistoriaHandler{service: service}
}

func (h *HistoriaHandler) GetByPaciente(w http.ResponseWriter, r *http.Request) {
	pacienteID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.Error(w, apperrors.NewBadRequest("ID de paciente inválido"))
		return
	}

	hist, err := h.service.GetByPacienteID(r.Context(), pacienteID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, hist)
}

func (h *HistoriaHandler) UpdateAntecedentes(w http.ResponseWriter, r *http.Request) {
	pacienteID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.Error(w, apperrors.NewBadRequest("ID de paciente inválido"))
		return
	}

	var req dto.UpdateAntecedentesRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.Error(w, err)
		return
	}

	hist, err := h.service.UpdateAntecedentes(r.Context(), pacienteID,
		req.AntecedentesPersonales, req.AntecedentesFamiliares,
		req.Alergias, req.MedicamentosActuales)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, hist)
}

func (h *HistoriaHandler) GetNotas(w http.ResponseWriter, r *http.Request) {
	pacienteID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.Error(w, apperrors.NewBadRequest("ID de paciente inválido"))
		return
	}

	notas, err := h.service.GetNotas(r.Context(), pacienteID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, notas)
}

func (h *HistoriaHandler) CreateNota(w http.ResponseWriter, r *http.Request) {
	pacienteID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.Error(w, apperrors.NewBadRequest("ID de paciente inválido"))
		return
	}

	var req dto.CreateNotaRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.Error(w, err)
		return
	}

	userID, err := uuid.Parse(middleware.GetUserID(r.Context()))
	if err != nil {
		response.Error(w, apperrors.NewUnauthorized("Usuario no identificado"))
		return
	}

	nota, err := h.service.CreateNota(r.Context(), pacienteID, req.Tipo, req.Contenido, userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, nota)
}
