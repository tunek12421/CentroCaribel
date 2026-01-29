package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tunek/centro-caribel/internal/application/auth"
)

type Service struct {
	secret          []byte
	expHours        int
	refreshExpHours int
}

func NewService(secret string, expHours, refreshExpHours int) *Service {
	return &Service{
		secret:          []byte(secret),
		expHours:        expHours,
		refreshExpHours: refreshExpHours,
	}
}

func (s *Service) GenerateToken(userID, rolNombre string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    userID,
		"rol_nombre": rolNombre,
		"exp":        time.Now().Add(time.Duration(s.expHours) * time.Hour).Unix(),
		"iat":        time.Now().Unix(),
		"type":       "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *Service) GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(s.refreshExpHours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"type":    "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *Service) ValidateToken(tokenStr string) (*auth.Claims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("token inválido")
	}

	tokenType, _ := claims["type"].(string)
	if tokenType != "access" {
		return nil, fmt.Errorf("tipo de token inválido")
	}

	return &auth.Claims{
		UserID:    claims["user_id"].(string),
		RolNombre: claims["rol_nombre"].(string),
	}, nil
}

func (s *Service) ValidateRefreshToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("token inválido")
	}

	tokenType, _ := claims["type"].(string)
	if tokenType != "refresh" {
		return "", fmt.Errorf("tipo de token inválido")
	}

	userID, _ := claims["user_id"].(string)
	return userID, nil
}
