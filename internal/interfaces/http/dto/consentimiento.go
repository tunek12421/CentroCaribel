package dto

import (
	"github.com/tunek/centro-caribel/pkg/validator"
)

type CreateConsentimientoRequest struct {
	FirmaDigital  string `json:"firma_digital"` // base64 encoded
	AutorizaFotos bool   `json:"autoriza_fotos"`
	Contenido     string `json:"contenido"`
}

func (r *CreateConsentimientoRequest) Validate() error {
	if err := validator.RequiredString(r.Contenido, "contenido"); err != nil {
		return err
	}
	return nil
}
