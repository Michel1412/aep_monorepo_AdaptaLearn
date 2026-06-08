import { createApiClient } from '@adaptalearn/api-client';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

let token: string | null = localStorage.getItem('teacher_token');

export const api = createApiClient({
  baseUrl: API_URL,
  getToken: () => token,
});

export function setToken(t: string | null) {
  token = t;
  if (t) localStorage.setItem('teacher_token', t);
  else localStorage.removeItem('teacher_token');
}

export function getStoredTeacher() {
  const raw = localStorage.getItem('teacher');
  return raw ? JSON.parse(raw) : null;
}

export function setStoredTeacher(teacher: { id: string; name: string; email: string } | null) {
  if (teacher) localStorage.setItem('teacher', JSON.stringify(teacher));
  else localStorage.removeItem('teacher');
}
