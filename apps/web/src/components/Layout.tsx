import { Link, Outlet, useNavigate } from 'react-router-dom';
import { getStoredTeacher, setStoredTeacher, setToken } from '../lib/api';

export function Layout() {
  const teacher = getStoredTeacher();
  const navigate = useNavigate();

  function logout() {
    setToken(null);
    setStoredTeacher(null);
    navigate('/login');
  }

  return (
    <div className="min-h-screen flex">
      <aside className="w-64 bg-white border-r border-gray-200 p-6">
        <h1 className="text-xl font-bold text-primary mb-8">AdaptaLearn</h1>
        <nav className="space-y-2">
          <Link to="/" className="block px-3 py-2 rounded hover:bg-gray-100">Dashboard</Link>
          <Link to="/classes" className="block px-3 py-2 rounded hover:bg-gray-100">Turmas</Link>
        </nav>
        <div className="mt-auto pt-8">
          <p className="text-sm text-gray-500 mb-2">{teacher?.name}</p>
          <button onClick={logout} className="text-sm text-red-600 hover:underline">Sair</button>
        </div>
      </aside>
      <main className="flex-1 p-8">
        <Outlet />
      </main>
    </div>
  );
}
