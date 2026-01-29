package dto

import (
	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/domain"
	apperrors "github.com/tunek/centro-caribel/pkg/errors"
	"github.com/tunek/centro-caribel/pkg/validator"
)

type CreateCitaRequest struct {
	PacienteID      uuid.UUID        `json:"paciente_id"`
	Fecha           string           `json:"fecha"` // formato: 2006-01-02
	Hora            string           `json:"hora"`  // formato: 15:04
	TipoTratamiento string           `json:"tipo_tratamiento"`
	Turno           domain.TurnoCita `json:"turno"`
	Observaciones   string           `json:"observaciones"`
	PaqueteID       *uuid.UUID       `json:"paquete_id,omitempty"`
}

func (r *CreateCitaRequest) Validate() error {
	if r.PacienteID == uuid.Nil {
		return validator.RequiredString("", "paciente_id")
	}
	if err := validator.RequiredString(r.Fecha, "fecha"); err != nil {
		return err
	}
	if err := validator.RequiredString(r.Hora, "hora"); err != nil {
		return err
	}
	if err := validator.RequiredString(r.TipoTratamiento, "tipo_tratamiento"); err != nil {
		return err
	}
	if !r.Turno.IsValid() {
		return apperrors.NewBadRequest("El turno debe ser 'AM' o 'PM'")
	}
	return nil
}

type UpdateEstadoCitaRequest struct {
	Estado domain.EstadoCita `json:"estado"`
	Fecha  string            `json:"fecha,omitempty"`
	Hora   string            `json:"hora,omitempty"`
	Turno  domain.TurnoCita  `json:"turno,omitempty"`
}

func (r *UpdateEstadoCitaRequest) Validate() error {
	if !r.Estado.IsValid() {
		return apperrors.NewBadRequest("Estado de cita inválido")
	}
	if r.Estado == domain.EstadoReagendada {
		if r.Fecha == "" {
			return apperrors.NewBadRequest("La fecha es requerida para reagendar")
		}
		if r.Hora == "" {
			return apperrors.NewBadRequest("La hora es requerida para reagendar")
		}
		if r.Turno == "" {
			r.Turno = domain.TurnoAM
		}
		if !r.Turno.IsValid() {
			return apperrors.NewBadRequest("El turno debe ser 'AM' o 'PM'")
		}
	}
	return nil
}
