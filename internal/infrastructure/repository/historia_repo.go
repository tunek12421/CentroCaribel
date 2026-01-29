package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/domain"
)

type HistoriaClinicaRepository struct {
	db *sql.DB
}

func NewHistoriaClinicaRepository(db *sql.DB) *HistoriaClinicaRepository {
	return &HistoriaClinicaRepository{db: db}
}

func (r *HistoriaClinicaRepository) Create(ctx context.Context, h *domain.HistoriaClinica) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO historias_clinicas (id, paciente_id, numero_historia, estado)
		 VALUES ($1, $2, $3, $4)`,
		h.ID, h.PacienteID, h.NumeroHistoria, h.Estado)
	return err
}

func (r *HistoriaClinicaRepository) GetByPacienteID(ctx context.Context, pacienteID uuid.UUID) (*domain.HistoriaClinica, error) {
	var h domain.HistoriaClinica
	err := r.db.QueryRowContext(ctx,
		`SELECT id, paciente_id, numero_historia, estado, created_at, updated_at
		 FROM historias_clinicas WHERE paciente_id = $1`, pacienteID).
		Scan(&h.ID, &h.PacienteID, &h.NumeroHistoria, &h.Estado, &h.CreatedAt, &h.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (r *HistoriaClinicaRepository) NextNumero(ctx context.Context) (string, error) {
	var seq int
	err := r.db.QueryRowContext(ctx, "SELECT nextval('historias_numero_seq')").Scan(&seq)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("HC-%05d", seq), nil
}
