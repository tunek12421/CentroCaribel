import api from '../lib/axios';
import type { ApiResponse, Cita, CreateCitaRequest, EstadoCita, TurnoCita } from '../types';

export interface UpdateEstadoData {
  estado: EstadoCita;
  fecha?: string;
  hora?: string;
  turno?: TurnoCita;
}

export const citasService = {
  getAll: async (page = 1, perPage = 20) => {
    const res = await api.get<ApiResponse<Cita[]>>('/citas', {
      params: { page, per_page: perPage },
    });
    return res.data;
  },

  create: async (data: CreateCitaRequest) => {
    const res = await api.post<ApiResponse<Cita>>('/citas', data);
    return res.data;
  },

  updateEstado: async (id: string, data: UpdateEstadoData) => {
    const res = await api.patch<ApiResponse<void>>(`/citas/${id}/estado`, data);
    return res.data;
  },
};
