package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type HistoriaClinica struct {
	ID             uuid.UUID `json:"id"`
	PacienteID     uuid.UUID `json:"paciente_id"`
	NumeroHistoria string    `json:"numero_historia"`
	Estado         string    `json:"estado"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type HistoriaClinicaRepository interface {
	Create(ctx context.Context, h *HistoriaClinica) error
	GetByPacienteID(ctx context.Context, pacienteID uuid.UUID) (*HistoriaClinica, error)
	NextNumero(ctx context.Context) (string, error)
}
