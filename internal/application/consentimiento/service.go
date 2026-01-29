package consentimiento

import (
	"context"
	"encoding/base64"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/domain"
	apperrors "github.com/tunek/centro-caribel/pkg/errors"
)

type Service struct {
	repo        domain.ConsentimientoRepository
	pacienteRepo domain.PacienteRepository
}

func NewService(repo domain.ConsentimientoRepository, pacienteRepo domain.PacienteRepository) *Service {
	return &Service{repo: repo, pacienteRepo: pacienteRepo}
}

func (s *Service) Create(ctx context.Context, pacienteID uuid.UUID, firmaB64, contenido string, autorizaFotos bool, registradoPor uuid.UUID) (*domain.Consentimiento, error) {
	if _, err := s.pacienteRepo.GetByID(ctx, pacienteID); err != nil {
		return nil, apperrors.NewNotFound("Paciente")
	}

	var firma []byte
	if firmaB64 != "" {
		var err error
		firma, err = base64.StdEncoding.DecodeString(firmaB64)
		if err != nil {
			return nil, apperrors.NewBadRequest("La firma digital no es un base64 válido")
		}
	}

	cons := &domain.Consentimiento{
		ID:            uuid.New(),
		PacienteID:    pacienteID,
		FirmaDigital:  firma,
		AutorizaFotos: autorizaFotos,
		Contenido:     contenido,
		RegistradoPor: registradoPor,
	}

	if err := s.repo.Create(ctx, cons); err != nil {
		return nil, apperrors.NewInternal("Error al registrar el consentimiento")
	}

	return cons, nil
}

func (s *Service) GetByPacienteID(ctx context.Context, pacienteID uuid.UUID) ([]domain.Consentimiento, error) {
	if _, err := s.pacienteRepo.GetByID(ctx, pacienteID); err != nil {
		return nil, apperrors.NewNotFound("Paciente")
	}
	return s.repo.GetByPacienteID(ctx, pacienteID)
}
