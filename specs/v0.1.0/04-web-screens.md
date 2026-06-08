# Telas Web (Professor) — v0.1.0

App: `apps/web` — React Router

## Rotas

| Rota | Tela | Descrição |
|------|------|-----------|
| `/login` | Login | Email e senha |
| `/` | Dashboard | Resumo geral |
| `/classes` | Minhas turmas | Grid de turmas |
| `/classes/:id` | Detalhe da turma | Abas Atividades / Alunos / Métricas |
| `/classes/:id/activities/new` | Nova atividade | Formulário template |
| `/classes/:id/activities/:activityId/edit` | Editar atividade | Formulário preenchido |

## Detalhamento por tela

### Login
- **Campos:** Email, Senha
- **Ações:** "Entrar" → POST `/auth/teacher`
- **Estados:** loading, erro, redirect para dashboard

### Dashboard
- **Cards:** Total de turmas, atividades publicadas na semana, alunos matriculados
- **Lista rápida:** Turmas com atividade recente
- **Navegação:** Sidebar → Turmas

### Minhas turmas
- **Filtro:** Por série (4, 5, 6, 7, Todas)
- **Grid:** Cards `TeoCom-S4-A` … `TeoCom-S7-C`
- **Card info:** Nome, qtd alunos, qtd atividades
- **Tap:** Detalhe da turma

### Detalhe da turma
- **Header:** Nome da turma, matéria
- **Abas:**
  1. **Atividades** — lista + botão "Nova atividade"
  2. **Alunos** — lista + formulário adicionar por RA
  3. **Métricas** — tabela aluno / atividades feitas / tempo total

### Nova / Editar atividade
- **Campos:**
  - Título (text)
  - Pergunta / enunciado (textarea)
  - Tempo estimado (number, minutos)
  - Tipo de resposta (select)
  - Opções (textarea JSON, visível se multiple_choice)
  - Tipo de complemento (select: none, video, image, link)
  - URL do complemento (text, condicional)
  - Semana de publicação (date)
- **Ações:** Salvar, Cancelar
- **Validação:** Zod schema compartilhado

### Gerenciar alunos (aba)
- **Formulário:** Campo RA + botão "Matricular"
- **Lista:** RA, nome, ações
- **Feedback:** Toast sucesso/erro

## Layout

- Sidebar fixa: Dashboard, Turmas, Sair
- Header com nome do professor
- Responsivo (mobile-friendly)
