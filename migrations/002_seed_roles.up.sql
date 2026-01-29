-- Seed: Roles iniciales del sistema
INSERT INTO roles (nombre, descripcion, permisos) VALUES
(
    'Administradora',
    'Acceso completo al sistema',
    '{
        "usuarios": ["crear", "leer", "actualizar", "eliminar"],
        "pacientes": ["crear", "leer", "actualizar", "eliminar"],
        "citas": ["crear", "leer", "actualizar", "eliminar"],
        "consentimientos": ["crear", "leer"],
        "historias": ["crear", "leer", "actualizar"],
        "tratamientos": ["crear", "leer", "actualizar", "eliminar"],
        "reportes": ["leer"]
    }'
),
(
    'Licenciada',
    'Gestión de pacientes, citas y tratamientos',
    '{
        "pacientes": ["crear", "leer", "actualizar"],
        "citas": ["crear", "leer", "actualizar"],
        "consentimientos": ["crear", "leer"],
        "historias": ["crear", "leer", "actualizar"],
        "tratamientos": ["crear", "leer", "actualizar"]
    }'
),
(
    'Interno',
    'Acceso de solo lectura y registro de asistencia',
    '{
        "pacientes": ["leer"],
        "citas": ["leer", "actualizar"],
        "historias": ["leer"]
    }'
),
(
    'Medico',
    'Evaluaciones y fichas clínicas',
    '{
        "pacientes": ["leer"],
        "citas": ["leer"],
        "historias": ["leer", "actualizar"],
        "tratamientos": ["crear", "leer", "actualizar"]
    }'
);
