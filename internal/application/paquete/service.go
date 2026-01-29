package paquete

import (
	"context"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/domain"
	apperrors "github.com/tunek/centro-caribel/pkg/errors"
)

type Service struct {
	repo         domain.PaqueteRepository
	pacienteRepo domain.PacienteRepository
}

func NewService(repo domain.PaqueteRepository, pacienteRepo domain.PacienteRepository) *Service {
	return &Service{repo: repo, pacienteRepo: pacienteRepo}
}

func (s *Service) Create(ctx context.Context, pacienteID uuid.UUID, tipoTratamiento string, totalSesiones int, notas string, createdBy uuid.UUID) (*domain.PaqueteTratamiento, error) {
	if _, err := s.pacienteRepo.GetByID(ctx, pacienteID); err != nil {
		return nil, apperrors.NewNotFound("Paciente")
	}

	if totalSesiones < 1 {
		return nil, apperrors.NewBadRequest("El total de sesiones debe ser al menos 1")
	}

	p := &domain.PaqueteTratamiento{
		ID:              uuid.New(),
		PacienteID:      pacienteID,
		TipoTratamiento: tipoTratamiento,
		TotalSesiones:   totalSesiones,
		Estado:          domain.PaqueteActivo,
		Notas:           notas,
		CreatedBy:       createdBy,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, apperrors.NewInternal("Error al crear el paquete de tratamiento")
	}

	return p, nil
}

func (s *Service) GetByPacienteID(ctx context.Context, pacienteID uuid.UUID) ([]domain.PaqueteTratamiento, error) {
	return s.repo.GetByPacienteID(ctx, pacienteID)
}

func (s *Service) GetActivosByPaciente(ctx context.Context, pacienteID uuid.UUID) ([]domain.PaqueteTratamiento, error) {
	return s.repo.GetActivosByPaciente(ctx, pacienteID)
}

func (s *Service) CancelPaquete(ctx context.Context, id uuid.UUID) error {
	paq, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFound("Paquete de tratamiento")
	}
	if paq.Estado != domain.PaqueteActivo {
		return apperrors.NewBadRequest("Solo se pueden cancelar paquetes activos")
	}
	return s.repo.UpdateEstado(ctx, id, domain.PaqueteCancelado)
}
