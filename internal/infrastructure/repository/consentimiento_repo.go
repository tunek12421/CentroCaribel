package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/domain"
)

type ConsentimientoRepository struct {
	db *sql.DB
}

func NewConsentimientoRepository(db *sql.DB) *ConsentimientoRepository {
	return &ConsentimientoRepository{db: db}
}

func (r *ConsentimientoRepository) Create(ctx context.Context, c *domain.Consentimiento) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO consentimientos (id, paciente_id, fecha_firma, firma_digital, autoriza_fotos, contenido, registrado_por)
		 VALUES ($1, $2, NOW(), $3, $4, $5, $6)`,
		c.ID, c.PacienteID, c.FirmaDigital, c.AutorizaFotos, c.Contenido, c.RegistradoPor)
	return err
}

func (r *ConsentimientoRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Consentimiento, error) {
	var c domain.Consentimiento
	err := r.db.QueryRowContext(ctx,
		`SELECT id, paciente_id, fecha_firma, firma_digital, autoriza_fotos, contenido, registrado_por, created_at
		 FROM consentimientos WHERE id = $1`, id).
		Scan(&c.ID, &c.PacienteID, &c.FechaFirma, &c.FirmaDigital, &c.AutorizaFotos, &c.Contenido, &c.RegistradoPor, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ConsentimientoRepository) GetByPacienteID(ctx context.Context, pacienteID uuid.UUID) ([]domain.Consentimiento, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, paciente_id, fecha_firma, autoriza_fotos, contenido, registrado_por, created_at
		 FROM consentimientos WHERE paciente_id = $1 ORDER BY created_at DESC`, pacienteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Consentimiento
	for rows.Next() {
		var c domain.Consentimiento
		if err := rows.Scan(&c.ID, &c.PacienteID, &c.FechaFirma, &c.AutorizaFotos, &c.Contenido, &c.RegistradoPor, &c.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}
