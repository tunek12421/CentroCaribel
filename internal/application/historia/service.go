package historia

import (
	"context"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/domain"
	apperrors "github.com/tunek/centro-caribel/pkg/errors"
)

type Service struct {
	repo         domain.HistoriaClinicaRepository
	pacienteRepo domain.PacienteRepository
}

func NewService(repo domain.HistoriaClinicaRepository, pacienteRepo domain.PacienteRepository) *Service {
	return &Service{repo: repo, pacienteRepo: pacienteRepo}
}

func (s *Service) GetByPacienteID(ctx context.Context, pacienteID uuid.UUID) (*domain.HistoriaClinica, error) {
	if _, err := s.pacienteRepo.GetByID(ctx, pacienteID); err != nil {
		return nil, apperrors.NewNotFound("Paciente")
	}

	historia, err := s.repo.GetByPacienteID(ctx, pacienteID)
	if err != nil {
		return nil, apperrors.NewNotFound("Historia clínica")
	}

	return historia, nil
}
