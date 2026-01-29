import { useQuery } from '@tanstack/react-query';
import { pacientesService } from '../../services/pacientes';
import { citasService } from '../../services/citas';
import { Card, CardContent } from '../../components/ui/Card';
import { Spinner } from '../../components/ui/Spinner';
import { Users, Calendar, ClipboardCheck, Clock } from 'lucide-react';
import { Badge } from '../../components/ui/Badge';
import { ESTADO_COLORS, ESTADO_LABELS, type EstadoCita } from '../../types';

export function DashboardPage() {
  const { data: pacientesData, isLoading: loadingPac } = useQuery({
    queryKey: ['pacientes', 1],
    queryFn: () => pacientesService.getAll(1, 1),
  });

  const { data: citasData, isLoading: loadingCitas } = useQuery({
    queryKey: ['citas', 1],
    queryFn: () => citasService.getAll(1, 100),
  });

  if (loadingPac || loadingCitas) return <Spinner />;

  const totalPacientes = pacientesData?.meta?.total ?? 0;
  const citas = citasData?.data ?? [];
  const hoy = new Date().toISOString().split('T')[0];
  const citasHoy = citas.filter((c) => c.fecha.startsWith(hoy));
  const citasPendientes = citas.filter((c) =>
    ['NUEVA', 'AGENDADA', 'CONFIRMADA'].includes(c.estado)
  );

  const stats = [
    { label: 'Total Pacientes', value: totalPacientes, icon: Users, color: 'text-primary' },
    { label: 'Citas Hoy', value: citasHoy.length, icon: Calendar, color: 'text-accent' },
    { label: 'Citas Pendientes', value: citasPendientes.length, icon: Clock, color: 'text-warning' },
    { label: 'Citas Atendidas', value: citas.filter((c) => c.estado === 'ATENDIDA').length, icon: ClipboardCheck, color: 'text-success' },
  ];

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat) => (
          <Card key={stat.label}>
            <CardContent className="flex items-center gap-4">
              <div className={`p-3 rounded-lg bg-gray-50 ${stat.color}`}>
                <stat.icon className="h-6 w-6" />
              </div>
              <div>
                <p className="text-2xl font-bold">{stat.value}</p>
                <p className="text-sm text-muted">{stat.label}</p>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      <Card>
        <div className="px-6 py-4 border-b border-border">
          <h2 className="text-base font-semibold">Próximas Citas</h2>
        </div>
        <CardContent>
          {citasPendientes.length === 0 ? (
            <p className="text-sm text-muted py-4 text-center">No hay citas pendientes</p>
          ) : (
            <div className="space-y-3">
              {citasPendientes.slice(0, 5).map((cita) => (
                <div key={cita.id} className="flex items-center justify-between py-2 border-b border-border last:border-0">
                  <div>
                    <p className="text-sm font-medium">{cita.tipo_tratamiento}</p>
                    <p className="text-xs text-muted">
                      {new Date(cita.fecha).toLocaleDateString('es-BO')} • {cita.hora} • {cita.turno}
                    </p>
                  </div>
                  <Badge className={ESTADO_COLORS[cita.estado as EstadoCita]}>
                    {ESTADO_LABELS[cita.estado as EstadoCita]}
                  </Badge>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
