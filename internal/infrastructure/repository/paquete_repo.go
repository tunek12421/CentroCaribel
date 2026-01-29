package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/domain"
)

type PaqueteRepository struct {
	db *sql.DB
}

func NewPaqueteRepository(db *sql.DB) *PaqueteRepository {
	return &PaqueteRepository{db: db}
}

func (r *PaqueteRepository) Create(ctx context.Context, p *domain.PaqueteTratamiento) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO paquetes_tratamiento (id, paciente_id, tipo_tratamiento, total_sesiones, sesiones_completadas, estado, notas, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		p.ID, p.PacienteID, p.TipoTratamiento, p.TotalSesiones, p.SesionesCompletadas, p.Estado, p.Notas, p.CreatedBy)
	return err
}

func (r *PaqueteRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.PaqueteTratamiento, error) {
	var p domain.PaqueteTratamiento
	err := r.db.QueryRowContext(ctx,
		`SELECT id, paciente_id, tipo_tratamiento, total_sesiones, sesiones_completadas, estado, notas, created_by, created_at, updated_at
		 FROM paquetes_tratamiento WHERE id = $1`, id).
		Scan(&p.ID, &p.PacienteID, &p.TipoTratamiento, &p.TotalSesiones, &p.SesionesCompletadas, &p.Estado, &p.Notas, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PaqueteRepository) GetByPacienteID(ctx context.Context, pacienteID uuid.UUID) ([]domain.PaqueteTratamiento, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, paciente_id, tipo_tratamiento, total_sesiones, sesiones_completadas, estado, notas, created_by, created_at, updated_at
		 FROM paquetes_tratamiento WHERE paciente_id = $1 ORDER BY created_at DESC`, pacienteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paquetes []domain.PaqueteTratamiento
	for rows.Next() {
		var p domain.PaqueteTratamiento
		if err := rows.Scan(&p.ID, &p.PacienteID, &p.TipoTratamiento, &p.TotalSesiones, &p.SesionesCompletadas, &p.Estado, &p.Notas, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		paquetes = append(paquetes, p)
	}
	return paquetes, nil
}

func (r *PaqueteRepository) GetActivosByPaciente(ctx context.Context, pacienteID uuid.UUID) ([]domain.PaqueteTratamiento, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, paciente_id, tipo_tratamiento, total_sesiones, sesiones_completadas, estado, notas, created_by, created_at, updated_at
		 FROM paquetes_tratamiento WHERE paciente_id = $1 AND estado = 'ACTIVO' ORDER BY created_at DESC`, pacienteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paquetes []domain.PaqueteTratamiento
	for rows.Next() {
		var p domain.PaqueteTratamiento
		if err := rows.Scan(&p.ID, &p.PacienteID, &p.TipoTratamiento, &p.TotalSesiones, &p.SesionesCompletadas, &p.Estado, &p.Notas, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		paquetes = append(paquetes, p)
	}
	return paquetes, nil
}

func (r *PaqueteRepository) IncrementSesiones(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE paquetes_tratamiento SET sesiones_completadas = sesiones_completadas + 1 WHERE id = $1`, id)
	return err
}

func (r *PaqueteRepository) UpdateEstado(ctx context.Context, id uuid.UUID, estado domain.EstadoPaquete) error {
	_, err := r.db.ExecContext(ctx, "UPDATE paquetes_tratamiento SET estado = $1 WHERE id = $2", estado, id)
	return err
}
