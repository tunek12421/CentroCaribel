package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Paciente struct {
	ID              uuid.UUID `json:"id"`
	Codigo          string    `json:"codigo"`
	NombreCompleto  string    `json:"nombre_completo"`
	CI              string    `json:"ci"`
	FechaNacimiento time.Time `json:"fecha_nacimiento"`
	Celular         string    `json:"celular"`
	Direccion       string    `json:"direccion,omitempty"`
	CreatedBy       uuid.UUID `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type PacienteRepository interface {
	Create(ctx context.Context, p *Paciente) error
	GetByID(ctx context.Context, id uuid.UUID) (*Paciente, error)
	GetByCI(ctx context.Context, ci string) (*Paciente, error)
	GetAll(ctx context.Context, offset, limit int) ([]Paciente, int64, error)
	Search(ctx context.Context, query string, offset, limit int) ([]Paciente, int64, error)
	Update(ctx context.Context, p *Paciente) error
	NextCodigo(ctx context.Context) (string, error)
}
