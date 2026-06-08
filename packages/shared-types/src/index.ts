import { z } from 'zod';

export const AnswerTypeSchema = z.enum([
  'multiple_choice',
  'true_false',
  'short_text',
  'essay',
]);
export type AnswerType = z.infer<typeof AnswerTypeSchema>;

export const ComplementTypeSchema = z.enum(['none', 'video', 'image', 'link']);
export type ComplementType = z.infer<typeof ComplementTypeSchema>;

export const StudentSchema = z.object({
  ra: z.string(),
  name: z.string(),
});
export type Student = z.infer<typeof StudentSchema>;

export const TeacherSchema = z.object({
  id: z.string(),
  name: z.string(),
  email: z.string(),
});
export type Teacher = z.infer<typeof TeacherSchema>;

export const ClassSchema = z.object({
  id: z.string(),
  name: z.string(),
  subject_id: z.string().optional(),
  teacher_id: z.string().optional(),
  series: z.number(),
  turma: z.string(),
  subject_code: z.string().optional(),
  subject_name: z.string().optional(),
  student_count: z.number().optional(),
  activity_count: z.number().optional(),
});
export type Class = z.infer<typeof ClassSchema>;

export const ActivitySchema = z.object({
  id: z.string(),
  class_id: z.string(),
  title: z.string(),
  question: z.string(),
  estimated_minutes: z.number(),
  answer_type: AnswerTypeSchema,
  complement_type: ComplementTypeSchema,
  complement_url: z.string().nullable().optional(),
  week_start: z.string(),
  created_by: z.string().optional(),
  options: z.string().nullable().optional(),
  status: z.enum(['pending', 'completed']).optional(),
});
export type Activity = z.infer<typeof ActivitySchema>;

export const CreateActivitySchema = z.object({
  title: z.string().min(1, 'Título obrigatório'),
  question: z.string().min(1, 'Pergunta obrigatória'),
  estimated_minutes: z.number().min(1, 'Tempo deve ser positivo'),
  answer_type: AnswerTypeSchema,
  complement_type: ComplementTypeSchema.default('none'),
  complement_url: z.string().url().optional().or(z.literal('')),
  week_start: z.string().min(1),
  options: z.string().nullable().optional(),
});
export type CreateActivityInput = z.infer<typeof CreateActivitySchema>;

export const SubjectProgressSchema = z.object({
  class_id: z.string(),
  class_name: z.string(),
  subject_name: z.string(),
  completed: z.number(),
  total: z.number(),
});

export const ExerciseBreakdownSchema = z.object({
  activity_id: z.string(),
  title: z.string(),
  estimated_minutes: z.number(),
  time_spent_seconds: z.number(),
});

export const StudentDashboardSchema = z.object({
  completed_activities: z.number(),
  total_activities: z.number(),
  planned_hours: z.number(),
  actual_hours: z.number(),
  subjects: z.array(SubjectProgressSchema),
  exercise_breakdown: z.array(ExerciseBreakdownSchema),
});
export type StudentDashboard = z.infer<typeof StudentDashboardSchema>;

export const ClassMetricsSchema = z.object({
  ra: z.string(),
  name: z.string(),
  completed_activities: z.number(),
  total_time_seconds: z.number(),
});
export type ClassMetrics = z.infer<typeof ClassMetricsSchema>;

export const TeacherDashboardSchema = z.object({
  classes: z.number(),
  activities: z.number(),
  students: z.number(),
});
export type TeacherDashboard = z.infer<typeof TeacherDashboardSchema>;

export function formatClassName(code: string, series: number, turma: string): string {
  return `${code}-S${series}-${turma}`;
}
