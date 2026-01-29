package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type EstadoCita string

const (
	EstadoNueva      EstadoCita = "NUEVA"
	EstadoAgendada   EstadoCita = "AGENDADA"
	EstadoConfirmada EstadoCita = "CONFIRMADA"
	EstadoAtendida   EstadoCita = "ATENDIDA"
	EstadoNoAsistio  EstadoCita = "NO_ASISTIO"
	EstadoCancelada  EstadoCita = "CANCELADA"
	EstadoReagendada EstadoCita = "REAGENDADA"
)

func (e EstadoCita) IsValid() bool {
	switch e {
	case EstadoNueva, EstadoAgendada, EstadoConfirmada,
		EstadoAtendida, EstadoNoAsistio, EstadoCancelada, EstadoReagendada:
		return true
	}
	return false
}

// TransicionesValidas define las transiciones de estado permitidas.
var TransicionesValidas = map[EstadoCita][]EstadoCita{
	EstadoNueva:      {EstadoAgendada, EstadoCancelada},
	EstadoAgendada:   {EstadoConfirmada, EstadoCancelada, EstadoReagendada},
	EstadoConfirmada: {EstadoAtendida, EstadoNoAsistio, EstadoCancelada},
	EstadoReagendada: {EstadoAgendada, EstadoCancelada},
}

func PuedeTransicionar(actual, nueva EstadoCita) bool {
	permitidas, ok := TransicionesValidas[actual]
	if !ok {
		return false
	}
	for _, e := range permitidas {
		if e == nueva {
			return true
		}
	}
	return false
}

type TurnoCita string

const (
	TurnoAM TurnoCita = "AM"
	TurnoPM TurnoCita = "PM"
)

func (t TurnoCita) IsValid() bool {
	return t == TurnoAM || t == TurnoPM
}

type Cita struct {
	ID              uuid.UUID  `json:"id"`
	PacienteID      uuid.UUID  `json:"paciente_id"`
	Fecha           time.Time  `json:"fecha"`
	Hora            string     `json:"hora"`
	TipoTratamiento string     `json:"tipo_tratamiento"`
	Estado          EstadoCita `json:"estado"`
	Turno           TurnoCita  `json:"turno"`
	Observaciones   string     `json:"observaciones,omitempty"`
	PaqueteID       *uuid.UUID `json:"paquete_id,omitempty"`
	CreatedBy       uuid.UUID  `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type CitaRepository interface {
	Create(ctx context.Context, c *Cita) error
	GetByID(ctx context.Context, id uuid.UUID) (*Cita, error)
	GetAll(ctx context.Context, offset, limit int) ([]Cita, int64, error)
	GetByPacienteID(ctx context.Context, pacienteID uuid.UUID) ([]Cita, error)
	GetByFecha(ctx context.Context, fecha time.Time) ([]Cita, error)
	UpdateEstado(ctx context.Context, id uuid.UUID, estado EstadoCita) error
	Reagendar(ctx context.Context, id uuid.UUID, fecha time.Time, hora string, turno TurnoCita) error
	GetAllFiltered(ctx context.Context, offset, limit int, fecha *time.Time, turno *TurnoCita, estado *EstadoCita) ([]Cita, int64, error)
	ExistsByFechaHora(ctx context.Context, fecha time.Time, hora string, excludeID *uuid.UUID) (bool, error)
}
