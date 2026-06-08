import { describe, it, expect, beforeAll, afterAll, afterEach } from 'vitest';
import { setupServer } from 'msw/node';
import { http, HttpResponse } from 'msw';
import { createApiClient } from './index';

const server = setupServer(
  http.post('http://localhost:8080/api/v1/auth/student', () =>
    HttpResponse.json({
      token: 'test-token',
      student: { ra: '23159293-2', name: 'Vinicius' },
    }),
  ),
  http.get('http://localhost:8080/api/v1/students/23159293-2/dashboard', () =>
    HttpResponse.json({
      completed_activities: 0,
      total_activities: 3,
      planned_hours: 1.08,
      actual_hours: 0,
      subjects: [],
      exercise_breakdown: [],
    }),
  ),
);

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe('api-client', () => {
  const client = createApiClient({ baseUrl: 'http://localhost:8080' });

  it('loginStudent returns token', async () => {
    const res = await client.auth.loginStudent('23159293-2');
    expect(res.token).toBe('test-token');
    expect(res.student.ra).toBe('23159293-2');
  });

  it('getDashboard returns metrics', async () => {
    const clientWithAuth = createApiClient({
      baseUrl: 'http://localhost:8080',
      getToken: () => 'test-token',
    });
    const dashboard = await clientWithAuth.students.getDashboard('23159293-2');
    expect(dashboard.total_activities).toBe(3);
  });
});
