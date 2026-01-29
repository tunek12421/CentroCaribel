package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/domain"
	apperrors "github.com/tunek/centro-caribel/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type JWTService interface {
	GenerateToken(userID, rolNombre string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateToken(token string) (*Claims, error)
	ValidateRefreshToken(token string) (string, error)
}

type Claims struct {
	UserID    string `json:"user_id"`
	RolNombre string `json:"rol_nombre"`
}

type Service struct {
	userRepo domain.UsuarioRepository
	rolRepo  domain.RolRepository
	jwt      JWTService
}

func NewService(userRepo domain.UsuarioRepository, rolRepo domain.RolRepository, jwt JWTService) *Service {
	return &Service{userRepo: userRepo, rolRepo: rolRepo, jwt: jwt}
}

func (s *Service) Login(ctx context.Context, email, password string) (token, refreshToken string, err error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", apperrors.NewUnauthorized("Credenciales inválidas")
	}

	if !user.Activo {
		return "", "", apperrors.NewUnauthorized("Usuario desactivado")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", apperrors.NewUnauthorized("Credenciales inválidas")
	}

	rol, err := s.rolRepo.GetByID(ctx, user.RolID)
	if err != nil {
		return "", "", apperrors.NewInternal("Error al obtener rol")
	}

	token, err = s.jwt.GenerateToken(user.ID.String(), rol.Nombre)
	if err != nil {
		return "", "", apperrors.NewInternal("Error al generar token")
	}

	refreshToken, err = s.jwt.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return "", "", apperrors.NewInternal("Error al generar refresh token")
	}

	return token, refreshToken, nil
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (newToken, newRefresh string, err error) {
	userIDStr, err := s.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", apperrors.NewUnauthorized("Refresh token inválido")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", "", apperrors.NewUnauthorized("Token inválido")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", "", apperrors.NewUnauthorized("Usuario no encontrado")
	}

	if !user.Activo {
		return "", "", apperrors.NewUnauthorized("Usuario desactivado")
	}

	rol, err := s.rolRepo.GetByID(ctx, user.RolID)
	if err != nil {
		return "", "", apperrors.NewInternal("Error al obtener rol")
	}

	newToken, err = s.jwt.GenerateToken(user.ID.String(), rol.Nombre)
	if err != nil {
		return "", "", apperrors.NewInternal("Error al generar token")
	}

	newRefresh, err = s.jwt.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return "", "", apperrors.NewInternal("Error al generar refresh token")
	}

	return newToken, newRefresh, nil
}
