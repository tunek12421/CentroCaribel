package router

import (
	"net/http"

	"github.com/tunek/centro-caribel/internal/application/auth"
	"github.com/tunek/centro-caribel/internal/interfaces/http/handler"
	"github.com/tunek/centro-caribel/internal/interfaces/http/middleware"
)

type Handlers struct {
	Auth           *handler.AuthHandler
	Usuario        *handler.UsuarioHandler
	Paciente       *handler.PacienteHandler
	Consentimiento *handler.ConsentimientoHandler
	Cita           *handler.CitaHandler
	Historia       *handler.HistoriaHandler
	Rol            *handler.RolHandler
}

func New(h Handlers, jwtSvc auth.JWTService) http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Auth (público)
	mux.HandleFunc("POST /auth/login", h.Auth.Login)
	mux.HandleFunc("POST /auth/refresh", h.Auth.Refresh)

	// Rutas protegidas
	authMw := middleware.AuthMiddleware(jwtSvc)
	adminOnly := middleware.RequireRoles("Administradora")
	staffRoles := middleware.RequireRoles("Administradora", "Licenciada")
	allRoles := middleware.RequireRoles("Administradora", "Licenciada", "Interno", "Medico")

	// Roles (autenticado)
	mux.Handle("GET /roles", authMw(allRoles(http.HandlerFunc(h.Rol.GetAll))))

	// Usuarios (solo admin)
	mux.Handle("GET /usuarios", authMw(adminOnly(http.HandlerFunc(h.Usuario.GetAll))))
	mux.Handle("POST /usuarios", authMw(adminOnly(http.HandlerFunc(h.Usuario.Create))))
	mux.Handle("GET /usuarios/{id}", authMw(adminOnly(http.HandlerFunc(h.Usuario.GetByID))))
	mux.Handle("PUT /usuarios/{id}", authMw(adminOnly(http.HandlerFunc(h.Usuario.Update))))
	mux.Handle("DELETE /usuarios/{id}", authMw(adminOnly(http.HandlerFunc(h.Usuario.Delete))))

	// Pacientes
	mux.Handle("GET /pacientes", authMw(allRoles(http.HandlerFunc(h.Paciente.GetAll))))
	mux.Handle("POST /pacientes", authMw(staffRoles(http.HandlerFunc(h.Paciente.Create))))
	mux.Handle("GET /pacientes/{id}", authMw(allRoles(http.HandlerFunc(h.Paciente.GetByID))))

	// Consentimientos
	mux.Handle("GET /pacientes/{id}/consentimientos", authMw(allRoles(http.HandlerFunc(h.Consentimiento.GetByPaciente))))
	mux.Handle("POST /pacientes/{id}/consentimientos", authMw(staffRoles(http.HandlerFunc(h.Consentimiento.Create))))

	// Historia clínica
	mux.Handle("GET /pacientes/{id}/historia", authMw(allRoles(http.HandlerFunc(h.Historia.GetByPaciente))))

	// Citas
	mux.Handle("GET /citas", authMw(allRoles(http.HandlerFunc(h.Cita.GetAll))))
	mux.Handle("POST /citas", authMw(staffRoles(http.HandlerFunc(h.Cita.Create))))
	mux.Handle("PATCH /citas/{id}/estado", authMw(staffRoles(http.HandlerFunc(h.Cita.UpdateEstado))))

	// Aplicar middlewares globales
	var handler http.Handler = mux
	handler = middleware.CORS(handler)
	handler = middleware.Logging(handler)
	handler = middleware.Recovery(handler)

	return handler
}
