import { useState } from 'react';
import { useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { pacientesService } from '../../services/pacientes';
import { Card, CardContent, CardHeader } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { Modal } from '../../components/ui/Modal';
import { Badge } from '../../components/ui/Badge';
import { Spinner } from '../../components/ui/Spinner';
import { formatDate, formatDateTime } from '../../lib/utils';
import { useAuthStore } from '../../store/auth';
import { ConsentimientoForm } from './ConsentimientoForm';
import { FileText, Phone, MapPin, Calendar, Hash, User } from 'lucide-react';

export function PacienteDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [showConsentimiento, setShowConsentimiento] = useState(false);
  const queryClient = useQueryClient();
  const { hasRole } = useAuthStore();

  const { data: pacienteData, isLoading: loadingPac } = useQuery({
    queryKey: ['paciente', id],
    queryFn: () => pacientesService.getById(id!),
    enabled: !!id,
  });

  const { data: historiaData } = useQuery({
    queryKey: ['historia', id],
    queryFn: () => pacientesService.getHistoria(id!),
    enabled: !!id,
  });

  const { data: consData } = useQuery({
    queryKey: ['consentimientos', id],
    queryFn: () => pacientesService.getConsentimientos(id!),
    enabled: !!id,
  });

  const consMutation = useMutation({
    mutationFn: (data: { firma_digital?: string; autoriza_fotos: boolean; contenido: string }) =>
      pacientesService.createConsentimiento(id!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['consentimientos', id] });
      setShowConsentimiento(false);
    },
  });

  if (loadingPac) return <Spinner />;

  const paciente = pacienteData?.data;
  const historia = historiaData?.data;
  const consentimientos = consData?.data ?? [];

  if (!paciente) return <p className="text-center text-muted py-8">Paciente no encontrado</p>;

  return (
    <div className="space-y-6 max-w-4xl">
      {/* Info del paciente */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-xl font-bold">{paciente.nombre_completo}</h2>
              <p className="text-sm text-muted font-mono">{paciente.codigo}</p>
            </div>
            {historia && (
              <Badge className="bg-green-100 text-green-800">
                {historia.numero_historia} • {historia.estado}
              </Badge>
            )}
          </div>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
            <div className="flex items-center gap-2 text-gray-600">
              <User className="h-4 w-4" />
              <span>CI: {paciente.ci}</span>
            </div>
            <div className="flex items-center gap-2 text-gray-600">
              <Calendar className="h-4 w-4" />
              <span>Nac: {formatDate(paciente.fecha_nacimiento)}</span>
            </div>
            <div className="flex items-center gap-2 text-gray-600">
              <Phone className="h-4 w-4" />
              <span>{paciente.celular}</span>
            </div>
            {paciente.direccion && (
              <div className="flex items-center gap-2 text-gray-600">
                <MapPin className="h-4 w-4" />
                <span>{paciente.direccion}</span>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Historia Clínica */}
      {historia && (
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <Hash className="h-5 w-5 text-primary" />
              <h3 className="font-semibold">Historia Clínica</h3>
            </div>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <p className="text-muted">Número</p>
                <p className="font-medium font-mono">{historia.numero_historia}</p>
              </div>
              <div>
                <p className="text-muted">Estado</p>
                <Badge className="bg-green-100 text-green-800">{historia.estado}</Badge>
              </div>
              <div>
                <p className="text-muted">Creada</p>
                <p className="font-medium">{formatDateTime(historia.created_at)}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Consentimientos */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <FileText className="h-5 w-5 text-primary" />
              <h3 className="font-semibold">Consentimientos Informados</h3>
            </div>
            {hasRole('Administradora', 'Licenciada') && (
              <Button size="sm" onClick={() => setShowConsentimiento(true)}>
                Nuevo
              </Button>
            )}
          </div>
        </CardHeader>
        <CardContent>
          {consentimientos.length === 0 ? (
            <p className="text-sm text-muted text-center py-4">Sin consentimientos registrados</p>
          ) : (
            <div className="space-y-3">
              {consentimientos.map((c) => (
                <div key={c.id} className="p-3 border border-border rounded-lg">
                  <div className="flex items-center justify-between mb-2">
                    <p className="text-xs text-muted">{formatDateTime(c.fecha_firma)}</p>
                    <Badge className={c.autoriza_fotos ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'}>
                      {c.autoriza_fotos ? 'Autoriza fotos' : 'No autoriza fotos'}
                    </Badge>
                  </div>
                  <p className="text-sm">{c.contenido}</p>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      <Modal open={showConsentimiento} onClose={() => setShowConsentimiento(false)} title="Nuevo Consentimiento">
        <ConsentimientoForm
          onSubmit={(data) => consMutation.mutate(data)}
          loading={consMutation.isPending}
          error={consMutation.error}
        />
      </Modal>
    </div>
  );
}
