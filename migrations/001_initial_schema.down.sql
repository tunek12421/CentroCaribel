DROP TRIGGER IF EXISTS tr_historias_updated_at ON historias_clinicas;
DROP TRIGGER IF EXISTS tr_citas_updated_at ON citas;
DROP TRIGGER IF EXISTS tr_pacientes_updated_at ON pacientes;
DROP TRIGGER IF EXISTS tr_usuarios_updated_at ON usuarios;
DROP TRIGGER IF EXISTS tr_roles_updated_at ON roles;
DROP FUNCTION IF EXISTS update_updated_at();

DROP TABLE IF EXISTS historias_clinicas;
DROP TABLE IF EXISTS citas;
DROP TABLE IF EXISTS consentimientos;
DROP TABLE IF EXISTS pacientes;
DROP TABLE IF EXISTS usuarios;
DROP TABLE IF EXISTS roles;

DROP SEQUENCE IF EXISTS historias_numero_seq;
DROP SEQUENCE IF EXISTS pacientes_codigo_seq;
DROP TYPE IF EXISTS turno_cita;
DROP TYPE IF EXISTS estado_cita;
