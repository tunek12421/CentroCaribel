-- HU-PAC-003: Historia clínica completa con antecedentes, evolución y tratamientos
CREATE EXTENSION IF NOT EXISTS unaccent;

-- Agregar campos de antecedentes a historias_clinicas
ALTER TABLE historias_clinicas
    ADD COLUMN antecedentes_personales TEXT NOT NULL DEFAULT '',
    ADD COLUMN antecedentes_familiares TEXT NOT NULL DEFAULT '',
    ADD COLUMN alergias TEXT NOT NULL DEFAULT '',
    ADD COLUMN medicamentos_actuales TEXT NOT NULL DEFAULT '';

-- Tabla de notas de evolución (tratamientos, evolución, notas generales)
CREATE TABLE notas_evolucion (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    historia_id UUID NOT NULL REFERENCES historias_clinicas(id),
    tipo VARCHAR(20) NOT NULL CHECK (tipo IN ('TRATAMIENTO', 'EVOLUCION', 'NOTA')),
    contenido TEXT NOT NULL,
    created_by UUID NOT NULL REFERENCES usuarios(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notas_historia ON notas_evolucion(historia_id);
CREATE INDEX idx_notas_tipo ON notas_evolucion(tipo);
CREATE INDEX idx_notas_created_at ON notas_evolucion(created_at DESC);

-- Índice para búsqueda de pacientes por nombre (lower)
CREATE INDEX idx_pacientes_nombre_lower ON pacientes (lower(nombre_completo));
