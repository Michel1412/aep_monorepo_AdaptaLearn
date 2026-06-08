import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { api } from '../../lib/api';

export function ClassesPage() {
  const [seriesFilter, setSeriesFilter] = useState<number | 'all'>('all');
  const { data: classes, isLoading } = useQuery({
    queryKey: ['teacher-classes'],
    queryFn: () => api.teachers.getClasses(),
  });

  const filtered = classes?.filter(
    (c) => seriesFilter === 'all' || c.series === seriesFilter,
  );

  if (isLoading) return <p>Carregando...</p>;

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">Minhas Turmas</h2>
      <div className="flex gap-2 mb-6">
        {(['all', 4, 5, 6, 7] as const).map((s) => (
          <button
            key={String(s)}
            onClick={() => setSeriesFilter(s)}
            className={`px-3 py-1 rounded text-sm ${
              seriesFilter === s ? 'bg-primary text-white' : 'bg-white border'
            }`}
          >
            {s === 'all' ? 'Todas' : `Série ${s}`}
          </button>
        ))}
      </div>
      <div className="grid grid-cols-3 gap-4">
        {filtered?.map((c) => (
          <Link
            key={c.id}
            to={`/classes/${c.id}`}
            className="bg-white p-4 rounded-lg border hover:border-primary transition"
          >
            <h3 className="font-semibold">{c.name}</h3>
            <p className="text-sm text-gray-500">{c.subject_name}</p>
            <p className="text-xs text-gray-400 mt-2">
              {c.student_count} alunos · {c.activity_count} atividades
            </p>
          </Link>
        ))}
      </div>
    </div>
  );
}
