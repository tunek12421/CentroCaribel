package paciente

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/domain"
	apperrors "github.com/tunek/centro-caribel/pkg/errors"
)

type Service struct {
	repo        domain.PacienteRepository
	historiaRepo domain.HistoriaClinicaRepository
}

func NewService(repo domain.PacienteRepository, historiaRepo domain.HistoriaClinicaRepository) *Service {
	return &Service{repo: repo, historiaRepo: historiaRepo}
}

func (s *Service) Create(ctx context.Context, nombre, ci, fechaNac, celular, direccion string, createdBy uuid.UUID) (*domain.Paciente, error) {
	if existing, _ := s.repo.GetByCI(ctx, ci); existing != nil {
		return nil, apperrors.NewConflict("Ya existe un paciente con ese CI")
	}

	fecha, err := time.Parse("2006-01-02", fechaNac)
	if err != nil {
		return nil, apperrors.NewBadRequest("Formato de fecha inválido. Use YYYY-MM-DD")
	}

	codigo, err := s.repo.NextCodigo(ctx)
	if err != nil {
		return nil, apperrors.NewInternal("Error al generar código de paciente")
	}

	pac := &domain.Paciente{
		ID:              uuid.New(),
		Codigo:          codigo,
		NombreCompleto:  nombre,
		CI:              ci,
		FechaNacimiento: fecha,
		Celular:         celular,
		Direccion:       direccion,
		CreatedBy:       createdBy,
	}

	if err := s.repo.Create(ctx, pac); err != nil {
		return nil, apperrors.NewInternal("Error al crear el paciente")
	}

	// Crear historia clínica automáticamente
	numHistoria, err := s.historiaRepo.NextNumero(ctx)
	if err == nil {
		historia := &domain.HistoriaClinica{
			ID:             uuid.New(),
			PacienteID:     pac.ID,
			NumeroHistoria: numHistoria,
			Estado:         "ACTIVA",
		}
		_ = s.historiaRepo.Create(ctx, historia)
	}

	return pac, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.Paciente, error) {
	pac, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFound("Paciente")
	}
	return pac, nil
}

func (s *Service) GetAll(ctx context.Context, page, perPage int) ([]domain.Paciente, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage
	return s.repo.GetAll(ctx, offset, perPage)
}
