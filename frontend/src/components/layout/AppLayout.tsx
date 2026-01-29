import { useState } from 'react';
import { Outlet, useLocation } from 'react-router-dom';
import { Sidebar } from './Sidebar';
import { Header } from './Header';

const titles: Record<string, string> = {
  '/': 'Dashboard',
  '/pacientes': 'Pacientes',
  '/citas': 'Citas',
  '/usuarios': 'Usuarios',
};

function getTitle(path: string) {
  if (titles[path]) return titles[path];
  if (path.startsWith('/pacientes/')) return 'Detalle de Paciente';
  return 'Centro Caribel';
}

export function AppLayout() {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const location = useLocation();

  return (
    <div className="flex h-screen overflow-hidden">
      <Sidebar open={sidebarOpen} onClose={() => setSidebarOpen(false)} />
      <div className="flex-1 flex flex-col overflow-hidden">
        <Header
          title={getTitle(location.pathname)}
          onMenuClick={() => setSidebarOpen(true)}
        />
        <main className="flex-1 overflow-y-auto p-4 lg:p-6">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
