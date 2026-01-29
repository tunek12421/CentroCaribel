package dto

import (
	"github.com/tunek/centro-caribel/pkg/validator"
)

type CreatePacienteRequest struct {
	NombreCompleto  string `json:"nombre_completo"`
	CI              string `json:"ci"`
	FechaNacimiento string `json:"fecha_nacimiento"` // formato: 2006-01-02
	Celular         string `json:"celular"`
	Direccion       string `json:"direccion"`
}

func (r *CreatePacienteRequest) Validate() error {
	if err := validator.RequiredString(r.NombreCompleto, "nombre_completo"); err != nil {
		return err
	}
	if err := validator.RequiredString(r.CI, "ci"); err != nil {
		return err
	}
	if err := validator.RequiredString(r.FechaNacimiento, "fecha_nacimiento"); err != nil {
		return err
	}
	if err := validator.RequiredString(r.Celular, "celular"); err != nil {
		return err
	}
	return nil
}
