package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/application/auth"
	"github.com/tunek/centro-caribel/internal/application/cita"
	"github.com/tunek/centro-caribel/internal/application/consentimiento"
	"github.com/tunek/centro-caribel/internal/application/historia"
	"github.com/tunek/centro-caribel/internal/application/paciente"
	"github.com/tunek/centro-caribel/internal/application/usuario"
	"github.com/tunek/centro-caribel/internal/domain"
	"github.com/tunek/centro-caribel/internal/infrastructure/config"
	"github.com/tunek/centro-caribel/internal/infrastructure/database"
	jwtinfra "github.com/tunek/centro-caribel/internal/infrastructure/jwt"
	"github.com/tunek/centro-caribel/internal/infrastructure/repository"
	"github.com/tunek/centro-caribel/internal/interfaces/http/handler"
	"github.com/tunek/centro-caribel/internal/interfaces/http/router"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Load()

	db, err := database.NewConnection(cfg.DB)
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}
	defer db.Close()

	migrationsDir := "migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		migrationsDir = "/migrations"
	}
	if err := database.RunMigrations(db, migrationsDir); err != nil {
		log.Fatalf("Error ejecutando migraciones: %v", err)
	}

	// Repositorios
	rolRepo := repository.NewRolRepository(db)
	usuarioRepo := repository.NewUsuarioRepository(db)
	pacienteRepo := repository.NewPacienteRepository(db)
	consentimientoRepo := repository.NewConsentimientoRepository(db)
	citaRepo := repository.NewCitaRepository(db)
	historiaRepo := repository.NewHistoriaClinicaRepository(db)

	// JWT
	jwtSvc := jwtinfra.NewService(cfg.JWT.Secret, cfg.JWT.ExpirationHours, cfg.JWT.RefreshExpirationHrs)

	// Servicios
	authSvc := auth.NewService(usuarioRepo, rolRepo, jwtSvc)
	usuarioSvc := usuario.NewService(usuarioRepo, rolRepo)
	pacienteSvc := paciente.NewService(pacienteRepo, historiaRepo)
	consentimientoSvc := consentimiento.NewService(consentimientoRepo, pacienteRepo)
	citaSvc := cita.NewService(citaRepo, pacienteRepo)
	historiaSvc := historia.NewService(historiaRepo, pacienteRepo)

	// Seed admin
	seedAdmin(usuarioRepo, rolRepo, cfg.Admin)

	// Handlers
	handlers := router.Handlers{
		Auth:           handler.NewAuthHandler(authSvc),
		Usuario:        handler.NewUsuarioHandler(usuarioSvc),
		Paciente:       handler.NewPacienteHandler(pacienteSvc),
		Consentimiento: handler.NewConsentimientoHandler(consentimientoSvc),
		Cita:           handler.NewCitaHandler(citaSvc),
		Historia:       handler.NewHistoriaHandler(historiaSvc),
	}

	mux := router.New(handlers, jwtSvc)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Servidor iniciado en puerto %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error en el servidor: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Apagando servidor...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error durante el apagado: %v", err)
	}
	log.Println("Servidor detenido")
}

func seedAdmin(userRepo *repository.UsuarioRepository, rolRepo *repository.RolRepository, adminCfg config.AdminConfig) {
	ctx := context.Background()

	if existing, _ := userRepo.GetByEmail(ctx, adminCfg.Email); existing != nil {
		return
	}

	rol, err := rolRepo.GetByNombre(ctx, "Administradora")
	if err != nil {
		log.Println("Rol Administradora no encontrado, omitiendo seed de admin")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(adminCfg.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error generando hash para admin: %v", err)
		return
	}

	user := &domain.Usuario{
		ID:             uuid.New(),
		NombreCompleto: "Administradora del Sistema",
		Email:          adminCfg.Email,
		PasswordHash:   string(hash),
		RolID:          rol.ID,
		Activo:         true,
	}

	if err := userRepo.Create(ctx, user); err != nil {
		log.Printf("Error creando usuario admin: %v", err)
		return
	}

	log.Printf("Usuario administrador creado: %s", adminCfg.Email)
}
