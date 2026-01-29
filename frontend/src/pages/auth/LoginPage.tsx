import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../store/auth';
import { authService } from '../../services/auth';
import { Button } from '../../components/ui/Button';

export function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { setTokens } = useAuthStore();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const res = await authService.login({ email, password });
      if (res.success) {
        setTokens(res.data.token, res.data.refresh_token);
        navigate('/');
      } else {
        setError(res.error?.detail || res.error?.message || 'Error al iniciar sesión');
      }
    } catch (err: any) {
      const msg = err.response?.data?.error?.detail || err.response?.data?.error?.message || 'Error de conexión';
      setError(msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-teal-50 to-cyan-50 px-4">
      <div className="w-full max-w-sm">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-primary">Centro Caribel</h1>
          <p className="text-muted text-sm mt-1">Fisioterapia Estética</p>
        </div>

        <div className="bg-white rounded-xl shadow-lg border border-border p-6">
          <h2 className="text-lg font-semibold mb-6">Iniciar Sesión</h2>

          {error && (
            <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-danger">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full"
                placeholder="correo@centrocaribel.com"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Contraseña</label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full"
                placeholder="••••••••"
                required
              />
            </div>

            <Button type="submit" loading={loading} className="w-full">
              Ingresar
            </Button>
          </form>
        </div>

        <p className="text-center text-xs text-muted mt-6">
          Cochabamba, Bolivia
        </p>
      </div>
    </div>
  );
}
