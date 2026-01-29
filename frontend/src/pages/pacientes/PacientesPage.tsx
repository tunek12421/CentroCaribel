import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { pacientesService } from '../../services/pacientes';
import { Card } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { Modal } from '../../components/ui/Modal';
import { Spinner } from '../../components/ui/Spinner';
import { Plus, Search, Eye } from 'lucide-react';
import { formatDate } from '../../lib/utils';
import { useAuthStore } from '../../store/auth';
import { PacienteForm } from './PacienteForm';

export function PacientesPage() {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const [showForm, setShowForm] = useState(false);
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { hasRole } = useAuthStore();

  const { data, isLoading } = useQuery({
    queryKey: ['pacientes', page],
    queryFn: () => pacientesService.getAll(page),
  });

  const createMutation = useMutation({
    mutationFn: pacientesService.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['pacientes'] });
      setShowForm(false);
    },
  });

  const pacientes = data?.data ?? [];
  const meta = data?.meta;

  const filtered = search
    ? pacientes.filter(
        (p) =>
          p.nombre_completo.toLowerCase().includes(search.toLowerCase()) ||
          p.ci.includes(search) ||
          p.codigo.toLowerCase().includes(search.toLowerCase())
      )
    : pacientes;

  if (isLoading) return <Spinner />;

  return (
    <div className="space-y-4">
      <div className="flex flex-col sm:flex-row gap-3 justify-between">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted" />
          <input
            type="text"
            placeholder="Buscar por nombre, CI o código..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full pl-10"
          />
        </div>
        {hasRole('Administradora', 'Licenciada') && (
          <Button onClick={() => setShowForm(true)}>
            <Plus className="h-4 w-4 mr-2" />
            Nuevo Paciente
          </Button>
        )}
      </div>

      <Card>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border bg-gray-50/50">
                <th className="text-left px-6 py-3 font-medium text-muted">Código</th>
                <th className="text-left px-6 py-3 font-medium text-muted">Nombre</th>
                <th className="text-left px-6 py-3 font-medium text-muted hidden sm:table-cell">CI</th>
                <th className="text-left px-6 py-3 font-medium text-muted hidden md:table-cell">Celular</th>
                <th className="text-left px-6 py-3 font-medium text-muted hidden lg:table-cell">Registro</th>
                <th className="text-right px-6 py-3 font-medium text-muted">Acción</th>
              </tr>
            </thead>
            <tbody>
              {filtered.map((p) => (
                <tr key={p.id} className="border-b border-border hover:bg-gray-50/50">
                  <td className="px-6 py-3 font-mono text-xs">{p.codigo}</td>
                  <td className="px-6 py-3 font-medium">{p.nombre_completo}</td>
                  <td className="px-6 py-3 hidden sm:table-cell">{p.ci}</td>
                  <td className="px-6 py-3 hidden md:table-cell">{p.celular}</td>
                  <td className="px-6 py-3 hidden lg:table-cell text-muted">{formatDate(p.created_at)}</td>
                  <td className="px-6 py-3 text-right">
                    <Button variant="ghost" size="sm" onClick={() => navigate(`/pacientes/${p.id}`)}>
                      <Eye className="h-4 w-4" />
                    </Button>
                  </td>
                </tr>
              ))}
              {filtered.length === 0 && (
                <tr>
                  <td colSpan={6} className="px-6 py-8 text-center text-muted">
                    No se encontraron pacientes
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        {meta && meta.total_pages > 1 && (
          <div className="flex items-center justify-between px-6 py-3 border-t border-border">
            <p className="text-xs text-muted">
              {meta.total} pacientes • Página {meta.page} de {meta.total_pages}
            </p>
            <div className="flex gap-2">
              <Button variant="secondary" size="sm" disabled={page <= 1} onClick={() => setPage(page - 1)}>
                Anterior
              </Button>
              <Button variant="secondary" size="sm" disabled={page >= meta.total_pages} onClick={() => setPage(page + 1)}>
                Siguiente
              </Button>
            </div>
          </div>
        )}
      </Card>

      <Modal open={showForm} onClose={() => setShowForm(false)} title="Nuevo Paciente">
        <PacienteForm
          onSubmit={(data) => createMutation.mutate(data)}
          loading={createMutation.isPending}
          error={createMutation.error}
        />
      </Modal>
    </div>
  );
}
