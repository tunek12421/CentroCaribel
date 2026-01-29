package historia

import (
	"context"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/domain"
	apperrors "github.com/tunek/centro-caribel/pkg/errors"
)

type Service struct {
	repo         domain.HistoriaClinicaRepository
	notaRepo     domain.NotaEvolucionRepository
	pacienteRepo domain.PacienteRepository
}

func NewService(repo domain.HistoriaClinicaRepository, notaRepo domain.NotaEvolucionRepository, pacienteRepo domain.PacienteRepository) *Service {
	return &Service{repo: repo, notaRepo: notaRepo, pacienteRepo: pacienteRepo}
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

func (s *Service) UpdateAntecedentes(ctx context.Context, pacienteID uuid.UUID, antPersonales, antFamiliares, alergias, medicamentos string) (*domain.HistoriaClinica, error) {
	historia, err := s.repo.GetByPacienteID(ctx, pacienteID)
	if err != nil {
		return nil, apperrors.NewNotFound("Historia clínica")
	}

	historia.AntecedentesPersonales = antPersonales
	historia.AntecedentesFamiliares = antFamiliares
	historia.Alergias = alergias
	historia.MedicamentosActuales = medicamentos

	if err := s.repo.UpdateAntecedentes(ctx, historia); err != nil {
		return nil, apperrors.NewInternal("Error al actualizar antecedentes")
	}

	return historia, nil
}

func (s *Service) GetNotas(ctx context.Context, pacienteID uuid.UUID) ([]domain.NotaEvolucion, error) {
	historia, err := s.repo.GetByPacienteID(ctx, pacienteID)
	if err != nil {
		return nil, apperrors.NewNotFound("Historia clínica")
	}

	notas, err := s.notaRepo.GetByHistoriaID(ctx, historia.ID)
	if err != nil {
		return nil, apperrors.NewInternal("Error al obtener notas de evolución")
	}
	if notas == nil {
		notas = []domain.NotaEvolucion{}
	}
	return notas, nil
}

func (s *Service) CreateNota(ctx context.Context, pacienteID uuid.UUID, tipo, contenido string, createdBy uuid.UUID) (*domain.NotaEvolucion, error) {
	historia, err := s.repo.GetByPacienteID(ctx, pacienteID)
	if err != nil {
		return nil, apperrors.NewNotFound("Historia clínica")
	}

	validTipos := map[string]bool{"TRATAMIENTO": true, "EVOLUCION": true, "NOTA": true}
	if !validTipos[tipo] {
		return nil, apperrors.NewBadRequest("Tipo de nota inválido. Use: TRATAMIENTO, EVOLUCION o NOTA")
	}

	nota := &domain.NotaEvolucion{
		ID:         uuid.New(),
		HistoriaID: historia.ID,
		Tipo:       tipo,
		Contenido:  contenido,
		CreatedBy:  createdBy,
	}

	if err := s.notaRepo.Create(ctx, nota); err != nil {
		return nil, apperrors.NewInternal("Error al crear nota de evolución")
	}

	return nota, nil
}
