import { useQuery } from '@tanstack/react-query';
import { api } from '../../lib/api';

export function DashboardPage() {
  const { data, isLoading } = useQuery({
    queryKey: ['teacher-dashboard'],
    queryFn: () => api.teachers.getDashboard(),
  });

  if (isLoading) return <p>Carregando...</p>;

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">Dashboard</h2>
      <div className="grid grid-cols-3 gap-4">
        <StatCard title="Turmas" value={data?.classes ?? 0} />
        <StatCard title="Atividades (semana)" value={data?.activities ?? 0} />
        <StatCard title="Alunos matriculados" value={data?.students ?? 0} />
      </div>
    </div>
  );
}

function StatCard({ title, value }: { title: string; value: number }) {
  return (
    <div className="bg-white p-6 rounded-lg shadow-sm border">
      <p className="text-gray-500 text-sm">{title}</p>
      <p className="text-3xl font-bold mt-2">{value}</p>
    </div>
  );
}
