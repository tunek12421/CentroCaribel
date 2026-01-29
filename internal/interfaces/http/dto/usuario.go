package dto

import (
	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/pkg/validator"
)

type CreateUsuarioRequest struct {
	NombreCompleto string    `json:"nombre_completo"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	RolID          uuid.UUID `json:"rol_id"`
}

func (r *CreateUsuarioRequest) Validate() error {
	if err := validator.RequiredString(r.NombreCompleto, "nombre_completo"); err != nil {
		return err
	}
	if err := validator.RequiredString(r.Email, "email"); err != nil {
		return err
	}
	if err := validator.ValidEmail(r.Email); err != nil {
		return err
	}
	if err := validator.MinLength(r.Password, "password", 8); err != nil {
		return err
	}
	if r.RolID == uuid.Nil {
		return validator.RequiredString("", "rol_id")
	}
	return nil
}

type UpdateUsuarioRequest struct {
	NombreCompleto string    `json:"nombre_completo,omitempty"`
	Email          string    `json:"email,omitempty"`
	RolID          uuid.UUID `json:"rol_id,omitempty"`
	Activo         *bool     `json:"activo,omitempty"`
}
