# Matriz de Testes — v0.1.0

## Backend Go (`apps/api`)

| Pacote | Teste | Tipo |
|--------|-------|------|
| service/dashboard | Cálculo horas previstas/realizadas | Unit |
| service/activity | Validação de campos obrigatórios | Unit |
| service/session | Validação tempo min/max | Unit |
| handler/auth | Login student/teacher | Integration |
| handler/students | Dashboard, classes, activities | Integration |
| handler/teachers | CRUD turmas, atividades, alunos | Integration |
| seed | Maurilio, 12 turmas, Vinicius em S7-B | Integration |

Ferramentas: `go test`, `testify/assert`, SQLite in-memory.

## packages/api-client

| Teste | Escopo |
|-------|--------|
| auth.loginStudent | POST /auth/student |
| auth.loginTeacher | POST /auth/teacher |
| students.getDashboard | GET dashboard |
| teachers.createActivity | POST activity |

Ferramentas: Vitest, MSW.

## apps/web

| Teste | Escopo |
|-------|--------|
| LoginForm | Render, validação, submit |
| ActivityForm | Validação Zod, campos condicionais |
| ClassGrid | Filtro por série |
| EnrollStudentForm | Submit RA |

Ferramentas: Vitest, React Testing Library, MSW.

## apps/mobile

| Teste | Escopo |
|-------|--------|
| LoginScreen | Render, validação RA |
| useStudyTimer | Incremento, pause, reset |
| QuestionInput | Render por answer_type |
| DashboardScreen | Render cards com mock MSW |

Ferramentas: Jest, React Native Testing Library, MSW.

## Critério de CI local

```bash
pnpm test        # todos os pacotes via turbo
cd apps/api && go test ./...
```

Todos devem passar antes de considerar v0.1.0 completa.

## Fora do escopo v0.1.0

- E2E com Detox ou Playwright
- Testes de carga
- Testes visuais (snapshot)
