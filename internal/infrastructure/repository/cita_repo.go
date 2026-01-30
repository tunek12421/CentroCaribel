package repository

import (
	"context"
	"database/sql"
	"fmt"
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

const citaColumns = `c.id, c.paciente_id, p.nombre_completo, c.fecha, TO_CHAR(c.hora, 'HH24:MI') as hora, c.tipo_tratamiento, c.estado, c.turno, c.observaciones, c.paquete_id, c.created_by, c.created_at, c.updated_at`
const citaFrom = `citas c JOIN pacientes p ON c.paciente_id = p.id`

func scanCita(row interface{ Scan(dest ...any) error }) (domain.Cita, error) {
	var c domain.Cita
	err := row.Scan(&c.ID, &c.PacienteID, &c.PacienteNombre, &c.Fecha, &c.Hora, &c.TipoTratamiento, &c.Estado, &c.Turno, &c.Observaciones, &c.PaqueteID, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *CitaRepository) Create(ctx context.Context, c *domain.Cita) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO citas (id, paciente_id, fecha, hora, tipo_tratamiento, estado, turno, observaciones, paquete_id, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		c.ID, c.PacienteID, c.Fecha, c.Hora, c.TipoTratamiento, c.Estado, c.Turno, c.Observaciones, c.PaqueteID, c.CreatedBy)
	return err
}

func (r *CitaRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Cita, error) {
	c, err := scanCita(r.db.QueryRowContext(ctx,
		`SELECT `+citaColumns+` FROM `+citaFrom+` WHERE c.id = $1`, id))
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CitaRepository) GetAll(ctx context.Context, offset, limit int) ([]domain.Cita, int64, error) {
	return r.GetAllFiltered(ctx, offset, limit, nil, nil, nil)
}

func (r *CitaRepository) GetAllFiltered(ctx context.Context, offset, limit int, fecha *time.Time, turno *domain.TurnoCita, estado *domain.EstadoCita) ([]domain.Cita, int64, error) {
	where := "WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	if fecha != nil {
		where += fmt.Sprintf(" AND c.fecha = $%d", argIdx)
		args = append(args, *fecha)
		argIdx++
	}
	if turno != nil {
		where += fmt.Sprintf(" AND c.turno = $%d", argIdx)
		args = append(args, *turno)
		argIdx++
	}
	if estado != nil {
		where += fmt.Sprintf(" AND c.estado = $%d", argIdx)
		args = append(args, *estado)
		argIdx++
	}

	var total int64
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+citaFrom+" "+where, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(
		`SELECT %s FROM %s %s ORDER BY c.fecha DESC, c.hora DESC LIMIT $%d OFFSET $%d`,
		citaColumns, citaFrom, where, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var citas []domain.Cita
	for rows.Next() {
		c, err := scanCita(rows)
		if err != nil {
			return nil, 0, err
		}
		citas = append(citas, c)
	}
	return citas, total, nil
}

func (r *CitaRepository) GetByPacienteID(ctx context.Context, pacienteID uuid.UUID) ([]domain.Cita, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+citaColumns+` FROM `+citaFrom+` WHERE c.paciente_id = $1 ORDER BY c.fecha DESC`, pacienteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var citas []domain.Cita
	for rows.Next() {
		c, err := scanCita(rows)
		if err != nil {
			return nil, err
		}
		citas = append(citas, c)
	}
	return citas, nil
}

func (r *CitaRepository) GetByFecha(ctx context.Context, fecha time.Time) ([]domain.Cita, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+citaColumns+` FROM `+citaFrom+` WHERE c.fecha = $1 ORDER BY c.hora`, fecha)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var citas []domain.Cita
	for rows.Next() {
		c, err := scanCita(rows)
		if err != nil {
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

func (r *CitaRepository) ExistsByFechaHora(ctx context.Context, fecha time.Time, hora string, excludeID *uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM citas WHERE fecha = $1 AND hora = $2 AND estado NOT IN ('CANCELADA', 'REAGENDADA')`
	args := []interface{}{fecha, hora}
	if excludeID != nil {
		query += ` AND id != $3`
		args = append(args, *excludeID)
	}
	query += `)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	return exists, err
}
