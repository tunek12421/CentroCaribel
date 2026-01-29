package dto

import (
	"github.com/google/uuid"
	apperrors "github.com/tunek/centro-caribel/pkg/errors"
	"github.com/tunek/centro-caribel/pkg/validator"
)

type CreatePaqueteRequest struct {
	PacienteID      uuid.UUID `json:"paciente_id"`
	TipoTratamiento string    `json:"tipo_tratamiento"`
	TotalSesiones   int       `json:"total_sesiones"`
	Notas           string    `json:"notas"`
}

func (r *CreatePaqueteRequest) Validate() error {
	if r.PacienteID == uuid.Nil {
		return validator.RequiredString("", "paciente_id")
	}
	if err := validator.RequiredString(r.TipoTratamiento, "tipo_tratamiento"); err != nil {
		return err
	}
	if r.TotalSesiones < 1 {
		return apperrors.NewBadRequest("total_sesiones debe ser al menos 1")
	}
	return nil
}
