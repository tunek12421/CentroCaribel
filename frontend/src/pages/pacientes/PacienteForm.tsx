import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Button } from '../../components/ui/Button';
import type { CreatePacienteRequest } from '../../types';

const schema = z.object({
  nombre_completo: z.string().min(3, 'Mínimo 3 caracteres'),
  ci: z.string().min(4, 'CI inválido'),
  fecha_nacimiento: z.string().min(1, 'Requerido'),
  celular: z.string().min(7, 'Celular inválido'),
  direccion: z.string().optional(),
});

interface PacienteFormProps {
  onSubmit: (data: CreatePacienteRequest) => void;
  loading: boolean;
  error: unknown;
}

export function PacienteForm({ onSubmit, loading, error }: PacienteFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreatePacienteRequest>({
    resolver: zodResolver(schema),
  });

  const apiError = (error as any)?.response?.data?.error;

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      {apiError && (
        <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-danger">
          {apiError.detail || apiError.message}
        </div>
      )}

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Nombre completo *</label>
        <input {...register('nombre_completo')} className="w-full" placeholder="Nombre del paciente" />
        {errors.nombre_completo && <p className="text-xs text-danger mt-1">{errors.nombre_completo.message}</p>}
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">CI *</label>
          <input {...register('ci')} className="w-full" placeholder="1234567" />
          {errors.ci && <p className="text-xs text-danger mt-1">{errors.ci.message}</p>}
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Fecha de nacimiento *</label>
          <input type="date" {...register('fecha_nacimiento')} className="w-full" />
          {errors.fecha_nacimiento && <p className="text-xs text-danger mt-1">{errors.fecha_nacimiento.message}</p>}
        </div>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Celular *</label>
        <input {...register('celular')} className="w-full" placeholder="70012345" />
        {errors.celular && <p className="text-xs text-danger mt-1">{errors.celular.message}</p>}
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Dirección</label>
        <input {...register('direccion')} className="w-full" placeholder="Av. América #123" />
      </div>

      <Button type="submit" loading={loading} className="w-full">
        Registrar Paciente
      </Button>
    </form>
  );
}
