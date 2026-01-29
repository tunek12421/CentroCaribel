import api from '../lib/axios';
import type { ApiResponse, PaqueteTratamiento, CreatePaqueteRequest } from '../types';

export const paquetesService = {
  getByPaciente: async (pacienteId: string, onlyActivos = false) => {
    const res = await api.get<ApiResponse<PaqueteTratamiento[]>>(
      `/pacientes/${pacienteId}/paquetes`,
      { params: onlyActivos ? { activos: 'true' } : {} }
    );
    return res.data;
  },

  create: async (data: CreatePaqueteRequest) => {
    const res = await api.post<ApiResponse<PaqueteTratamiento>>('/paquetes', data);
    return res.data;
  },
};
