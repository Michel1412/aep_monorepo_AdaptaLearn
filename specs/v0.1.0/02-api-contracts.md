# Contratos de API — v0.1.0

Base URL: `http://localhost:8080/api/v1`

Autenticação: Bearer JWT no header `Authorization`.

## Auth

### POST /auth/student
```json
// Request
{ "ra": "23159293-2" }

// Response 200
{
  "token": "eyJ...",
  "student": { "ra": "23159293-2", "name": "Vinicius Reginaldo Ferrarini" }
}

// Response 401
{ "error": "RA não encontrado" }
```

### POST /auth/teacher
```json
// Request
{ "email": "maurilio@adaptalearn.local", "password": "senha123" }

// Response 200
{
  "token": "eyJ...",
  "teacher": { "id": "...", "name": "Maurilio", "email": "..." }
}
```

## Student

### GET /students/{ra}/dashboard
Headers: Authorization (student token)

```json
// Response 200
{
  "completed_activities": 2,
  "total_activities": 5,
  "planned_hours": 3.5,
  "actual_hours": 1.2,
  "subjects": [
    {
      "class_id": "...",
      "class_name": "TeoCom-S7-B",
      "subject_name": "Teoria da Computação",
      "completed": 2,
      "total": 3
    }
  ],
  "exercise_breakdown": [
    {
      "activity_id": "...",
      "title": "Autômatos finitos",
      "estimated_minutes": 30,
      "time_spent_seconds": 1200
    }
  ]
}
```

### GET /students/{ra}/classes
```json
// Response 200
[
  {
    "id": "...",
    "name": "TeoCom-S7-B",
    "subject_code": "TeoCom",
    "subject_name": "Teoria da Computação",
    "series": 7,
    "turma": "B"
  }
]
```

### GET /students/{ra}/classes/{classId}/activities
Query: `week=current` (default)

```json
// Response 200
[
  {
    "id": "...",
    "title": "Autômatos finitos",
    "question": "Defina autômato finito determinístico...",
    "estimated_minutes": 30,
    "answer_type": "essay",
    "complement_type": "video",
    "complement_url": "https://youtube.com/...",
    "week_start": "2026-06-02",
    "status": "pending",
    "options": null
  }
]
```

## Activities / Sessions

### POST /activities/{id}/sessions/start
```json
// Response 201
{ "session_id": "...", "started_at": "2026-06-08T10:00:00Z" }
```

### POST /activities/{id}/sessions/{sessionId}/submit
```json
// Request
{ "answer": "Um AFD é uma tupla (Q, Σ, δ, q0, F)...", "time_spent_seconds": 1200 }

// Response 201
{
  "submission_id": "...",
  "time_spent_seconds": 1200,
  "submitted_at": "2026-06-08T10:20:00Z"
}
```

## Teacher

### GET /teachers/me/classes
```json
// Response 200
[
  {
    "id": "...",
    "name": "TeoCom-S7-B",
    "subject_code": "TeoCom",
    "subject_name": "Teoria da Computação",
    "series": 7,
    "turma": "B",
    "student_count": 1,
    "activity_count": 3
  }
]
```

### GET /teachers/me/classes/{classId}/students
```json
// Response 200
[
  { "ra": "23159293-2", "name": "Vinicius Reginaldo Ferrarini" }
]
```

### POST /teachers/me/classes/{classId}/students
```json
// Request
{ "ra": "23159293-2" }

// Response 201
{ "ra": "23159293-2", "name": "Vinicius Reginaldo Ferrarini" }
```

### GET /teachers/me/classes/{classId}/metrics
```json
// Response 200
[
  {
    "ra": "23159293-2",
    "name": "Vinicius Reginaldo Ferrarini",
    "completed_activities": 2,
    "total_time_seconds": 3600
  }
]
```

### POST /teachers/me/classes/{classId}/activities
```json
// Request
{
  "title": "Autômatos finitos",
  "question": "Defina autômato finito determinístico...",
  "estimated_minutes": 30,
  "answer_type": "essay",
  "complement_type": "video",
  "complement_url": "https://youtube.com/...",
  "week_start": "2026-06-02",
  "options": null
}

// Response 201
{ "id": "...", ... }
```

### PUT /activities/{id}
Mesmo payload do POST (parcial permitido).

## Códigos de erro

| Código | Significado |
|--------|-------------|
| 400 | Payload inválido |
| 401 | Não autenticado |
| 403 | Sem permissão |
| 404 | Recurso não encontrado |
| 409 | Conflito (ex.: já matriculado) |

Formato: `{ "error": "mensagem descritiva" }`
