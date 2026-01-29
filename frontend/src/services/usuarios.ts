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
    // Roles come from the usuarios list - we extract unique roles
    // For now we use a hardcoded approach since there's no /roles endpoint
    const res = await api.get<ApiResponse<Usuario[]>>('/usuarios', { params: { per_page: 100 } });
    const roles: Rol[] = [];
    const seen = new Set<string>();
    for (const u of res.data.data || []) {
      if (u.rol && !seen.has(u.rol.id)) {
        seen.add(u.rol.id);
        roles.push(u.rol);
      }
    }
    return roles;
  },
};
