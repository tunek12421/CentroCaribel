-- Tabla de paquetes de tratamiento
CREATE TABLE paquetes_tratamiento (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    paciente_id UUID NOT NULL REFERENCES pacientes(id),
    tipo_tratamiento VARCHAR(100) NOT NULL,
    total_sesiones INT NOT NULL CHECK (total_sesiones > 0),
    sesiones_completadas INT NOT NULL DEFAULT 0 CHECK (sesiones_completadas >= 0),
    estado VARCHAR(20) NOT NULL DEFAULT 'ACTIVO',
    notas TEXT,
    created_by UUID NOT NULL REFERENCES usuarios(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_paquetes_paciente ON paquetes_tratamiento(paciente_id);
CREATE INDEX idx_paquetes_estado ON paquetes_tratamiento(estado);

CREATE TRIGGER tr_paquetes_updated_at BEFORE UPDATE ON paquetes_tratamiento
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Agregar paquete_id nullable a citas
ALTER TABLE citas ADD COLUMN paquete_id UUID REFERENCES paquetes_tratamiento(id);
CREATE INDEX idx_citas_paquete ON citas(paquete_id);

-- Índice único parcial para prevenir doble reserva en misma fecha+hora
CREATE UNIQUE INDEX idx_citas_fecha_hora_unique
    ON citas (fecha, hora)
    WHERE estado NOT IN ('CANCELADA', 'REAGENDADA');
