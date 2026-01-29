package cita

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/domain"
	apperrors "github.com/tunek/centro-caribel/pkg/errors"
)

type Service struct {
	repo         domain.CitaRepository
	pacienteRepo domain.PacienteRepository
}

func NewService(repo domain.CitaRepository, pacienteRepo domain.PacienteRepository) *Service {
	return &Service{repo: repo, pacienteRepo: pacienteRepo}
}

func (s *Service) Create(ctx context.Context, pacienteID uuid.UUID, fecha, hora, tipoTratamiento string, turno domain.TurnoCita, observaciones string, createdBy uuid.UUID) (*domain.Cita, error) {
	if _, err := s.pacienteRepo.GetByID(ctx, pacienteID); err != nil {
		return nil, apperrors.NewNotFound("Paciente")
	}

	fechaParsed, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		return nil, apperrors.NewBadRequest("Formato de fecha inválido. Use YYYY-MM-DD")
	}

	if _, err := time.Parse("15:04", hora); err != nil {
		return nil, apperrors.NewBadRequest("Formato de hora inválido. Use HH:MM")
	}

	c := &domain.Cita{
		ID:              uuid.New(),
		PacienteID:      pacienteID,
		Fecha:           fechaParsed,
		Hora:            hora,
		TipoTratamiento: tipoTratamiento,
		Estado:          domain.EstadoNueva,
		Turno:           turno,
		Observaciones:   observaciones,
		CreatedBy:       createdBy,
	}

	if err := s.repo.Create(ctx, c); err != nil {
		return nil, apperrors.NewInternal("Error al crear la cita")
	}

	return c, nil
}

func (s *Service) GetAll(ctx context.Context, page, perPage int) ([]domain.Cita, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage
	return s.repo.GetAll(ctx, offset, perPage)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.Cita, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFound("Cita")
	}
	return c, nil
}

func (s *Service) UpdateEstado(ctx context.Context, id uuid.UUID, nuevoEstado domain.EstadoCita) error {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFound("Cita")
	}

	if !domain.PuedeTransicionar(c.Estado, nuevoEstado) {
		return apperrors.NewBadRequest("Transición de estado no permitida: " + string(c.Estado) + " -> " + string(nuevoEstado))
	}

	return s.repo.UpdateEstado(ctx, id, nuevoEstado)
}

func (s *Service) Reagendar(ctx context.Context, id uuid.UUID, fecha, hora string, turno domain.TurnoCita) error {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFound("Cita")
	}

	if !domain.PuedeTransicionar(c.Estado, domain.EstadoReagendada) {
		return apperrors.NewBadRequest("No se puede reagendar desde el estado: " + string(c.Estado))
	}

	fechaParsed, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		return apperrors.NewBadRequest("Formato de fecha inválido. Use YYYY-MM-DD")
	}

	if _, err := time.Parse("15:04", hora); err != nil {
		return apperrors.NewBadRequest("Formato de hora inválido. Use HH:MM")
	}

	return s.repo.Reagendar(ctx, id, fechaParsed, hora, turno)
}
