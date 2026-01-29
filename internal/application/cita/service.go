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
	paqueteRepo  domain.PaqueteRepository
}

func NewService(repo domain.CitaRepository, pacienteRepo domain.PacienteRepository, paqueteRepo domain.PaqueteRepository) *Service {
	return &Service{repo: repo, pacienteRepo: pacienteRepo, paqueteRepo: paqueteRepo}
}

func validarHorarioAtencion(fecha time.Time, hora string) error {
	weekday := fecha.Weekday()

	if weekday == time.Sunday {
		return apperrors.NewBadRequest("No se atiende los domingos")
	}

	horaTime, err := time.Parse("15:04", hora)
	if err != nil {
		return apperrors.NewBadRequest("Formato de hora inválido")
	}
	totalMinutes := horaTime.Hour()*60 + horaTime.Minute()

	if weekday == time.Saturday {
		if totalMinutes >= 12*60 {
			return apperrors.NewBadRequest("Los sábados se atiende solo hasta las 12:00")
		}
		return nil
	}

	// Lunes a viernes hasta las 20:00
	if totalMinutes >= 20*60 {
		return apperrors.NewBadRequest("De lunes a viernes se atiende hasta las 20:00")
	}

	return nil
}

func (s *Service) Create(ctx context.Context, pacienteID uuid.UUID, fecha, hora, tipoTratamiento string, turno domain.TurnoCita, observaciones string, paqueteID *uuid.UUID, createdBy uuid.UUID) (*domain.Cita, error) {
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

	if err := validarHorarioAtencion(fechaParsed, hora); err != nil {
		return nil, err
	}

	exists, err := s.repo.ExistsByFechaHora(ctx, fechaParsed, hora, nil)
	if err != nil {
		return nil, apperrors.NewInternal("Error verificando disponibilidad")
	}
	if exists {
		return nil, apperrors.NewConflict("Ya existe una cita agendada en esa fecha y hora")
	}

	if paqueteID != nil {
		paq, err := s.paqueteRepo.GetByID(ctx, *paqueteID)
		if err != nil {
			return nil, apperrors.NewNotFound("Paquete de tratamiento")
		}
		if paq.PacienteID != pacienteID {
			return nil, apperrors.NewBadRequest("El paquete no pertenece al paciente")
		}
		if paq.Estado != domain.PaqueteActivo {
			return nil, apperrors.NewBadRequest("El paquete no está activo")
		}
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
		PaqueteID:       paqueteID,
		CreatedBy:       createdBy,
	}

	if err := s.repo.Create(ctx, c); err != nil {
		return nil, apperrors.NewInternal("Error al crear la cita")
	}

	return c, nil
}

func (s *Service) GetAll(ctx context.Context, page, perPage int, fecha *time.Time, turno *domain.TurnoCita, estado *domain.EstadoCita) ([]domain.Cita, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage
	return s.repo.GetAllFiltered(ctx, offset, perPage, fecha, turno, estado)
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

	if err := s.repo.UpdateEstado(ctx, id, nuevoEstado); err != nil {
		return apperrors.NewInternal("Error actualizando estado")
	}

	// Al marcar como ATENDIDA, incrementar sesiones del paquete si existe
	if nuevoEstado == domain.EstadoAtendida && c.PaqueteID != nil {
		_ = s.paqueteRepo.IncrementSesiones(ctx, *c.PaqueteID)
		paq, err := s.paqueteRepo.GetByID(ctx, *c.PaqueteID)
		if err == nil && paq.SesionesCompletadas >= paq.TotalSesiones {
			_ = s.paqueteRepo.UpdateEstado(ctx, paq.ID, domain.PaqueteCompletado)
		}
	}

	return nil
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

	if err := validarHorarioAtencion(fechaParsed, hora); err != nil {
		return err
	}

	exists, err := s.repo.ExistsByFechaHora(ctx, fechaParsed, hora, &id)
	if err != nil {
		return apperrors.NewInternal("Error verificando disponibilidad")
	}
	if exists {
		return apperrors.NewConflict("Ya existe una cita agendada en esa fecha y hora")
	}

	return s.repo.Reagendar(ctx, id, fechaParsed, hora, turno)
}
