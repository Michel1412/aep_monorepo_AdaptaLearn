import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Link, useParams } from 'react-router-dom';
import { api } from '../../lib/api';

type Tab = 'activities' | 'students' | 'metrics';

export function ClassDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [tab, setTab] = useState<Tab>('activities');
  const queryClient = useQueryClient();

  const { data: classes } = useQuery({
    queryKey: ['teacher-classes'],
    queryFn: () => api.teachers.getClasses(),
  });
  const classInfo = classes?.find((c) => c.id === id);

  const { data: students } = useQuery({
    queryKey: ['class-students', id],
    queryFn: () => api.teachers.getStudents(id!),
    enabled: !!id && tab === 'students',
  });

  const { data: metrics } = useQuery({
    queryKey: ['class-metrics', id],
    queryFn: () => api.teachers.getMetrics(id!),
    enabled: !!id && tab === 'metrics',
  });

  const [ra, setRa] = useState('');
  const enrollMutation = useMutation({
    mutationFn: (studentRa: string) => api.teachers.enrollStudent(id!, studentRa),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['class-students', id] });
      setRa('');
    },
  });

  return (
    <div>
      <Link to="/classes" className="text-primary text-sm hover:underline">← Voltar</Link>
      <h2 className="text-2xl font-bold mt-2">{classInfo?.name ?? 'Turma'}</h2>
      <p className="text-gray-500 mb-6">{classInfo?.subject_name}</p>

      <div className="flex gap-4 border-b mb-6">
        {(['activities', 'students', 'metrics'] as Tab[]).map((t) => (
          <button
            key={t}
            onClick={() => setTab(t)}
            className={`pb-2 px-1 capitalize ${
              tab === t ? 'border-b-2 border-primary font-medium' : 'text-gray-500'
            }`}
          >
            {t === 'activities' ? 'Atividades' : t === 'students' ? 'Alunos' : 'Métricas'}
          </button>
        ))}
      </div>

      {tab === 'activities' && (
        <div>
          <Link
            to={`/classes/${id}/activities/new`}
            className="inline-block bg-primary text-white px-4 py-2 rounded mb-4"
          >
            Nova atividade
          </Link>
          <p className="text-gray-500">Atividades da semana atual são exibidas no app do aluno.</p>
        </div>
      )}

      {tab === 'students' && (
        <div>
          <form
            onSubmit={(e) => { e.preventDefault(); enrollMutation.mutate(ra); }}
            className="flex gap-2 mb-4"
          >
            <input
              value={ra}
              onChange={(e) => setRa(e.target.value)}
              placeholder="RA do aluno"
              className="border rounded px-3 py-2"
            />
            <button type="submit" className="bg-primary text-white px-4 py-2 rounded">
              Matricular
            </button>
          </form>
          <ul className="space-y-2">
            {students?.map((s) => (
              <li key={s.ra} className="bg-white p-3 rounded border">
                {s.name} — RA {s.ra}
              </li>
            ))}
          </ul>
        </div>
      )}

      {tab === 'metrics' && (
        <table className="w-full bg-white rounded border">
          <thead>
            <tr className="border-b text-left">
              <th className="p-3">Aluno</th>
              <th className="p-3">Atividades</th>
              <th className="p-3">Tempo total</th>
            </tr>
          </thead>
          <tbody>
            {metrics?.map((m) => (
              <tr key={m.ra} className="border-b">
                <td className="p-3">{m.name}</td>
                <td className="p-3">{m.completed_activities}</td>
                <td className="p-3">{Math.round(m.total_time_seconds / 60)} min</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
