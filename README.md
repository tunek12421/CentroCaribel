# Centro Caribel - API Backend

Sistema de gestión para centro de fisioterapia estética en Cochabamba, Bolivia.

## Stack

- Go 1.22+ (net/http con ServeMux mejorado)
- PostgreSQL 15
- JWT (golang-jwt/jwt/v5)
- Docker + Docker Compose

## Requisitos previos

- Docker y Docker Compose
- Go 1.22+ (para desarrollo local)

## Setup rápido con Docker

```bash
# Copiar variables de entorno
cp .env.example .env

# Levantar servicios
docker compose up -d

# La API estará disponible en http://localhost:8080
```

## Desarrollo local

```bash
# Copiar variables de entorno
cp .env.example .env

# Levantar solo PostgreSQL
docker compose up -d postgres

# Ejecutar la API
go run ./cmd/api

# Build
go build -o api ./cmd/api
```

## Endpoints

### Autenticación (público)

| Método | Ruta           | Descripción       |
|--------|----------------|--------------------|
| POST   | /auth/login    | Iniciar sesión     |
| POST   | /auth/refresh  | Refrescar token    |

### Usuarios (solo Administradora)

| Método | Ruta            | Descripción         |
|--------|-----------------|----------------------|
| GET    | /usuarios       | Listar usuarios      |
| POST   | /usuarios       | Crear usuario        |
| GET    | /usuarios/:id   | Obtener usuario      |
| PUT    | /usuarios/:id   | Actualizar usuario   |
| DELETE | /usuarios/:id   | Desactivar usuario   |

### Pacientes

| Método | Ruta            | Descripción         |
|--------|-----------------|----------------------|
| GET    | /pacientes      | Listar pacientes     |
| POST   | /pacientes      | Registrar paciente   |
| GET    | /pacientes/:id  | Obtener paciente     |

### Consentimientos

| Método | Ruta                              | Descripción              |
|--------|-----------------------------------|--------------------------|
| GET    | /pacientes/:id/consentimientos    | Listar consentimientos   |
| POST   | /pacientes/:id/consentimientos    | Registrar consentimiento |

### Historia Clínica

| Método | Ruta                       | Descripción             |
|--------|----------------------------|--------------------------|
| GET    | /pacientes/:id/historia    | Consultar historia       |

### Citas

| Método | Ruta                  | Descripción          |
|--------|-----------------------|-----------------------|
| GET    | /citas                | Listar citas          |
| POST   | /citas                | Agendar cita          |
| PATCH  | /citas/:id/estado     | Cambiar estado        |

### Health Check

| Método | Ruta     | Descripción       |
|--------|----------|--------------------|
| GET    | /health  | Estado del servicio|

## Credenciales iniciales

```
Email:    admin@centrocaribel.com
Password: Admin123!
```

## Roles del sistema

- **Administradora**: Acceso completo
- **Licenciada**: Pacientes, citas, tratamientos
- **Interno**: Solo lectura + registro de asistencia
- **Medico**: Evaluaciones y fichas clínicas

## Estructura del proyecto

```
cmd/api/                        → Punto de entrada
internal/
  domain/                       → Entidades y contratos (interfaces)
  application/                  → Casos de uso / servicios
    auth/                       → Autenticación
    usuario/                    → Gestión de usuarios
    paciente/                   → Gestión de pacientes
    consentimiento/             → Consentimientos informados
    cita/                       → Gestión de citas
    historia/                   → Historias clínicas
  infrastructure/
    config/                     → Configuración desde env
    database/                   → Conexión y migraciones
    jwt/                        → Implementación JWT
    repository/                 → Implementaciones de repositorios
  interfaces/http/
    handler/                    → Controladores HTTP
    middleware/                 → Auth, CORS, logging, recovery
    router/                     → Definición de rutas
    dto/                        → Request/Response DTOs
pkg/
  errors/                       → Errores de aplicación
  response/                     → Respuestas JSON estandarizadas
  validator/                    → Validación de entrada
migrations/                     → Scripts SQL
```
