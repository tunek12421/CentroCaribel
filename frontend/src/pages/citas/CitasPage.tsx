import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { citasService, type CitasFilters } from '../../services/citas';
import { Card } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { Modal } from '../../components/ui/Modal';
import { Badge } from '../../components/ui/Badge';
import { Spinner } from '../../components/ui/Spinner';
import { Plus, Filter, X } from 'lucide-react';
import { useAuthStore } from '../../store/auth';
import { ESTADO_COLORS, ESTADO_LABELS, TRANSICIONES_VALIDAS, type EstadoCita, type TurnoCita } from '../../types';
import { CitaForm } from './CitaForm';

function formatHora(hora: string): string {
  if (!hora) return '';
  if (/^\d{2}:\d{2}$/.test(hora)) return hora;
  const match = hora.match(/T(\d{2}:\d{2})/);
  if (match) return match[1];
  const timeParts = hora.match(/^(\d{2}:\d{2}):\d{2}$/);
  if (timeParts) return timeParts[1];
  return hora;
}

export function CitasPage() {
  const [page, setPage] = useState(1);
  const [showForm, setShowForm] = useState(false);
  const [estadoModal, setEstadoModal] = useState<{ id: string; estado: EstadoCita } | null>(null);
  const [reagendarForm, setReagendarForm] = useState<{ id: string } | null>(null);
  const [reagendarData, setReagendarData] = useState({ fecha: '', hora: '', turno: 'AM' as TurnoCita });
  const [filters, setFilters] = useState<CitasFilters>({});
  const [showFilters, setShowFilters] = useState(false);
  const queryClient = useQueryClient();
  const { hasRole } = useAuthStore();

  const hasActiveFilters = !!(filters.fecha || filters.turno || filters.estado);

  const { data, isLoading } = useQuery({
    queryKey: ['citas', page, filters],
    queryFn: () => citasService.getAll(page, 20, filters),
    refetchInterval: 30000,
  });

  const createMutation = useMutation({
    mutationFn: citasService.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['citas'] });
      setShowForm(false);
    },
  });

  const estadoMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Parameters<typeof citasService.updateEstado>[1] }) =>
      citasService.updateEstado(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['citas'] });
      setEstadoModal(null);
      setReagendarForm(null);
      setReagendarData({ fecha: '', hora: '', turno: 'AM' });
    },
  });

  if (isLoading) return <Spinner />;

  const citas = data?.data ?? [];
  const meta = data?.meta;

  const handleEstadoClick = (nuevoEstado: EstadoCita, citaId: string) => {
    if (nuevoEstado === 'REAGENDADA') {
      setEstadoModal(null);
      setReagendarForm({ id: citaId });
    } else {
      estadoMutation.mutate({ id: citaId, data: { estado: nuevoEstado } });
    }
  };

  const handleReagendar = () => {
    if (!reagendarForm || !reagendarData.fecha || !reagendarData.hora) return;
    estadoMutation.mutate({
      id: reagendarForm.id,
      data: {
        estado: 'REAGENDADA',
        fecha: reagendarData.fecha,
        hora: reagendarData.hora,
        turno: reagendarData.turno,
      },
    });
  };

  const clearFilters = () => {
    setFilters({});
    setPage(1);
  };

  return (
    <div className="space-y-4">
      {/* Info de turnos */}
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-3 text-sm text-blue-800">
        <strong>Turno AM:</strong> Propietaria &nbsp;|&nbsp; <strong>Turno PM:</strong> Lic. Tatiana Mart&iacute;nez
      </div>

      {/* Header con acciones */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-2">
        <div className="flex items-center gap-2">
          <p className="text-sm text-muted">{meta?.total ?? 0} citas</p>
          <Button
            variant="secondary"
            size="sm"
            onClick={() => setShowFilters(!showFilters)}
          >
            <Filter className="h-3.5 w-3.5 mr-1" />
            Filtros
            {hasActiveFilters && <span className="ml-1 bg-primary text-white rounded-full w-4 h-4 text-xs flex items-center justify-center">!</span>}
          </Button>
          {hasActiveFilters && (
            <Button variant="secondary" size="sm" onClick={clearFilters}>
              <X className="h-3.5 w-3.5 mr-1" />
              Limpiar
            </Button>
          )}
        </div>
        {hasRole('Administradora', 'Licenciada') && (
          <Button onClick={() => setShowForm(true)}>
            <Plus className="h-4 w-4 mr-2" />
            Nueva Cita
          </Button>
        )}
      </div>

      {/* Barra de filtros */}
      {showFilters && (
        <Card className="p-4">
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div>
              <label className="block text-xs font-medium text-muted mb-1">Fecha</label>
              <input
                type="date"
                className="w-full"
                value={filters.fecha ?? ''}
                onChange={(e) => { setFilters({ ...filters, fecha: e.target.value || undefined }); setPage(1); }}
              />
            </div>
            <div>
              <label className="block text-xs font-medium text-muted mb-1">Turno</label>
              <select
                className="w-full"
                value={filters.turno ?? ''}
                onChange={(e) => { setFilters({ ...filters, turno: (e.target.value || undefined) as TurnoCita | undefined }); setPage(1); }}
              >
                <option value="">Todos</option>
                <option value="AM">AM (Mañana)</option>
                <option value="PM">PM (Tarde)</option>
              </select>
            </div>
            <div>
              <label className="block text-xs font-medium text-muted mb-1">Estado</label>
              <select
                className="w-full"
                value={filters.estado ?? ''}
                onChange={(e) => { setFilters({ ...filters, estado: (e.target.value || undefined) as EstadoCita | undefined }); setPage(1); }}
              >
                <option value="">Todos</option>
                {Object.entries(ESTADO_LABELS).map(([key, label]) => (
                  <option key={key} value={key}>{label}</option>
                ))}
              </select>
            </div>
          </div>
        </Card>
      )}

      {/* Tabla de citas */}
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
                    <td className="px-6 py-3">{formatHora(c.hora)}</td>
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

      {/* Modal nueva cita */}
      <Modal open={showForm} onClose={() => setShowForm(false)} title="Nueva Cita">
        <CitaForm
          onSubmit={(data) => createMutation.mutate(data)}
          loading={createMutation.isPending}
          error={createMutation.error}
        />
      </Modal>

      {/* Modal cambiar estado */}
      <Modal open={!!estadoModal} onClose={() => setEstadoModal(null)} title="Cambiar Estado de Cita">
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
                  onClick={() => handleEstadoClick(nuevoEstado, estadoModal.id)}
                >
                  {ESTADO_LABELS[nuevoEstado]}
                </Button>
              ))}
            </div>
          </div>
        )}
      </Modal>

      {/* Modal reagendar */}
      <Modal
        open={!!reagendarForm}
        onClose={() => {
          setReagendarForm(null);
          setReagendarData({ fecha: '', hora: '', turno: 'AM' });
        }}
        title="Reagendar Cita"
      >
        {reagendarForm && (
          <div className="space-y-4">
            <p className="text-sm text-muted">
              Ingrese la nueva fecha y hora para la cita:
            </p>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Nueva Fecha *</label>
                <input
                  type="date"
                  className="w-full"
                  value={reagendarData.fecha}
                  onChange={(e) => setReagendarData({ ...reagendarData, fecha: e.target.value })}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Nueva Hora *</label>
                <input
                  type="time"
                  className="w-full"
                  value={reagendarData.hora}
                  onChange={(e) => setReagendarData({ ...reagendarData, hora: e.target.value })}
                />
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Turno *</label>
              <select
                className="w-full"
                value={reagendarData.turno}
                onChange={(e) => setReagendarData({ ...reagendarData, turno: e.target.value as TurnoCita })}
              >
                <option value="AM">AM (Mañana)</option>
                <option value="PM">PM (Tarde)</option>
              </select>
            </div>
            <Button
              onClick={handleReagendar}
              loading={estadoMutation.isPending}
              disabled={!reagendarData.fecha || !reagendarData.hora}
              className="w-full"
            >
              Confirmar Reagendamiento
            </Button>
          </div>
        )}
      </Modal>
    </div>
  );
}
