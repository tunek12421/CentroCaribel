package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Consentimiento struct {
	ID             uuid.UUID `json:"id"`
	PacienteID     uuid.UUID `json:"paciente_id"`
	FechaFirma     time.Time `json:"fecha_firma"`
	FirmaDigital   []byte    `json:"firma_digital,omitempty"`
	AutorizaFotos  bool      `json:"autoriza_fotos"`
	Contenido      string    `json:"contenido"`
	RegistradoPor  uuid.UUID `json:"registrado_por"`
	CreatedAt      time.Time `json:"created_at"`
}

type ConsentimientoRepository interface {
	Create(ctx context.Context, c *Consentimiento) error
	GetByID(ctx context.Context, id uuid.UUID) (*Consentimiento, error)
	GetByPacienteID(ctx context.Context, pacienteID uuid.UUID) ([]Consentimiento, error)
}
