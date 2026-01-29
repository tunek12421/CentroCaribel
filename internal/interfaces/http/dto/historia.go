package dto

import (
	"github.com/tunek/centro-caribel/pkg/validator"
)

type UpdateAntecedentesRequest struct {
	AntecedentesPersonales string `json:"antecedentes_personales"`
	AntecedentesFamiliares string `json:"antecedentes_familiares"`
	Alergias               string `json:"alergias"`
	MedicamentosActuales   string `json:"medicamentos_actuales"`
}

func (r *UpdateAntecedentesRequest) Validate() error {
	return nil
}

type CreateNotaRequest struct {
	Tipo      string `json:"tipo"`
	Contenido string `json:"contenido"`
}

func (r *CreateNotaRequest) Validate() error {
	if err := validator.RequiredString(r.Tipo, "tipo"); err != nil {
		return err
	}
	if err := validator.RequiredString(r.Contenido, "contenido"); err != nil {
		return err
	}
	return nil
}
