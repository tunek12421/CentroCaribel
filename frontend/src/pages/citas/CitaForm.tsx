import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useQuery } from '@tanstack/react-query';
import { pacientesService } from '../../services/pacientes';
import { paquetesService } from '../../services/paquetes';
import { Button } from '../../components/ui/Button';
import type { CreateCitaRequest } from '../../types';

const schema = z.object({
  paciente_id: z.string().min(1, 'Seleccione un paciente'),
  fecha: z.string().min(1, 'Requerido'),
  hora: z.string().min(1, 'Requerido'),
  tipo_tratamiento: z.string().min(1, 'Requerido'),
  turno: z.enum(['AM', 'PM'], { message: 'Seleccione turno' }),
  observaciones: z.string().optional(),
  paquete_id: z.string().optional(),
});

interface CitaFormProps {
  onSubmit: (data: CreateCitaRequest) => void;
  loading: boolean;
  error: unknown;
}

export function CitaForm({ onSubmit, loading, error }: CitaFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
    watch,
  } = useForm<CreateCitaRequest>({
    resolver: zodResolver(schema),
    defaultValues: { turno: 'AM' },
  });

  const selectedPacienteId = watch('paciente_id');

  const { data: pacientesData } = useQuery({
    queryKey: ['pacientes-all'],
    queryFn: () => pacientesService.getAll(1, 100),
  });

  const { data: paquetesData } = useQuery({
    queryKey: ['paquetes-activos', selectedPacienteId],
    queryFn: () => paquetesService.getByPaciente(selectedPacienteId, true),
    enabled: !!selectedPacienteId,
  });

  const pacientes = pacientesData?.data ?? [];
  const paquetes = paquetesData?.data ?? [];
  const apiError = (error as any)?.response?.data?.error;

  const handleFormSubmit = (data: CreateCitaRequest) => {
    const submitData = { ...data };
    if (!submitData.paquete_id) {
      delete submitData.paquete_id;
    }
    onSubmit(submitData);
  };

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
      {apiError && (
        <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-danger">
          {apiError.detail || apiError.message}
        </div>
      )}

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Paciente *</label>
        <select {...register('paciente_id')} className="w-full">
          <option value="">Seleccionar paciente...</option>
          {pacientes.map((p) => (
            <option key={p.id} value={p.id}>
              {p.codigo} - {p.nombre_completo}
            </option>
          ))}
        </select>
        {errors.paciente_id && <p className="text-xs text-danger mt-1">{errors.paciente_id.message}</p>}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Fecha *</label>
          <input type="date" {...register('fecha')} className="w-full" />
          {errors.fecha && <p className="text-xs text-danger mt-1">{errors.fecha.message}</p>}
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Hora *</label>
          <input type="time" {...register('hora')} className="w-full" />
          {errors.hora && <p className="text-xs text-danger mt-1">{errors.hora.message}</p>}
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Tipo de tratamiento *</label>
          <input {...register('tipo_tratamiento')} className="w-full" placeholder="Masaje terapéutico" />
          {errors.tipo_tratamiento && <p className="text-xs text-danger mt-1">{errors.tipo_tratamiento.message}</p>}
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Turno *</label>
          <select {...register('turno')} className="w-full">
            <option value="AM">AM (Mañana)</option>
            <option value="PM">PM (Tarde)</option>
          </select>
          {errors.turno && <p className="text-xs text-danger mt-1">{errors.turno.message}</p>}
        </div>
      </div>

      {paquetes.length > 0 && (
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Paquete de Tratamiento (opcional)</label>
          <select {...register('paquete_id')} className="w-full">
            <option value="">Sin paquete</option>
            {paquetes.map((p) => (
              <option key={p.id} value={p.id}>
                {p.tipo_tratamiento} ({p.sesiones_completadas}/{p.total_sesiones} sesiones)
              </option>
            ))}
          </select>
        </div>
      )}

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Observaciones</label>
        <textarea {...register('observaciones')} className="w-full" rows={2} placeholder="Notas adicionales..." />
      </div>

      <Button type="submit" loading={loading} className="w-full">
        Agendar Cita
      </Button>
    </form>
  );
}
