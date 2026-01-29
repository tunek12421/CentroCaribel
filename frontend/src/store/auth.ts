import { create } from 'zustand';
import { jwtDecode } from '../lib/jwt';

interface AuthState {
  token: string | null;
  refreshToken: string | null;
  userId: string | null;
  rolNombre: string | null;
  isAuthenticated: boolean;
  setTokens: (token: string, refreshToken: string) => void;
  logout: () => void;
  hasRole: (...roles: string[]) => boolean;
}

export const useAuthStore = create<AuthState>((set, get) => {
  const savedToken = localStorage.getItem('token');
  const savedRefresh = localStorage.getItem('refresh_token');
  let decoded: { user_id: string; rol_nombre: string } | null = null;

  if (savedToken) {
    decoded = jwtDecode(savedToken);
    if (!decoded) {
      localStorage.removeItem('token');
      localStorage.removeItem('refresh_token');
    }
  }

  return {
    token: decoded ? savedToken : null,
    refreshToken: decoded ? savedRefresh : null,
    userId: decoded?.user_id ?? null,
    rolNombre: decoded?.rol_nombre ?? null,
    isAuthenticated: !!decoded,

    setTokens: (token, refreshToken) => {
      localStorage.setItem('token', token);
      localStorage.setItem('refresh_token', refreshToken);
      const claims = jwtDecode(token);
      set({
        token,
        refreshToken,
        userId: claims?.user_id ?? null,
        rolNombre: claims?.rol_nombre ?? null,
        isAuthenticated: true,
      });
    },

    logout: () => {
      localStorage.removeItem('token');
      localStorage.removeItem('refresh_token');
      set({
        token: null,
        refreshToken: null,
        userId: null,
        rolNombre: null,
        isAuthenticated: false,
      });
    },

    hasRole: (...roles) => {
      const { rolNombre } = get();
      return rolNombre !== null && roles.includes(rolNombre);
    },
  };
});
