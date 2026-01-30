export interface ApiResponse<T> {
  success: boolean;
  data: T;
  error?: { message: string; detail?: string };
  meta?: PaginationMeta;
}

export interface PaginationMeta {
  page: number;
  per_page: number;
  total: number;
  total_pages: number;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  refresh_token: string;
  expires_in: number;
}

export interface Rol {
  id: string;
  nombre: string;
  descripcion: string;
  permisos: Record<string, string[]>;
  activo: boolean;
}

export interface Usuario {
  id: string;
  nombre_completo: string;
  email: string;
  rol_id: string;
  rol?: Rol;
  activo: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateUsuarioRequest {
  nombre_completo: string;
  email: string;
  password: string;
  rol_id: string;
}

export interface UpdateUsuarioRequest {
  nombre_completo?: string;
  email?: string;
  rol_id?: string;
  activo?: boolean;
}

export interface Paciente {
  id: string;
  codigo: string;
  nombre_completo: string;
  ci: string;
  fecha_nacimiento: string;
  celular: string;
  direccion?: string;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface CreatePacienteRequest {
  nombre_completo: string;
  ci: string;
  fecha_nacimiento: string;
  celular: string;
  direccion?: string;
}

export interface Consentimiento {
  id: string;
  paciente_id: string;
  fecha_firma: string;
  firma_digital?: string;
  autoriza_fotos: boolean;
  contenido: string;
  registrado_por: string;
  created_at: string;
}

export interface CreateConsentimientoRequest {
  firma_digital?: string;
  autoriza_fotos: boolean;
  contenido: string;
}

export interface HistoriaClinica {
  id: string;
  paciente_id: string;
  numero_historia: string;
  estado: string;
  antecedentes_personales: string;
  antecedentes_familiares: string;
  alergias: string;
  medicamentos_actuales: string;
  created_at: string;
  updated_at: string;
}

export interface NotaEvolucion {
  id: string;
  historia_id: string;
  tipo: 'TRATAMIENTO' | 'EVOLUCION' | 'NOTA';
  contenido: string;
  created_by: string;
  created_at: string;
}

export interface UpdateAntecedentesRequest {
  antecedentes_personales: string;
  antecedentes_familiares: string;
  alergias: string;
  medicamentos_actuales: string;
}

export interface CreateNotaRequest {
  tipo: 'TRATAMIENTO' | 'EVOLUCION' | 'NOTA';
  contenido: string;
}

export type EstadoCita =
  | 'NUEVA'
  | 'AGENDADA'
  | 'CONFIRMADA'
  | 'ATENDIDA'
  | 'NO_ASISTIO'
  | 'CANCELADA'
  | 'REAGENDADA';

export type TurnoCita = 'AM' | 'PM';

export interface Cita {
  id: string;
  paciente_id: string;
  paciente_nombre?: string;
  fecha: string;
  hora: string;
  tipo_tratamiento: string;
  estado: EstadoCita;
  turno: TurnoCita;
  observaciones?: string;
  paquete_id?: string;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface CreateCitaRequest {
  paciente_id: string;
  fecha: string;
  hora: string;
  tipo_tratamiento: string;
  turno: TurnoCita;
  observaciones?: string;
  paquete_id?: string;
}

export interface PaqueteTratamiento {
  id: string;
  paciente_id: string;
  tipo_tratamiento: string;
  total_sesiones: number;
  sesiones_completadas: number;
  estado: 'ACTIVO' | 'COMPLETADO' | 'CANCELADO';
  notas?: string;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface CreatePaqueteRequest {
  paciente_id: string;
  tipo_tratamiento: string;
  total_sesiones: number;
  notas?: string;
}

export const TRANSICIONES_VALIDAS: Record<EstadoCita, EstadoCita[]> = {
  NUEVA: ['AGENDADA', 'CANCELADA'],
  AGENDADA: ['CONFIRMADA', 'CANCELADA', 'REAGENDADA'],
  CONFIRMADA: ['ATENDIDA', 'NO_ASISTIO', 'CANCELADA'],
  REAGENDADA: ['AGENDADA', 'CANCELADA'],
  ATENDIDA: [],
  NO_ASISTIO: [],
  CANCELADA: [],
};

export const ESTADO_COLORS: Record<EstadoCita, string> = {
  NUEVA: 'bg-blue-100 text-blue-800',
  AGENDADA: 'bg-yellow-100 text-yellow-800',
  CONFIRMADA: 'bg-green-100 text-green-800',
  ATENDIDA: 'bg-emerald-100 text-emerald-800',
  NO_ASISTIO: 'bg-red-100 text-red-800',
  CANCELADA: 'bg-gray-100 text-gray-800',
  REAGENDADA: 'bg-purple-100 text-purple-800',
};

export const ESTADO_LABELS: Record<EstadoCita, string> = {
  NUEVA: 'Nueva',
  AGENDADA: 'Agendada',
  CONFIRMADA: 'Confirmada',
  ATENDIDA: 'Atendida',
  NO_ASISTIO: 'No Asistió',
  CANCELADA: 'Cancelada',
  REAGENDADA: 'Reagendada',
};
