import api from '../lib/axios';
import type { ApiResponse, LoginRequest, LoginResponse } from '../types';

export const authService = {
  login: async (data: LoginRequest) => {
    const res = await api.post<ApiResponse<LoginResponse>>('/auth/login', data);
    return res.data;
  },

  refresh: async (refreshToken: string) => {
    const res = await api.post<ApiResponse<LoginResponse>>('/auth/refresh', {
      refresh_token: refreshToken,
    });
    return res.data;
  },
};
