# Modelo de Domínio — v0.1.0

## Entidades

### Teacher (Professor)
| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | Identificador único |
| name | string | Nome completo |
| email | string | Email de login (único) |
| password_hash | string | Senha hasheada (bcrypt) |

### Student (Aluno)
| Campo | Tipo | Descrição |
|-------|------|-----------|
| ra | string | Registro Acadêmico (PK) |
| name | string | Nome completo |

### Subject (Matéria)
| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | Identificador |
| code | string | Código curto (ex.: `TeoCom`) |
| name | string | Nome completo |

### Class (Turma)
| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | Identificador |
| name | string | Nome composto (ex.: `TeoCom-S7-B`) |
| subject_id | UUID | FK para Subject |
| teacher_id | UUID | FK para Teacher |
| series | int | Série (4, 5, 6, 7) |
| turma | string | Turma (A, B, C) |

**Regra de nomenclatura:** `{subject.code}-S{series}-{turma}`

### Enrollment (Matrícula)
| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | Identificador |
| student_ra | string | FK para Student |
| class_id | UUID | FK para Class |

### Activity (Atividade)
| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | Identificador |
| class_id | UUID | FK para Class |
| title | string | Título curto |
| question | string | Enunciado |
| estimated_minutes | int | Tempo estimado em minutos |
| answer_type | enum | Tipo de resposta |
| complement_type | enum | Tipo de complemento |
| complement_url | string? | URL do complemento |
| week_start | date | Início da semana de publicação |
| created_by | UUID | FK para Teacher |

**answer_type:** `multiple_choice` | `true_false` | `short_text` | `essay`

**complement_type:** `none` | `video` | `image` | `link`

Para `multiple_choice`, opções são armazenadas em JSON no campo `options`.

### StudySession (Sessão de Estudo)
| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | Identificador |
| activity_id | UUID | FK para Activity |
| student_ra | string | FK para Student |
| started_at | datetime | Início da sessão |
| ended_at | datetime? | Fim (null enquanto ativa) |

### Submission (Submissão)
| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | Identificador |
| activity_id | UUID | FK para Activity |
| student_ra | string | FK para Student |
| session_id | UUID | FK para StudySession |
| answer | string | Resposta do aluno |
| time_spent_seconds | int | Tempo gasto |
| submitted_at | datetime | Data/hora do envio |

## Relacionamentos

```
Teacher 1──* Class
Subject 1──* Class
Class 1──* Enrollment *──1 Student
Class 1──* Activity
Activity 1──* StudySession
Activity 1──* Submission
Student 1──* Submission
StudySession 1──1 Submission
```

## Regras de negócio

1. Um aluno só vê atividades das turmas em que está matriculado
2. Um professor só gerencia turmas que leciona
3. Tempo de estudo: mínimo 10s, máximo 4h por sessão
4. Atividades da semana: `week_start <= hoje < week_start + 7 dias`
5. Um aluno pode ter apenas uma submissão por atividade (v0.1.0)
