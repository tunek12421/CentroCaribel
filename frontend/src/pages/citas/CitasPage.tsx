import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { citasService } from '../../services/citas';
import { Card } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { Modal } from '../../components/ui/Modal';
import { Badge } from '../../components/ui/Badge';
import { Spinner } from '../../components/ui/Spinner';
import { Plus } from 'lucide-react';
import { useAuthStore } from '../../store/auth';
import { ESTADO_COLORS, ESTADO_LABELS, TRANSICIONES_VALIDAS, type EstadoCita } from '../../types';
import { CitaForm } from './CitaForm';

export function CitasPage() {
  const [page, setPage] = useState(1);
  const [showForm, setShowForm] = useState(false);
  const [estadoModal, setEstadoModal] = useState<{ id: string; estado: EstadoCita } | null>(null);
  const queryClient = useQueryClient();
  const { hasRole } = useAuthStore();

  const { data, isLoading } = useQuery({
    queryKey: ['citas', page],
    queryFn: () => citasService.getAll(page),
  });

  const createMutation = useMutation({
    mutationFn: citasService.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['citas'] });
      setShowForm(false);
    },
  });

  const estadoMutation = useMutation({
    mutationFn: ({ id, estado }: { id: string; estado: EstadoCita }) =>
      citasService.updateEstado(id, estado),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['citas'] });
      setEstadoModal(null);
    },
  });

  if (isLoading) return <Spinner />;

  const citas = data?.data ?? [];
  const meta = data?.meta;

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <p className="text-sm text-muted">{meta?.total ?? 0} citas</p>
        {hasRole('Administradora', 'Licenciada') && (
          <Button onClick={() => setShowForm(true)}>
            <Plus className="h-4 w-4 mr-2" />
            Nueva Cita
          </Button>
        )}
      </div>

      <Card>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border bg-gray-50/50">
                <th className="text-left px-6 py-3 font-medium text-muted">Fecha</th>
                <th className="text-left px-6 py-3 font-medium text-muted">Hora</th>
                <th className="text-left px-6 py-3 font-medium text-muted hidden sm:table-cell">Turno</th>
                <th className="text-left px-6 py-3 font-medium text-muted">Tratamiento</th>
                <th className="text-left px-6 py-3 font-medium text-muted">Estado</th>
                <th className="text-right px-6 py-3 font-medium text-muted">Acción</th>
              </tr>
            </thead>
            <tbody>
              {citas.map((c) => {
                const transiciones = TRANSICIONES_VALIDAS[c.estado] ?? [];
                return (
                  <tr key={c.id} className="border-b border-border hover:bg-gray-50/50">
                    <td className="px-6 py-3">{new Date(c.fecha).toLocaleDateString('es-BO')}</td>
                    <td className="px-6 py-3">{c.hora}</td>
                    <td className="px-6 py-3 hidden sm:table-cell">
                      <Badge className={c.turno === 'AM' ? 'bg-sky-100 text-sky-800' : 'bg-orange-100 text-orange-800'}>
                        {c.turno}
                      </Badge>
                    </td>
                    <td className="px-6 py-3">{c.tipo_tratamiento}</td>
                    <td className="px-6 py-3">
                      <Badge className={ESTADO_COLORS[c.estado]}>{ESTADO_LABELS[c.estado]}</Badge>
                    </td>
                    <td className="px-6 py-3 text-right">
                      {transiciones.length > 0 && hasRole('Administradora', 'Licenciada') && (
                        <Button
                          variant="secondary"
                          size="sm"
                          onClick={() => setEstadoModal({ id: c.id, estado: c.estado })}
                        >
                          Cambiar
                        </Button>
                      )}
                    </td>
                  </tr>
                );
              })}
              {citas.length === 0 && (
                <tr>
                  <td colSpan={6} className="px-6 py-8 text-center text-muted">
                    No hay citas registradas
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        {meta && meta.total_pages > 1 && (
          <div className="flex items-center justify-between px-6 py-3 border-t border-border">
            <p className="text-xs text-muted">Página {meta.page} de {meta.total_pages}</p>
            <div className="flex gap-2">
              <Button variant="secondary" size="sm" disabled={page <= 1} onClick={() => setPage(page - 1)}>Anterior</Button>
              <Button variant="secondary" size="sm" disabled={page >= meta.total_pages} onClick={() => setPage(page + 1)}>Siguiente</Button>
            </div>
          </div>
        )}
      </Card>

      <Modal open={showForm} onClose={() => setShowForm(false)} title="Nueva Cita">
        <CitaForm
          onSubmit={(data) => createMutation.mutate(data)}
          loading={createMutation.isPending}
          error={createMutation.error}
        />
      </Modal>

      <Modal
        open={!!estadoModal}
        onClose={() => setEstadoModal(null)}
        title="Cambiar Estado de Cita"
      >
        {estadoModal && (
          <div className="space-y-3">
            <p className="text-sm text-muted">
              Estado actual: <Badge className={ESTADO_COLORS[estadoModal.estado]}>{ESTADO_LABELS[estadoModal.estado]}</Badge>
            </p>
            <p className="text-sm font-medium">Seleccione el nuevo estado:</p>
            <div className="flex flex-wrap gap-2">
              {TRANSICIONES_VALIDAS[estadoModal.estado].map((nuevoEstado) => (
                <Button
                  key={nuevoEstado}
                  variant="secondary"
                  size="sm"
                  loading={estadoMutation.isPending}
                  onClick={() => estadoMutation.mutate({ id: estadoModal.id, estado: nuevoEstado })}
                >
                  {ESTADO_LABELS[nuevoEstado]}
                </Button>
              ))}
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
}
