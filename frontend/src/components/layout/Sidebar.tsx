import { NavLink } from 'react-router-dom';
import { useAuthStore } from '../../store/auth';
import {
  LayoutDashboard,
  Users,
  UserPlus,
  Calendar,
  LogOut,
  X,
} from 'lucide-react';
import { cn } from '../../lib/utils';

interface SidebarProps {
  open: boolean;
  onClose: () => void;
}

const navItems = [
  { to: '/', label: 'Dashboard', icon: LayoutDashboard, roles: null },
  { to: '/pacientes', label: 'Pacientes', icon: UserPlus, roles: null },
  { to: '/citas', label: 'Citas', icon: Calendar, roles: null },
  { to: '/usuarios', label: 'Usuarios', icon: Users, roles: ['Administradora'] },
];

export function Sidebar({ open, onClose }: SidebarProps) {
  const { logout, hasRole, rolNombre } = useAuthStore();

  return (
    <>
      {open && (
        <div className="fixed inset-0 bg-black/30 z-40 lg:hidden" onClick={onClose} />
      )}

      <aside
        className={cn(
          'fixed top-0 left-0 z-50 h-full w-64 bg-white border-r border-border flex flex-col transition-transform duration-200 lg:translate-x-0 lg:static lg:z-auto',
          open ? 'translate-x-0' : '-translate-x-full'
        )}
      >
        <div className="flex items-center justify-between px-6 py-5 border-b border-border">
          <div>
            <h1 className="text-lg font-bold text-primary">Centro Caribel</h1>
            <p className="text-xs text-muted">Fisioterapia Estética</p>
          </div>
          <button className="lg:hidden p-1 hover:bg-gray-100 rounded" onClick={onClose}>
            <X className="h-5 w-5" />
          </button>
        </div>

        <nav className="flex-1 px-3 py-4 space-y-1">
          {navItems.map((item) => {
            if (item.roles && !hasRole(...item.roles)) return null;
            return (
              <NavLink
                key={item.to}
                to={item.to}
                onClick={onClose}
                className={({ isActive }) =>
                  cn(
                    'flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors',
                    isActive
                      ? 'bg-primary/10 text-primary'
                      : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
                  )
                }
                end={item.to === '/'}
              >
                <item.icon className="h-5 w-5" />
                {item.label}
              </NavLink>
            );
          })}
        </nav>

        <div className="px-3 py-4 border-t border-border">
          <p className="px-3 mb-2 text-xs text-muted truncate">{rolNombre}</p>
          <button
            onClick={logout}
            className="flex items-center gap-3 w-full px-3 py-2.5 rounded-lg text-sm font-medium text-gray-600 hover:bg-red-50 hover:text-danger transition-colors"
          >
            <LogOut className="h-5 w-5" />
            Cerrar sesión
          </button>
        </div>
      </aside>
    </>
  );
}
