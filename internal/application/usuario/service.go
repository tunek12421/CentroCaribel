package usuario

import (
	"context"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/domain"
	apperrors "github.com/tunek/centro-caribel/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo    domain.UsuarioRepository
	rolRepo domain.RolRepository
}

func NewService(repo domain.UsuarioRepository, rolRepo domain.RolRepository) *Service {
	return &Service{repo: repo, rolRepo: rolRepo}
}

func (s *Service) Create(ctx context.Context, nombreCompleto, email, password string, rolID uuid.UUID) (*domain.Usuario, error) {
	if _, err := s.rolRepo.GetByID(ctx, rolID); err != nil {
		return nil, apperrors.NewBadRequest("El rol especificado no existe")
	}

	if existing, _ := s.repo.GetByEmail(ctx, email); existing != nil {
		return nil, apperrors.NewConflict("Ya existe un usuario con ese email")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.NewInternal("Error al procesar la contraseña")
	}

	user := &domain.Usuario{
		ID:             uuid.New(),
		NombreCompleto: nombreCompleto,
		Email:          email,
		PasswordHash:   string(hash),
		RolID:          rolID,
		Activo:         true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, apperrors.NewInternal("Error al crear el usuario")
	}

	return user, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.Usuario, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFound("Usuario")
	}

	rol, err := s.rolRepo.GetByID(ctx, user.RolID)
	if err == nil {
		user.Rol = rol
	}

	return user, nil
}

func (s *Service) GetAll(ctx context.Context, page, perPage int) ([]domain.Usuario, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage
	return s.repo.GetAll(ctx, offset, perPage)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, nombreCompleto, email string, rolID uuid.UUID, activo *bool) (*domain.Usuario, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFound("Usuario")
	}

	if nombreCompleto != "" {
		user.NombreCompleto = nombreCompleto
	}
	if email != "" && email != user.Email {
		if existing, _ := s.repo.GetByEmail(ctx, email); existing != nil {
			return nil, apperrors.NewConflict("Ya existe un usuario con ese email")
		}
		user.Email = email
	}
	if rolID != uuid.Nil {
		if _, err := s.rolRepo.GetByID(ctx, rolID); err != nil {
			return nil, apperrors.NewBadRequest("El rol especificado no existe")
		}
		user.RolID = rolID
	}
	if activo != nil {
		user.Activo = *activo
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, apperrors.NewInternal("Error al actualizar el usuario")
	}

	return user, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return apperrors.NewNotFound("Usuario")
	}
	return s.repo.Delete(ctx, id)
}
