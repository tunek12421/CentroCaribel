package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type HistoriaClinica struct {
	ID                      uuid.UUID `json:"id"`
	PacienteID              uuid.UUID `json:"paciente_id"`
	NumeroHistoria          string    `json:"numero_historia"`
	Estado                  string    `json:"estado"`
	AntecedentesPersonales  string    `json:"antecedentes_personales"`
	AntecedentesFamiliares  string    `json:"antecedentes_familiares"`
	Alergias                string    `json:"alergias"`
	MedicamentosActuales    string    `json:"medicamentos_actuales"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type NotaEvolucion struct {
	ID        uuid.UUID `json:"id"`
	HistoriaID uuid.UUID `json:"historia_id"`
	Tipo      string    `json:"tipo"` // TRATAMIENTO, EVOLUCION, NOTA
	Contenido string    `json:"contenido"`
	CreatedBy uuid.UUID `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type HistoriaClinicaRepository interface {
	Create(ctx context.Context, h *HistoriaClinica) error
	GetByPacienteID(ctx context.Context, pacienteID uuid.UUID) (*HistoriaClinica, error)
	UpdateAntecedentes(ctx context.Context, h *HistoriaClinica) error
	NextNumero(ctx context.Context) (string, error)
}

type NotaEvolucionRepository interface {
	Create(ctx context.Context, n *NotaEvolucion) error
	GetByHistoriaID(ctx context.Context, historiaID uuid.UUID) ([]NotaEvolucion, error)
}
