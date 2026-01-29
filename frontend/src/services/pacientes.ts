import api from '../lib/axios';
import type {
  ApiResponse,
  Paciente,
  CreatePacienteRequest,
  Consentimiento,
  CreateConsentimientoRequest,
  HistoriaClinica,
} from '../types';

export const pacientesService = {
  getAll: async (page = 1, perPage = 20) => {
    const res = await api.get<ApiResponse<Paciente[]>>('/pacientes', {
      params: { page, per_page: perPage },
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
};
