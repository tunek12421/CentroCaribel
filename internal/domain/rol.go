package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Rol struct {
	ID          uuid.UUID       `json:"id"`
	Nombre      string          `json:"nombre"`
	Descripcion string          `json:"descripcion"`
	Permisos    json.RawMessage `json:"permisos"`
	Activo      bool            `json:"activo"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type RolRepository interface {
	GetAll(ctx context.Context) ([]Rol, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Rol, error)
	GetByNombre(ctx context.Context, nombre string) (*Rol, error)
}
