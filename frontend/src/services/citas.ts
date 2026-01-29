import api from '../lib/axios';
import type { ApiResponse, Cita, CreateCitaRequest, EstadoCita } from '../types';

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

  updateEstado: async (id: string, estado: EstadoCita) => {
    const res = await api.patch<ApiResponse<void>>(`/citas/${id}/estado`, { estado });
    return res.data;
  },
};
