package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/domain"
)

type PacienteRepository struct {
	db *sql.DB
}

func NewPacienteRepository(db *sql.DB) *PacienteRepository {
	return &PacienteRepository{db: db}
}

func (r *PacienteRepository) Create(ctx context.Context, p *domain.Paciente) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO pacientes (id, codigo, nombre_completo, ci, fecha_nacimiento, celular, direccion, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		p.ID, p.Codigo, p.NombreCompleto, p.CI, p.FechaNacimiento, p.Celular, p.Direccion, p.CreatedBy)
	return err
}

func (r *PacienteRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Paciente, error) {
	var p domain.Paciente
	err := r.db.QueryRowContext(ctx,
		`SELECT id, codigo, nombre_completo, ci, fecha_nacimiento, celular, direccion, created_by, created_at, updated_at
		 FROM pacientes WHERE id = $1`, id).
		Scan(&p.ID, &p.Codigo, &p.NombreCompleto, &p.CI, &p.FechaNacimiento, &p.Celular, &p.Direccion, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PacienteRepository) GetByCI(ctx context.Context, ci string) (*domain.Paciente, error) {
	var p domain.Paciente
	err := r.db.QueryRowContext(ctx,
		`SELECT id, codigo, nombre_completo, ci, fecha_nacimiento, celular, direccion, created_by, created_at, updated_at
		 FROM pacientes WHERE ci = $1`, ci).
		Scan(&p.ID, &p.Codigo, &p.NombreCompleto, &p.CI, &p.FechaNacimiento, &p.Celular, &p.Direccion, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PacienteRepository) GetAll(ctx context.Context, offset, limit int) ([]domain.Paciente, int64, error) {
	var total int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM pacientes").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, codigo, nombre_completo, ci, fecha_nacimiento, celular, direccion, created_by, created_at, updated_at
		 FROM pacientes ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var pacientes []domain.Paciente
	for rows.Next() {
		var p domain.Paciente
		if err := rows.Scan(&p.ID, &p.Codigo, &p.NombreCompleto, &p.CI, &p.FechaNacimiento, &p.Celular, &p.Direccion, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		pacientes = append(pacientes, p)
	}
	return pacientes, total, nil
}

func (r *PacienteRepository) Update(ctx context.Context, p *domain.Paciente) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE pacientes SET nombre_completo = $1, celular = $2, direccion = $3
		 WHERE id = $4`,
		p.NombreCompleto, p.Celular, p.Direccion, p.ID)
	return err
}

func (r *PacienteRepository) Search(ctx context.Context, query string, offset, limit int) ([]domain.Paciente, int64, error) {
	like := "%" + query + "%"

	var total int64
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM pacientes WHERE unaccent(lower(nombre_completo)) LIKE unaccent(lower($1)) OR ci LIKE $2`, like, like).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, codigo, nombre_completo, ci, fecha_nacimiento, celular, direccion, created_by, created_at, updated_at
		 FROM pacientes WHERE unaccent(lower(nombre_completo)) LIKE unaccent(lower($1)) OR ci LIKE $2
		 ORDER BY nombre_completo ASC LIMIT $3 OFFSET $4`, like, like, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var pacientes []domain.Paciente
	for rows.Next() {
		var p domain.Paciente
		if err := rows.Scan(&p.ID, &p.Codigo, &p.NombreCompleto, &p.CI, &p.FechaNacimiento, &p.Celular, &p.Direccion, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		pacientes = append(pacientes, p)
	}
	return pacientes, total, nil
}

func (r *PacienteRepository) NextCodigo(ctx context.Context) (string, error) {
	var seq int
	err := r.db.QueryRowContext(ctx, "SELECT nextval('pacientes_codigo_seq')").Scan(&seq)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("PAC-%05d", seq), nil
}
