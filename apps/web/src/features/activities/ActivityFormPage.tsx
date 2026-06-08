import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { CreateActivitySchema, type CreateActivityInput } from '@adaptalearn/shared-types';
import { useMutation } from '@tanstack/react-query';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { api } from '../../lib/api';

function getMonday(): string {
  const now = new Date();
  const day = now.getDay();
  const diff = day === 0 ? -6 : 1 - day;
  const monday = new Date(now);
  monday.setDate(now.getDate() + diff);
  return monday.toISOString().split('T')[0];
}

export function ActivityFormPage() {
  const { id: classId } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const { register, handleSubmit, watch, formState: { errors } } = useForm<CreateActivityInput>({
    resolver: zodResolver(CreateActivitySchema),
    defaultValues: {
      answer_type: 'essay',
      complement_type: 'none',
      week_start: getMonday(),
      estimated_minutes: 30,
    },
  });

  const answerType = watch('answer_type');
  const complementType = watch('complement_type');

  const mutation = useMutation({
    mutationFn: (data: CreateActivityInput) => api.teachers.createActivity(classId!, data),
    onSuccess: () => navigate(`/classes/${classId}`),
  });

  return (
    <div className="max-w-2xl">
      <Link to={`/classes/${classId}`} className="text-primary text-sm hover:underline">← Voltar</Link>
      <h2 className="text-2xl font-bold mt-2 mb-6">Nova Atividade</h2>
      <form onSubmit={handleSubmit((d) => mutation.mutate(d))} className="space-y-4 bg-white p-6 rounded border">
        <label className="block">
          <span className="text-sm font-medium">Título</span>
          <input {...register('title')} className="mt-1 w-full border rounded px-3 py-2" />
          {errors.title && <p className="text-red-600 text-xs">{errors.title.message}</p>}
        </label>
        <label className="block">
          <span className="text-sm font-medium">Pergunta</span>
          <textarea {...register('question')} rows={4} className="mt-1 w-full border rounded px-3 py-2" />
          {errors.question && <p className="text-red-600 text-xs">{errors.question.message}</p>}
        </label>
        <label className="block">
          <span className="text-sm font-medium">Tempo estimado (min)</span>
          <input type="number" {...register('estimated_minutes', { valueAsNumber: true })} className="mt-1 w-full border rounded px-3 py-2" />
        </label>
        <label className="block">
          <span className="text-sm font-medium">Tipo de resposta</span>
          <select {...register('answer_type')} className="mt-1 w-full border rounded px-3 py-2">
            <option value="essay">Dissertativa</option>
            <option value="short_text">Texto curto</option>
            <option value="true_false">Verdadeiro/Falso</option>
            <option value="multiple_choice">Múltipla escolha</option>
          </select>
        </label>
        {answerType === 'multiple_choice' && (
          <label className="block">
            <span className="text-sm font-medium">Opções (JSON)</span>
            <textarea {...register('options')} placeholder='["A","B","C"]' className="mt-1 w-full border rounded px-3 py-2" />
          </label>
        )}
        <label className="block">
          <span className="text-sm font-medium">Complemento</span>
          <select {...register('complement_type')} className="mt-1 w-full border rounded px-3 py-2">
            <option value="none">Nenhum</option>
            <option value="video">Vídeo</option>
            <option value="image">Imagem</option>
            <option value="link">Link</option>
          </select>
        </label>
        {complementType !== 'none' && (
          <label className="block">
            <span className="text-sm font-medium">URL do complemento</span>
            <input {...register('complement_url')} className="mt-1 w-full border rounded px-3 py-2" />
          </label>
        )}
        <label className="block">
          <span className="text-sm font-medium">Semana de publicação</span>
          <input type="date" {...register('week_start')} className="mt-1 w-full border rounded px-3 py-2" />
        </label>
        {mutation.error && (
          <p className="text-red-600 text-sm">{(mutation.error as Error).message}</p>
        )}
        <button type="submit" disabled={mutation.isPending} className="bg-primary text-white px-6 py-2 rounded">
          {mutation.isPending ? 'Salvando...' : 'Salvar'}
        </button>
      </form>
    </div>
  );
}
