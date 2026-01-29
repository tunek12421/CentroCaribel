package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Usuario struct {
	ID             uuid.UUID `json:"id"`
	NombreCompleto string    `json:"nombre_completo"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	RolID          uuid.UUID `json:"rol_id"`
	Rol            *Rol      `json:"rol,omitempty"`
	Activo         bool      `json:"activo"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UsuarioRepository interface {
	Create(ctx context.Context, u *Usuario) error
	GetByID(ctx context.Context, id uuid.UUID) (*Usuario, error)
	GetByEmail(ctx context.Context, email string) (*Usuario, error)
	GetAll(ctx context.Context, offset, limit int) ([]Usuario, int64, error)
	Update(ctx context.Context, u *Usuario) error
	Delete(ctx context.Context, id uuid.UUID) error
}
