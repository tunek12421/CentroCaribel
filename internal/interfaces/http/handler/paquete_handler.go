package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/application/paquete"
	"github.com/tunek/centro-caribel/internal/interfaces/http/dto"
	"github.com/tunek/centro-caribel/internal/interfaces/http/middleware"
	apperrors "github.com/tunek/centro-caribel/pkg/errors"
	"github.com/tunek/centro-caribel/pkg/response"
	"github.com/tunek/centro-caribel/pkg/validator"
)

type PaqueteHandler struct {
	service *paquete.Service
}

func NewPaqueteHandler(service *paquete.Service) *PaqueteHandler {
	return &PaqueteHandler{service: service}
}

func (h *PaqueteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePaqueteRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.Error(w, err)
		return
	}

	userID, err := uuid.Parse(middleware.GetUserID(r.Context()))
	if err != nil {
		response.Error(w, apperrors.NewUnauthorized("Usuario no identificado"))
		return
	}

	p, err := h.service.Create(r.Context(), req.PacienteID, req.TipoTratamiento, req.TotalSesiones, req.Notas, userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, p)
}

func (h *PaqueteHandler) GetByPaciente(w http.ResponseWriter, r *http.Request) {
	pacienteID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.Error(w, apperrors.NewBadRequest("ID inválido"))
		return
	}

	onlyActivos := r.URL.Query().Get("activos") == "true"

	var paquetes []interface{}
	if onlyActivos {
		result, err := h.service.GetActivosByPaciente(r.Context(), pacienteID)
		if err != nil {
			response.Error(w, err)
			return
		}
		for _, p := range result {
			paquetes = append(paquetes, p)
		}
	} else {
		result, err := h.service.GetByPacienteID(r.Context(), pacienteID)
		if err != nil {
			response.Error(w, err)
			return
		}
		for _, p := range result {
			paquetes = append(paquetes, p)
		}
	}

	if paquetes == nil {
		paquetes = []interface{}{}
	}
	response.JSON(w, http.StatusOK, paquetes)
}
