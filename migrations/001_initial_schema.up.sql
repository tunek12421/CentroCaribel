-- Centro Caribel - Initial Schema
-- Sprint 1: Roles, Usuarios, Pacientes, Consentimientos, Citas, Historias Clínicas

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enum para estados de cita
CREATE TYPE estado_cita AS ENUM (
    'NUEVA',
    'AGENDADA',
    'CONFIRMADA',
    'ATENDIDA',
    'NO_ASISTIO',
    'CANCELADA',
    'REAGENDADA'
);

-- Enum para turno
CREATE TYPE turno_cita AS ENUM ('AM', 'PM');

-- Tabla de roles
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nombre VARCHAR(50) NOT NULL UNIQUE,
    descripcion TEXT,
    permisos JSONB NOT NULL DEFAULT '{}',
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Tabla de usuarios
CREATE TABLE usuarios (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nombre_completo VARCHAR(150) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    rol_id UUID NOT NULL REFERENCES roles(id),
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_usuarios_email ON usuarios(email);
CREATE INDEX idx_usuarios_rol_id ON usuarios(rol_id);

-- Tabla de pacientes
CREATE TABLE pacientes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    codigo VARCHAR(20) NOT NULL UNIQUE,
    nombre_completo VARCHAR(150) NOT NULL,
    ci VARCHAR(20) NOT NULL UNIQUE,
    fecha_nacimiento DATE NOT NULL,
    celular VARCHAR(20) NOT NULL,
    direccion TEXT,
    created_by UUID NOT NULL REFERENCES usuarios(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pacientes_codigo ON pacientes(codigo);
CREATE INDEX idx_pacientes_ci ON pacientes(ci);

-- Secuencia para código de paciente
CREATE SEQUENCE pacientes_codigo_seq START 1;

-- Tabla de consentimientos
CREATE TABLE consentimientos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    paciente_id UUID NOT NULL REFERENCES pacientes(id),
    fecha_firma TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    firma_digital BYTEA,
    autoriza_fotos BOOLEAN NOT NULL DEFAULT FALSE,
    contenido TEXT NOT NULL,
    registrado_por UUID NOT NULL REFERENCES usuarios(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_consentimientos_paciente ON consentimientos(paciente_id);

-- Tabla de citas
CREATE TABLE citas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    paciente_id UUID NOT NULL REFERENCES pacientes(id),
    fecha DATE NOT NULL,
    hora TIME NOT NULL,
    tipo_tratamiento VARCHAR(100) NOT NULL,
    estado estado_cita NOT NULL DEFAULT 'NUEVA',
    turno turno_cita NOT NULL,
    observaciones TEXT,
    created_by UUID NOT NULL REFERENCES usuarios(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_citas_paciente ON citas(paciente_id);
CREATE INDEX idx_citas_fecha ON citas(fecha);
CREATE INDEX idx_citas_estado ON citas(estado);

-- Tabla de historias clínicas
CREATE TABLE historias_clinicas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    paciente_id UUID NOT NULL UNIQUE REFERENCES pacientes(id),
    numero_historia VARCHAR(20) NOT NULL UNIQUE,
    estado VARCHAR(20) NOT NULL DEFAULT 'ACTIVA',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_historias_paciente ON historias_clinicas(paciente_id);

-- Secuencia para número de historia
CREATE SEQUENCE historias_numero_seq START 1;

-- Trigger para updated_at
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_roles_updated_at BEFORE UPDATE ON roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER tr_usuarios_updated_at BEFORE UPDATE ON usuarios
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER tr_pacientes_updated_at BEFORE UPDATE ON pacientes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER tr_citas_updated_at BEFORE UPDATE ON citas
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER tr_historias_updated_at BEFORE UPDATE ON historias_clinicas
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
