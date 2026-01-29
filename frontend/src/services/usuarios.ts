import api from '../lib/axios';
import type { ApiResponse, Usuario, CreateUsuarioRequest, UpdateUsuarioRequest, Rol } from '../types';

export const usuariosService = {
  getAll: async (page = 1, perPage = 20) => {
    const res = await api.get<ApiResponse<Usuario[]>>('/usuarios', {
      params: { page, per_page: perPage },
    });
    return res.data;
  },

  getById: async (id: string) => {
    const res = await api.get<ApiResponse<Usuario>>(`/usuarios/${id}`);
    return res.data;
  },

  create: async (data: CreateUsuarioRequest) => {
    const res = await api.post<ApiResponse<Usuario>>('/usuarios', data);
    return res.data;
  },

  update: async (id: string, data: UpdateUsuarioRequest) => {
    const res = await api.put<ApiResponse<Usuario>>(`/usuarios/${id}`, data);
    return res.data;
  },

  delete: async (id: string) => {
    const res = await api.delete<ApiResponse<void>>(`/usuarios/${id}`);
    return res.data;
  },

  getRoles: async () => {
    const res = await api.get<ApiResponse<Rol[]>>('/roles');
    return res.data.data ?? [];
  },
};
