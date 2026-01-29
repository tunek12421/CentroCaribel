import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { usuariosService } from '../../services/usuarios';
import { Card } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { Modal } from '../../components/ui/Modal';
import { Badge } from '../../components/ui/Badge';
import { Spinner } from '../../components/ui/Spinner';
import { Plus, Pencil, Trash2 } from 'lucide-react';
import { formatDate } from '../../lib/utils';
import { UsuarioForm } from './UsuarioForm';
import type { Usuario, UpdateUsuarioRequest } from '../../types';

export function UsuariosPage() {
  const [page, setPage] = useState(1);
  const [showForm, setShowForm] = useState(false);
  const [editing, setEditing] = useState<Usuario | null>(null);
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['usuarios', page],
    queryFn: () => usuariosService.getAll(page),
  });

  const createMutation = useMutation({
    mutationFn: usuariosService.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['usuarios'] });
      setShowForm(false);
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateUsuarioRequest }) =>
      usuariosService.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['usuarios'] });
      setEditing(null);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: usuariosService.delete,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['usuarios'] }),
  });

  if (isLoading) return <Spinner />;

  const usuarios = data?.data ?? [];
  const meta = data?.meta;

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <p className="text-sm text-muted">{meta?.total ?? 0} usuarios</p>
        <Button onClick={() => setShowForm(true)}>
          <Plus className="h-4 w-4 mr-2" />
          Nuevo Usuario
        </Button>
      </div>

      <Card>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border bg-gray-50/50">
                <th className="text-left px-6 py-3 font-medium text-muted">Nombre</th>
                <th className="text-left px-6 py-3 font-medium text-muted hidden sm:table-cell">Email</th>
                <th className="text-left px-6 py-3 font-medium text-muted">Rol</th>
                <th className="text-left px-6 py-3 font-medium text-muted">Estado</th>
                <th className="text-left px-6 py-3 font-medium text-muted hidden md:table-cell">Creado</th>
                <th className="text-right px-6 py-3 font-medium text-muted">Acciones</th>
              </tr>
            </thead>
            <tbody>
              {usuarios.map((u) => (
                <tr key={u.id} className="border-b border-border hover:bg-gray-50/50">
                  <td className="px-6 py-3 font-medium">{u.nombre_completo}</td>
                  <td className="px-6 py-3 hidden sm:table-cell text-muted">{u.email}</td>
                  <td className="px-6 py-3">
                    <Badge className="bg-primary/10 text-primary">{u.rol?.nombre || '-'}</Badge>
                  </td>
                  <td className="px-6 py-3">
                    <Badge className={u.activo ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}>
                      {u.activo ? 'Activo' : 'Inactivo'}
                    </Badge>
                  </td>
                  <td className="px-6 py-3 hidden md:table-cell text-muted">{formatDate(u.created_at)}</td>
                  <td className="px-6 py-3 text-right">
                    <div className="flex justify-end gap-1">
                      <Button variant="ghost" size="sm" onClick={() => setEditing(u)}>
                        <Pencil className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => {
                          if (confirm('¿Desactivar este usuario?')) {
                            deleteMutation.mutate(u.id);
                          }
                        }}
                      >
                        <Trash2 className="h-4 w-4 text-danger" />
                      </Button>
                    </div>
                  </td>
                </tr>
              ))}
              {usuarios.length === 0 && (
                <tr>
                  <td colSpan={6} className="px-6 py-8 text-center text-muted">
                    No hay usuarios
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

      <Modal open={showForm} onClose={() => setShowForm(false)} title="Nuevo Usuario">
        <UsuarioForm
          onSubmit={(data) => createMutation.mutate(data)}
          loading={createMutation.isPending}
          error={createMutation.error}
        />
      </Modal>

      <Modal open={!!editing} onClose={() => setEditing(null)} title="Editar Usuario">
        {editing && (
          <UsuarioForm
            usuario={editing}
            onSubmitUpdate={(data) => updateMutation.mutate({ id: editing.id, data })}
            loading={updateMutation.isPending}
            error={updateMutation.error}
          />
        )}
      </Modal>
    </div>
  );
}
