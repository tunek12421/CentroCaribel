import { useState } from 'react';
import { useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { pacientesService } from '../../services/pacientes';
import { paquetesService } from '../../services/paquetes';
import { Card, CardContent, CardHeader } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { Modal } from '../../components/ui/Modal';
import { Badge } from '../../components/ui/Badge';
import { Spinner } from '../../components/ui/Spinner';
import { formatDate, formatDateTime } from '../../lib/utils';
import { useAuthStore } from '../../store/auth';
import { ConsentimientoForm } from './ConsentimientoForm';
import { FileText, Phone, MapPin, Calendar, Hash, User, Package, ClipboardList, Plus, Edit3, Stethoscope, Activity, StickyNote } from 'lucide-react';
import type { CreatePaqueteRequest, NotaEvolucion } from '../../types';

const ESTADO_PAQUETE_COLORS: Record<string, string> = {
  ACTIVO: 'bg-green-100 text-green-800',
  COMPLETADO: 'bg-blue-100 text-blue-800',
  CANCELADO: 'bg-gray-100 text-gray-800',
};

const NOTA_TIPO_LABELS: Record<string, string> = {
  TRATAMIENTO: 'Tratamiento',
  EVOLUCION: 'Evolución',
  NOTA: 'Nota',
};

const NOTA_TIPO_COLORS: Record<string, string> = {
  TRATAMIENTO: 'bg-blue-100 text-blue-800',
  EVOLUCION: 'bg-green-100 text-green-800',
  NOTA: 'bg-gray-100 text-gray-800',
};

const NOTA_TIPO_ICONS: Record<string, typeof Stethoscope> = {
  TRATAMIENTO: Stethoscope,
  EVOLUCION: Activity,
  NOTA: StickyNote,
};

export function PacienteDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [showConsentimiento, setShowConsentimiento] = useState(false);
  const [showPaqueteForm, setShowPaqueteForm] = useState(false);
  const [showAntecedentesForm, setShowAntecedentesForm] = useState(false);
  const [showNotaForm, setShowNotaForm] = useState(false);
  const [paqueteForm, setPaqueteForm] = useState({ tipo_tratamiento: '', total_sesiones: 1, notas: '' });
  const [notaForm, setNotaForm] = useState({ tipo: 'EVOLUCION' as 'TRATAMIENTO' | 'EVOLUCION' | 'NOTA', contenido: '' });
  const [antForm, setAntForm] = useState({
    antecedentes_personales: '',
    antecedentes_familiares: '',
    alergias: '',
    medicamentos_actuales: '',
  });
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

  const { data: notasData } = useQuery({
    queryKey: ['notas', id],
    queryFn: () => pacientesService.getNotas(id!),
    enabled: !!id,
  });

  const { data: consData } = useQuery({
    queryKey: ['consentimientos', id],
    queryFn: () => pacientesService.getConsentimientos(id!),
    enabled: !!id,
  });

  const { data: paquetesData } = useQuery({
    queryKey: ['paquetes', id],
    queryFn: () => paquetesService.getByPaciente(id!),
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

  const paqueteMutation = useMutation({
    mutationFn: (data: CreatePaqueteRequest) => paquetesService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['paquetes', id] });
      setShowPaqueteForm(false);
      setPaqueteForm({ tipo_tratamiento: '', total_sesiones: 1, notas: '' });
    },
  });

  const antecedentesMutation = useMutation({
    mutationFn: (data: typeof antForm) => pacientesService.updateAntecedentes(id!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['historia', id] });
      setShowAntecedentesForm(false);
    },
  });

  const notaMutation = useMutation({
    mutationFn: (data: typeof notaForm) => pacientesService.createNota(id!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notas', id] });
      setShowNotaForm(false);
      setNotaForm({ tipo: 'EVOLUCION', contenido: '' });
    },
  });

  if (loadingPac) return <Spinner />;

  const paciente = pacienteData?.data;
  const historia = historiaData?.data;
  const notas: NotaEvolucion[] = notasData?.data ?? [];
  const consentimientos = consData?.data ?? [];
  const paquetes = paquetesData?.data ?? [];

  if (!paciente) return <p className="text-center text-muted py-8">Paciente no encontrado</p>;

  const handlePaqueteSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    paqueteMutation.mutate({
      paciente_id: id!,
      tipo_tratamiento: paqueteForm.tipo_tratamiento,
      total_sesiones: paqueteForm.total_sesiones,
      notas: paqueteForm.notas || undefined,
    });
  };

  const handleOpenAntecedentes = () => {
    if (historia) {
      setAntForm({
        antecedentes_personales: historia.antecedentes_personales || '',
        antecedentes_familiares: historia.antecedentes_familiares || '',
        alergias: historia.alergias || '',
        medicamentos_actuales: historia.medicamentos_actuales || '',
      });
    }
    setShowAntecedentesForm(true);
  };

  const handleAntecedentesSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    antecedentesMutation.mutate(antForm);
  };

  const handleNotaSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    notaMutation.mutate(notaForm);
  };

  const paqueteApiError = (paqueteMutation.error as any)?.response?.data?.error;
  const antApiError = (antecedentesMutation.error as any)?.response?.data?.error;
  const notaApiError = (notaMutation.error as any)?.response?.data?.error;

  const hasAntecedentes = historia && (
    historia.antecedentes_personales || historia.antecedentes_familiares ||
    historia.alergias || historia.medicamentos_actuales
  );

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

      {/* Historia Clínica - Antecedentes */}
      {historia && (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <ClipboardList className="h-5 w-5 text-primary" />
                <h3 className="font-semibold">Historia Clínica — {historia.numero_historia}</h3>
              </div>
              {hasRole('Administradora', 'Licenciada') && (
                <Button size="sm" variant="secondary" onClick={handleOpenAntecedentes}>
                  <Edit3 className="h-3.5 w-3.5 mr-1" />
                  {hasAntecedentes ? 'Editar' : 'Registrar'}
                </Button>
              )}
            </div>
          </CardHeader>
          <CardContent>
            {!hasAntecedentes ? (
              <p className="text-sm text-muted text-center py-4">Sin antecedentes registrados</p>
            ) : (
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
                {historia.antecedentes_personales && (
                  <div>
                    <p className="text-muted font-medium mb-1">Antecedentes Personales</p>
                    <p className="whitespace-pre-wrap">{historia.antecedentes_personales}</p>
                  </div>
                )}
                {historia.antecedentes_familiares && (
                  <div>
                    <p className="text-muted font-medium mb-1">Antecedentes Familiares</p>
                    <p className="whitespace-pre-wrap">{historia.antecedentes_familiares}</p>
                  </div>
                )}
                {historia.alergias && (
                  <div>
                    <p className="text-muted font-medium mb-1">Alergias</p>
                    <p className="whitespace-pre-wrap">{historia.alergias}</p>
                  </div>
                )}
                {historia.medicamentos_actuales && (
                  <div>
                    <p className="text-muted font-medium mb-1">Medicamentos Actuales</p>
                    <p className="whitespace-pre-wrap">{historia.medicamentos_actuales}</p>
                  </div>
                )}
              </div>
            )}
            <div className="mt-3 pt-3 border-t border-border flex items-center gap-4 text-xs text-muted">
              <span>Creada: {formatDateTime(historia.created_at)}</span>
              <span>Actualizada: {formatDateTime(historia.updated_at)}</span>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Notas de Evolución */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Hash className="h-5 w-5 text-primary" />
              <h3 className="font-semibold">Notas de Evolución</h3>
            </div>
            {hasRole('Administradora', 'Licenciada', 'Medico') && historia && (
              <Button size="sm" onClick={() => setShowNotaForm(true)}>
                <Plus className="h-3.5 w-3.5 mr-1" />
                Nueva Nota
              </Button>
            )}
          </div>
        </CardHeader>
        <CardContent>
          {notas.length === 0 ? (
            <p className="text-sm text-muted text-center py-4">Sin notas de evolución registradas</p>
          ) : (
            <div className="space-y-3">
              {notas.map((n) => {
                const Icon = NOTA_TIPO_ICONS[n.tipo] ?? StickyNote;
                return (
                  <div key={n.id} className="p-3 border border-border rounded-lg">
                    <div className="flex items-center justify-between mb-2">
                      <div className="flex items-center gap-2">
                        <Icon className="h-4 w-4 text-muted" />
                        <Badge className={NOTA_TIPO_COLORS[n.tipo] ?? 'bg-gray-100 text-gray-800'}>
                          {NOTA_TIPO_LABELS[n.tipo] ?? n.tipo}
                        </Badge>
                      </div>
                      <p className="text-xs text-muted">{formatDateTime(n.created_at)}</p>
                    </div>
                    <p className="text-sm whitespace-pre-wrap">{n.contenido}</p>
                  </div>
                );
              })}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Paquetes de Tratamiento */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Package className="h-5 w-5 text-primary" />
              <h3 className="font-semibold">Paquetes de Tratamiento</h3>
            </div>
            {hasRole('Administradora', 'Licenciada') && (
              <Button size="sm" onClick={() => setShowPaqueteForm(true)}>
                Nuevo
              </Button>
            )}
          </div>
        </CardHeader>
        <CardContent>
          {paquetes.length === 0 ? (
            <p className="text-sm text-muted text-center py-4">Sin paquetes de tratamiento</p>
          ) : (
            <div className="space-y-3">
              {paquetes.map((p) => (
                <div key={p.id} className="p-3 border border-border rounded-lg">
                  <div className="flex items-center justify-between mb-2">
                    <p className="text-sm font-medium">{p.tipo_tratamiento}</p>
                    <Badge className={ESTADO_PAQUETE_COLORS[p.estado] ?? 'bg-gray-100 text-gray-800'}>
                      {p.estado}
                    </Badge>
                  </div>
                  <div className="flex items-center gap-4 text-xs text-muted">
                    <span>Sesiones: {p.sesiones_completadas}/{p.total_sesiones}</span>
                    <div className="flex-1 bg-gray-200 rounded-full h-2">
                      <div
                        className="bg-primary rounded-full h-2 transition-all"
                        style={{ width: `${Math.min(100, (p.sesiones_completadas / p.total_sesiones) * 100)}%` }}
                      />
                    </div>
                  </div>
                  {p.notas && <p className="text-xs text-muted mt-1">{p.notas}</p>}
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

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
                  {c.firma_digital && (
                    <div className="mt-3 pt-3 border-t border-border">
                      <p className="text-xs text-muted mb-1">Firma del paciente:</p>
                      <img
                        src={`data:image/png;base64,${c.firma_digital}`}
                        alt="Firma del paciente"
                        className="max-h-24 border border-border rounded bg-white"
                      />
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Modal consentimiento */}
      <Modal open={showConsentimiento} onClose={() => setShowConsentimiento(false)} title="Nuevo Consentimiento">
        <ConsentimientoForm
          onSubmit={(data) => consMutation.mutate(data)}
          loading={consMutation.isPending}
          error={consMutation.error}
        />
      </Modal>

      {/* Modal paquete */}
      <Modal open={showPaqueteForm} onClose={() => { setShowPaqueteForm(false); setPaqueteForm({ tipo_tratamiento: '', total_sesiones: 1, notas: '' }); }} title="Nuevo Paquete de Tratamiento">
        <form onSubmit={handlePaqueteSubmit} className="space-y-4">
          {paqueteApiError && (
            <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-danger">
              {paqueteApiError.detail || paqueteApiError.message}
            </div>
          )}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Tipo de Tratamiento *</label>
            <input
              className="w-full"
              placeholder="Masaje terapéutico"
              value={paqueteForm.tipo_tratamiento}
              onChange={(e) => setPaqueteForm({ ...paqueteForm, tipo_tratamiento: e.target.value })}
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Total de Sesiones *</label>
            <input
              type="number"
              min={1}
              className="w-full"
              value={paqueteForm.total_sesiones}
              onChange={(e) => setPaqueteForm({ ...paqueteForm, total_sesiones: parseInt(e.target.value) || 1 })}
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Notas</label>
            <textarea
              className="w-full"
              rows={2}
              placeholder="Indicaciones adicionales..."
              value={paqueteForm.notas}
              onChange={(e) => setPaqueteForm({ ...paqueteForm, notas: e.target.value })}
            />
          </div>
          <Button type="submit" loading={paqueteMutation.isPending} className="w-full">
            Crear Paquete
          </Button>
        </form>
      </Modal>

      {/* Modal antecedentes */}
      <Modal open={showAntecedentesForm} onClose={() => setShowAntecedentesForm(false)} title="Antecedentes del Paciente">
        <form onSubmit={handleAntecedentesSubmit} className="space-y-4">
          {antApiError && (
            <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-danger">
              {antApiError.detail || antApiError.message}
            </div>
          )}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Antecedentes Personales</label>
            <textarea
              className="w-full"
              rows={3}
              placeholder="Enfermedades previas, cirugías, condiciones crónicas..."
              value={antForm.antecedentes_personales}
              onChange={(e) => setAntForm({ ...antForm, antecedentes_personales: e.target.value })}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Antecedentes Familiares</label>
            <textarea
              className="w-full"
              rows={3}
              placeholder="Enfermedades hereditarias, condiciones familiares..."
              value={antForm.antecedentes_familiares}
              onChange={(e) => setAntForm({ ...antForm, antecedentes_familiares: e.target.value })}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Alergias</label>
            <textarea
              className="w-full"
              rows={2}
              placeholder="Alergias conocidas a medicamentos, alimentos, etc."
              value={antForm.alergias}
              onChange={(e) => setAntForm({ ...antForm, alergias: e.target.value })}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Medicamentos Actuales</label>
            <textarea
              className="w-full"
              rows={2}
              placeholder="Medicamentos que está tomando actualmente..."
              value={antForm.medicamentos_actuales}
              onChange={(e) => setAntForm({ ...antForm, medicamentos_actuales: e.target.value })}
            />
          </div>
          <Button type="submit" loading={antecedentesMutation.isPending} className="w-full">
            Guardar Antecedentes
          </Button>
        </form>
      </Modal>

      {/* Modal nueva nota */}
      <Modal open={showNotaForm} onClose={() => { setShowNotaForm(false); setNotaForm({ tipo: 'EVOLUCION', contenido: '' }); }} title="Nueva Nota de Evolución">
        <form onSubmit={handleNotaSubmit} className="space-y-4">
          {notaApiError && (
            <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-danger">
              {notaApiError.detail || notaApiError.message}
            </div>
          )}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Tipo *</label>
            <select
              className="w-full"
              value={notaForm.tipo}
              onChange={(e) => setNotaForm({ ...notaForm, tipo: e.target.value as 'TRATAMIENTO' | 'EVOLUCION' | 'NOTA' })}
            >
              <option value="EVOLUCION">Evolución</option>
              <option value="TRATAMIENTO">Tratamiento</option>
              <option value="NOTA">Nota General</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Contenido *</label>
            <textarea
              className="w-full"
              rows={4}
              placeholder="Descripción detallada..."
              value={notaForm.contenido}
              onChange={(e) => setNotaForm({ ...notaForm, contenido: e.target.value })}
              required
            />
          </div>
          <Button type="submit" loading={notaMutation.isPending} className="w-full">
            Guardar Nota
          </Button>
        </form>
      </Modal>
    </div>
  );
}
