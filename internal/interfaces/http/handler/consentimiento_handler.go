package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/application/consentimiento"
	"github.com/tunek/centro-caribel/internal/interfaces/http/dto"
	"github.com/tunek/centro-caribel/internal/interfaces/http/middleware"
	apperrors "github.com/tunek/centro-caribel/pkg/errors"
	"github.com/tunek/centro-caribel/pkg/response"
	"github.com/tunek/centro-caribel/pkg/validator"
)

type ConsentimientoHandler struct {
	service *consentimiento.Service
}

func NewConsentimientoHandler(service *consentimiento.Service) *ConsentimientoHandler {
	return &ConsentimientoHandler{service: service}
}

func (h *ConsentimientoHandler) Create(w http.ResponseWriter, r *http.Request) {
	pacienteID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.Error(w, apperrors.NewBadRequest("ID de paciente inválido"))
		return
	}

	var req dto.CreateConsentimientoRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.Error(w, err)
		return
	}

	userID, err := uuid.Parse(middleware.GetUserID(r.Context()))
	if err != nil {
		response.Error(w, apperrors.NewUnauthorized("Usuario no identificado"))
		return
	}

	cons, err := h.service.Create(r.Context(), pacienteID, req.FirmaDigital, req.Contenido, req.AutorizaFotos, userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, cons)
}

func (h *ConsentimientoHandler) GetByPaciente(w http.ResponseWriter, r *http.Request) {
	pacienteID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.Error(w, apperrors.NewBadRequest("ID de paciente inválido"))
		return
	}

	list, err := h.service.GetByPacienteID(r.Context(), pacienteID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, list)
}
