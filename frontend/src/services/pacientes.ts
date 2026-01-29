import api from '../lib/axios';
import type {
  ApiResponse,
  Paciente,
  CreatePacienteRequest,
  Consentimiento,
  CreateConsentimientoRequest,
  HistoriaClinica,
  NotaEvolucion,
  UpdateAntecedentesRequest,
  CreateNotaRequest,
} from '../types';

export const pacientesService = {
  getAll: async (page = 1, perPage = 20, query?: string) => {
    const res = await api.get<ApiResponse<Paciente[]>>('/pacientes', {
      params: { page, per_page: perPage, ...(query ? { q: query } : {}) },
    });
    return res.data;
  },

  getById: async (id: string) => {
    const res = await api.get<ApiResponse<Paciente>>(`/pacientes/${id}`);
    return res.data;
  },

  create: async (data: CreatePacienteRequest) => {
    const res = await api.post<ApiResponse<Paciente>>('/pacientes', data);
    return res.data;
  },

  getConsentimientos: async (pacienteId: string) => {
    const res = await api.get<ApiResponse<Consentimiento[]>>(
      `/pacientes/${pacienteId}/consentimientos`
    );
    return res.data;
  },

  createConsentimiento: async (pacienteId: string, data: CreateConsentimientoRequest) => {
    const res = await api.post<ApiResponse<Consentimiento>>(
      `/pacientes/${pacienteId}/consentimientos`,
      data
    );
    return res.data;
  },

  getHistoria: async (pacienteId: string) => {
    const res = await api.get<ApiResponse<HistoriaClinica>>(
      `/pacientes/${pacienteId}/historia`
    );
    return res.data;
  },

  updateAntecedentes: async (pacienteId: string, data: UpdateAntecedentesRequest) => {
    const res = await api.put<ApiResponse<HistoriaClinica>>(
      `/pacientes/${pacienteId}/historia/antecedentes`,
      data
    );
    return res.data;
  },

  getNotas: async (pacienteId: string) => {
    const res = await api.get<ApiResponse<NotaEvolucion[]>>(
      `/pacientes/${pacienteId}/historia/notas`
    );
    return res.data;
  },

  createNota: async (pacienteId: string, data: CreateNotaRequest) => {
    const res = await api.post<ApiResponse<NotaEvolucion>>(
      `/pacientes/${pacienteId}/historia/notas`,
      data
    );
    return res.data;
  },
};
