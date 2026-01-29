package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/domain"
)

type RolRepository struct {
	db *sql.DB
}

func NewRolRepository(db *sql.DB) *RolRepository {
	return &RolRepository{db: db}
}

func (r *RolRepository) GetAll(ctx context.Context) ([]domain.Rol, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, nombre, descripcion, permisos, activo, created_at, updated_at FROM roles WHERE activo = true ORDER BY nombre")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []domain.Rol
	for rows.Next() {
		var rol domain.Rol
		if err := rows.Scan(&rol.ID, &rol.Nombre, &rol.Descripcion, &rol.Permisos, &rol.Activo, &rol.CreatedAt, &rol.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, rol)
	}
	return roles, nil
}

func (r *RolRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Rol, error) {
	var rol domain.Rol
	err := r.db.QueryRowContext(ctx,
		"SELECT id, nombre, descripcion, permisos, activo, created_at, updated_at FROM roles WHERE id = $1", id).
		Scan(&rol.ID, &rol.Nombre, &rol.Descripcion, &rol.Permisos, &rol.Activo, &rol.CreatedAt, &rol.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &rol, nil
}

func (r *RolRepository) GetByNombre(ctx context.Context, nombre string) (*domain.Rol, error) {
	var rol domain.Rol
	err := r.db.QueryRowContext(ctx,
		"SELECT id, nombre, descripcion, permisos, activo, created_at, updated_at FROM roles WHERE nombre = $1", nombre).
		Scan(&rol.ID, &rol.Nombre, &rol.Descripcion, &rol.Permisos, &rol.Activo, &rol.CreatedAt, &rol.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &rol, nil
}
