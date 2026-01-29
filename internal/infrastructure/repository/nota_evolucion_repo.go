package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/domain"
)

type NotaEvolucionRepository struct {
	db *sql.DB
}

func NewNotaEvolucionRepository(db *sql.DB) *NotaEvolucionRepository {
	return &NotaEvolucionRepository{db: db}
}

func (r *NotaEvolucionRepository) Create(ctx context.Context, n *domain.NotaEvolucion) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO notas_evolucion (id, historia_id, tipo, contenido, created_by)
		 VALUES ($1, $2, $3, $4, $5)`,
		n.ID, n.HistoriaID, n.Tipo, n.Contenido, n.CreatedBy)
	return err
}

func (r *NotaEvolucionRepository) GetByHistoriaID(ctx context.Context, historiaID uuid.UUID) ([]domain.NotaEvolucion, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, historia_id, tipo, contenido, created_by, created_at
		 FROM notas_evolucion WHERE historia_id = $1
		 ORDER BY created_at DESC`, historiaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notas []domain.NotaEvolucion
	for rows.Next() {
		var n domain.NotaEvolucion
		if err := rows.Scan(&n.ID, &n.HistoriaID, &n.Tipo, &n.Contenido, &n.CreatedBy, &n.CreatedAt); err != nil {
			return nil, err
		}
		notas = append(notas, n)
	}
	return notas, nil
}
