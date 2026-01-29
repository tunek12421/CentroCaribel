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
		`INSERT INTO historias_clinicas (id, paciente_id, numero_historia, estado,
		 antecedentes_personales, antecedentes_familiares, alergias, medicamentos_actuales)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		h.ID, h.PacienteID, h.NumeroHistoria, h.Estado,
		h.AntecedentesPersonales, h.AntecedentesFamiliares, h.Alergias, h.MedicamentosActuales)
	return err
}

func (r *HistoriaClinicaRepository) GetByPacienteID(ctx context.Context, pacienteID uuid.UUID) (*domain.HistoriaClinica, error) {
	var h domain.HistoriaClinica
	err := r.db.QueryRowContext(ctx,
		`SELECT id, paciente_id, numero_historia, estado,
		 antecedentes_personales, antecedentes_familiares, alergias, medicamentos_actuales,
		 created_at, updated_at
		 FROM historias_clinicas WHERE paciente_id = $1`, pacienteID).
		Scan(&h.ID, &h.PacienteID, &h.NumeroHistoria, &h.Estado,
			&h.AntecedentesPersonales, &h.AntecedentesFamiliares, &h.Alergias, &h.MedicamentosActuales,
			&h.CreatedAt, &h.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (r *HistoriaClinicaRepository) UpdateAntecedentes(ctx context.Context, h *domain.HistoriaClinica) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE historias_clinicas SET
		 antecedentes_personales = $1, antecedentes_familiares = $2,
		 alergias = $3, medicamentos_actuales = $4
		 WHERE id = $5`,
		h.AntecedentesPersonales, h.AntecedentesFamiliares, h.Alergias, h.MedicamentosActuales, h.ID)
	return err
}

func (r *HistoriaClinicaRepository) NextNumero(ctx context.Context) (string, error) {
	var seq int
	err := r.db.QueryRowContext(ctx, "SELECT nextval('historias_numero_seq')").Scan(&seq)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("HC-%05d", seq), nil
}
