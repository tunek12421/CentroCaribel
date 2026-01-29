import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useQuery } from '@tanstack/react-query';
import { usuariosService } from '../../services/usuarios';
import { Button } from '../../components/ui/Button';
import type { CreateUsuarioRequest, UpdateUsuarioRequest, Usuario } from '../../types';

const createSchema = z.object({
  nombre_completo: z.string().min(3, 'Mínimo 3 caracteres'),
  email: z.string().email('Email inválido'),
  password: z.string().min(8, 'Mínimo 8 caracteres'),
  rol_id: z.string().min(1, 'Seleccione un rol'),
});

const updateSchema = z.object({
  nombre_completo: z.string().min(3, 'Mínimo 3 caracteres'),
  email: z.string().email('Email inválido'),
  rol_id: z.string().min(1, 'Seleccione un rol'),
  activo: z.boolean(),
});

interface UsuarioFormProps {
  usuario?: Usuario;
  onSubmit?: (data: CreateUsuarioRequest) => void;
  onSubmitUpdate?: (data: UpdateUsuarioRequest) => void;
  loading: boolean;
  error: unknown;
}

export function UsuarioForm({ usuario, onSubmit, onSubmitUpdate, loading, error }: UsuarioFormProps) {
  const isEdit = !!usuario;

  const { data: roles } = useQuery({
    queryKey: ['roles'],
    queryFn: () => usuariosService.getRoles(),
  });

  const createForm = useForm<CreateUsuarioRequest>({
    resolver: zodResolver(createSchema),
  });

  const updateForm = useForm({
    resolver: zodResolver(updateSchema),
    defaultValues: usuario
      ? { nombre_completo: usuario.nombre_completo, email: usuario.email, rol_id: usuario.rol_id, activo: usuario.activo }
      : undefined,
  });

  const apiError = (error as any)?.response?.data?.error;

  if (isEdit) {
    return (
      <form onSubmit={updateForm.handleSubmit((data) => onSubmitUpdate?.(data))} className="space-y-4">
        {apiError && (
          <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-danger">
            {apiError.detail || apiError.message}
          </div>
        )}

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Nombre completo *</label>
          <input {...updateForm.register('nombre_completo')} className="w-full" />
          {updateForm.formState.errors.nombre_completo && (
            <p className="text-xs text-danger mt-1">{updateForm.formState.errors.nombre_completo.message}</p>
          )}
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Email *</label>
          <input type="email" {...updateForm.register('email')} className="w-full" />
          {updateForm.formState.errors.email && (
            <p className="text-xs text-danger mt-1">{updateForm.formState.errors.email.message}</p>
          )}
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Rol *</label>
          <select {...updateForm.register('rol_id')} className="w-full">
            <option value="">Seleccionar rol...</option>
            {(roles ?? []).map((r) => (
              <option key={r.id} value={r.id}>{r.nombre}</option>
            ))}
          </select>
        </div>

        <div className="flex items-center gap-2">
          <input type="checkbox" id="activo" {...updateForm.register('activo')} className="h-4 w-4" />
          <label htmlFor="activo" className="text-sm text-gray-700">Activo</label>
        </div>

        <Button type="submit" loading={loading} className="w-full">Guardar Cambios</Button>
      </form>
    );
  }

  return (
    <form onSubmit={createForm.handleSubmit((data) => onSubmit?.(data))} className="space-y-4">
      {apiError && (
        <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-danger">
          {apiError.detail || apiError.message}
        </div>
      )}

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Nombre completo *</label>
        <input {...createForm.register('nombre_completo')} className="w-full" placeholder="Nombre del usuario" />
        {createForm.formState.errors.nombre_completo && (
          <p className="text-xs text-danger mt-1">{createForm.formState.errors.nombre_completo.message}</p>
        )}
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Email *</label>
        <input type="email" {...createForm.register('email')} className="w-full" placeholder="correo@centrocaribel.com" />
        {createForm.formState.errors.email && (
          <p className="text-xs text-danger mt-1">{createForm.formState.errors.email.message}</p>
        )}
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Contraseña *</label>
        <input type="password" {...createForm.register('password')} className="w-full" placeholder="Mínimo 8 caracteres" />
        {createForm.formState.errors.password && (
          <p className="text-xs text-danger mt-1">{createForm.formState.errors.password.message}</p>
        )}
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Rol *</label>
        <select {...createForm.register('rol_id')} className="w-full">
          <option value="">Seleccionar rol...</option>
          {(roles ?? []).map((r) => (
            <option key={r.id} value={r.id}>{r.nombre}</option>
          ))}
        </select>
        {createForm.formState.errors.rol_id && (
          <p className="text-xs text-danger mt-1">{createForm.formState.errors.rol_id.message}</p>
        )}
      </div>

      <Button type="submit" loading={loading} className="w-full">Crear Usuario</Button>
    </form>
  );
}
