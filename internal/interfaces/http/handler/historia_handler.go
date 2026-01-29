package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/application/historia"
	apperrors "github.com/tunek/centro-caribel/pkg/errors"
	"github.com/tunek/centro-caribel/pkg/response"
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
