import type {
  Activity,
  Class,
  ClassMetrics,
  CreateActivityInput,
  Student,
  StudentDashboard,
  Teacher,
  TeacherDashboard,
} from '@adaptalearn/shared-types';

export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

export interface ApiClientConfig {
  baseUrl: string;
  getToken?: () => string | null;
}

export function createApiClient(config: ApiClientConfig) {
  const { baseUrl, getToken } = config;

  async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options.headers as Record<string, string>),
    };
    const token = getToken?.();
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }
    const res = await fetch(`${baseUrl}${path}`, { ...options, headers });
    if (!res.ok) {
      const body = await res.json().catch(() => ({ error: res.statusText }));
      throw new ApiError(res.status, body.error || res.statusText);
    }
    if (res.status === 204) return undefined as T;
    return res.json();
  }

  return {
    auth: {
      loginStudent: (ra: string) =>
        request<{ token: string; student: Student }>('/api/v1/auth/student', {
          method: 'POST',
          body: JSON.stringify({ ra }),
        }),
      loginTeacher: (email: string, password: string) =>
        request<{ token: string; teacher: Teacher }>('/api/v1/auth/teacher', {
          method: 'POST',
          body: JSON.stringify({ email, password }),
        }),
    },
    students: {
      getDashboard: (ra: string) =>
        request<StudentDashboard>(`/api/v1/students/${ra}/dashboard`),
      getClasses: (ra: string) =>
        request<Class[]>(`/api/v1/students/${ra}/classes`),
      getActivities: (ra: string, classId: string) =>
        request<Activity[]>(`/api/v1/students/${ra}/classes/${classId}/activities`),
    },
    activities: {
      startSession: (activityId: string) =>
        request<{ session_id: string; started_at: string }>(
          `/api/v1/activities/${activityId}/sessions/start`,
          { method: 'POST' },
        ),
      submit: (activityId: string, sessionId: string, answer: string, timeSpentSeconds: number) =>
        request<{ submission_id: string; time_spent_seconds: number; submitted_at: string }>(
          `/api/v1/activities/${activityId}/sessions/${sessionId}/submit`,
          {
            method: 'POST',
            body: JSON.stringify({ answer, time_spent_seconds: timeSpentSeconds }),
          },
        ),
    },
    teachers: {
      getDashboard: () => request<TeacherDashboard>('/api/v1/teachers/me/dashboard'),
      getClasses: () => request<Class[]>('/api/v1/teachers/me/classes'),
      getStudents: (classId: string) =>
        request<Student[]>(`/api/v1/teachers/me/classes/${classId}/students`),
      enrollStudent: (classId: string, ra: string) =>
        request<Student>(`/api/v1/teachers/me/classes/${classId}/students`, {
          method: 'POST',
          body: JSON.stringify({ ra }),
        }),
      getMetrics: (classId: string) =>
        request<ClassMetrics[]>(`/api/v1/teachers/me/classes/${classId}/metrics`),
      createActivity: (classId: string, input: CreateActivityInput) =>
        request<Activity>(`/api/v1/teachers/me/classes/${classId}/activities`, {
          method: 'POST',
          body: JSON.stringify(input),
        }),
      updateActivity: (activityId: string, input: CreateActivityInput) =>
        request<Activity>(`/api/v1/activities/${activityId}`, {
          method: 'PUT',
          body: JSON.stringify(input),
        }),
    },
  };
}

export type ApiClient = ReturnType<typeof createApiClient>;
