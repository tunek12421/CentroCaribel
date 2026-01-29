package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/domain"
)

type UsuarioRepository struct {
	db *sql.DB
}

func NewUsuarioRepository(db *sql.DB) *UsuarioRepository {
	return &UsuarioRepository{db: db}
}

func (r *UsuarioRepository) Create(ctx context.Context, u *domain.Usuario) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO usuarios (id, nombre_completo, email, password_hash, rol_id, activo)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		u.ID, u.NombreCompleto, u.Email, u.PasswordHash, u.RolID, u.Activo)
	return err
}

func (r *UsuarioRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Usuario, error) {
	var u domain.Usuario
	var rol domain.Rol
	err := r.db.QueryRowContext(ctx,
		`SELECT u.id, u.nombre_completo, u.email, u.password_hash, u.rol_id, u.activo, u.created_at, u.updated_at,
		        r.id, r.nombre, r.descripcion, r.permisos, r.activo
		 FROM usuarios u JOIN roles r ON u.rol_id = r.id
		 WHERE u.id = $1`, id).
		Scan(&u.ID, &u.NombreCompleto, &u.Email, &u.PasswordHash, &u.RolID, &u.Activo, &u.CreatedAt, &u.UpdatedAt,
			&rol.ID, &rol.Nombre, &rol.Descripcion, &rol.Permisos, &rol.Activo)
	if err != nil {
		return nil, err
	}
	u.Rol = &rol
	return &u, nil
}

func (r *UsuarioRepository) GetByEmail(ctx context.Context, email string) (*domain.Usuario, error) {
	var u domain.Usuario
	err := r.db.QueryRowContext(ctx,
		`SELECT id, nombre_completo, email, password_hash, rol_id, activo, created_at, updated_at
		 FROM usuarios WHERE email = $1`, email).
		Scan(&u.ID, &u.NombreCompleto, &u.Email, &u.PasswordHash, &u.RolID, &u.Activo, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UsuarioRepository) GetAll(ctx context.Context, offset, limit int) ([]domain.Usuario, int64, error) {
	var total int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM usuarios").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT u.id, u.nombre_completo, u.email, u.rol_id, u.activo, u.created_at, u.updated_at,
		        r.nombre
		 FROM usuarios u JOIN roles r ON u.rol_id = r.id
		 ORDER BY u.created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var usuarios []domain.Usuario
	for rows.Next() {
		var u domain.Usuario
		var rolNombre string
		if err := rows.Scan(&u.ID, &u.NombreCompleto, &u.Email, &u.RolID, &u.Activo, &u.CreatedAt, &u.UpdatedAt, &rolNombre); err != nil {
			return nil, 0, err
		}
		u.Rol = &domain.Rol{Nombre: rolNombre}
		usuarios = append(usuarios, u)
	}
	return usuarios, total, nil
}

func (r *UsuarioRepository) Update(ctx context.Context, u *domain.Usuario) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE usuarios SET nombre_completo = $1, email = $2, rol_id = $3, activo = $4
		 WHERE id = $5`,
		u.NombreCompleto, u.Email, u.RolID, u.Activo, u.ID)
	return err
}

func (r *UsuarioRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "UPDATE usuarios SET activo = false WHERE id = $1", id)
	return err
}
