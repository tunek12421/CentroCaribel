package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/domain"
)

type CitaRepository struct {
	db *sql.DB
}

func NewCitaRepository(db *sql.DB) *CitaRepository {
	return &CitaRepository{db: db}
}

func (r *CitaRepository) Create(ctx context.Context, c *domain.Cita) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO citas (id, paciente_id, fecha, hora, tipo_tratamiento, estado, turno, observaciones, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		c.ID, c.PacienteID, c.Fecha, c.Hora, c.TipoTratamiento, c.Estado, c.Turno, c.Observaciones, c.CreatedBy)
	return err
}

func (r *CitaRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Cita, error) {
	var c domain.Cita
	err := r.db.QueryRowContext(ctx,
		`SELECT id, paciente_id, fecha, TO_CHAR(hora, 'HH24:MI') as hora, tipo_tratamiento, estado, turno, observaciones, created_by, created_at, updated_at
		 FROM citas WHERE id = $1`, id).
		Scan(&c.ID, &c.PacienteID, &c.Fecha, &c.Hora, &c.TipoTratamiento, &c.Estado, &c.Turno, &c.Observaciones, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CitaRepository) GetAll(ctx context.Context, offset, limit int) ([]domain.Cita, int64, error) {
	var total int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM citas").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, paciente_id, fecha, TO_CHAR(hora, 'HH24:MI') as hora, tipo_tratamiento, estado, turno, observaciones, created_by, created_at, updated_at
		 FROM citas ORDER BY fecha DESC, hora DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var citas []domain.Cita
	for rows.Next() {
		var c domain.Cita
		if err := rows.Scan(&c.ID, &c.PacienteID, &c.Fecha, &c.Hora, &c.TipoTratamiento, &c.Estado, &c.Turno, &c.Observaciones, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, 0, err
		}
		citas = append(citas, c)
	}
	return citas, total, nil
}

func (r *CitaRepository) GetByPacienteID(ctx context.Context, pacienteID uuid.UUID) ([]domain.Cita, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, paciente_id, fecha, TO_CHAR(hora, 'HH24:MI') as hora, tipo_tratamiento, estado, turno, observaciones, created_by, created_at, updated_at
		 FROM citas WHERE paciente_id = $1 ORDER BY fecha DESC`, pacienteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var citas []domain.Cita
	for rows.Next() {
		var c domain.Cita
		if err := rows.Scan(&c.ID, &c.PacienteID, &c.Fecha, &c.Hora, &c.TipoTratamiento, &c.Estado, &c.Turno, &c.Observaciones, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		citas = append(citas, c)
	}
	return citas, nil
}

func (r *CitaRepository) GetByFecha(ctx context.Context, fecha time.Time) ([]domain.Cita, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, paciente_id, fecha, TO_CHAR(hora, 'HH24:MI') as hora, tipo_tratamiento, estado, turno, observaciones, created_by, created_at, updated_at
		 FROM citas WHERE fecha = $1 ORDER BY hora`, fecha)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var citas []domain.Cita
	for rows.Next() {
		var c domain.Cita
		if err := rows.Scan(&c.ID, &c.PacienteID, &c.Fecha, &c.Hora, &c.TipoTratamiento, &c.Estado, &c.Turno, &c.Observaciones, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		citas = append(citas, c)
	}
	return citas, nil
}

func (r *CitaRepository) UpdateEstado(ctx context.Context, id uuid.UUID, estado domain.EstadoCita) error {
	_, err := r.db.ExecContext(ctx, "UPDATE citas SET estado = $1 WHERE id = $2", estado, id)
	return err
}

func (r *CitaRepository) Reagendar(ctx context.Context, id uuid.UUID, fecha time.Time, hora string, turno domain.TurnoCita) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE citas SET estado = $1, fecha = $2, hora = $3, turno = $4 WHERE id = $5",
		domain.EstadoReagendada, fecha, hora, turno, id)
	return err
}
