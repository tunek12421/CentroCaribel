package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type EstadoPaquete string

const (
	PaqueteActivo     EstadoPaquete = "ACTIVO"
	PaqueteCompletado EstadoPaquete = "COMPLETADO"
	PaqueteCancelado  EstadoPaquete = "CANCELADO"
)

func (e EstadoPaquete) IsValid() bool {
	switch e {
	case PaqueteActivo, PaqueteCompletado, PaqueteCancelado:
		return true
	}
	return false
}

type PaqueteTratamiento struct {
	ID                  uuid.UUID     `json:"id"`
	PacienteID          uuid.UUID     `json:"paciente_id"`
	TipoTratamiento     string        `json:"tipo_tratamiento"`
	TotalSesiones       int           `json:"total_sesiones"`
	SesionesCompletadas int           `json:"sesiones_completadas"`
	Estado              EstadoPaquete `json:"estado"`
	Notas               string        `json:"notas,omitempty"`
	CreatedBy           uuid.UUID     `json:"created_by"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
}

type PaqueteRepository interface {
	Create(ctx context.Context, p *PaqueteTratamiento) error
	GetByID(ctx context.Context, id uuid.UUID) (*PaqueteTratamiento, error)
	GetByPacienteID(ctx context.Context, pacienteID uuid.UUID) ([]PaqueteTratamiento, error)
	GetActivosByPaciente(ctx context.Context, pacienteID uuid.UUID) ([]PaqueteTratamiento, error)
	IncrementSesiones(ctx context.Context, id uuid.UUID) error
	UpdateEstado(ctx context.Context, id uuid.UUID, estado EstadoPaquete) error
}
