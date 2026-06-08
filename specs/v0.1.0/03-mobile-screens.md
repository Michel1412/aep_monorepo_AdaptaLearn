# Telas Mobile (Aluno) — v0.1.0

App: `apps/mobile` — Expo Router

## Rotas

| Rota | Tela | Descrição |
|------|------|-----------|
| `/(auth)/login` | Login | Entrada por RA |
| `/(student)/` | Dashboard | Métricas e gráfico |
| `/(student)/classes` | Minhas turmas | Lista de turmas matriculadas |
| `/(student)/classes/[id]` | Atividades da semana | Lista por turma |
| `/(student)/activities/[id]` | Realizar questão | Timer + resposta |
| `/(student)/activities/[id]/result` | Feedback | Confirmação pós-envio |

## Detalhamento por tela

### Login
- **Campos:** RA (texto, máscara livre)
- **Ações:** Botão "Entrar" → POST `/auth/student`
- **Estados:** loading, erro (RA inválido), sucesso → redirect dashboard
- **Validação:** RA não vazio, mínimo 5 caracteres

### Dashboard
- **Header:** Saudação com nome do aluno
- **Cards:**
  - Atividades concluídas: `X / Y` na semana
  - Horas de estudo: previstas vs realizadas (número + barra)
- **Gráfico:** Barras horizontais — previsto (cinza) vs realizado (azul) por matéria
- **Lista:** Breakdown por exercício (título, tempo estimado, tempo gasto)
- **Navegação:** Botão/link "Minhas turmas"
- **Pull-to-refresh:** Recarrega GET `/students/{ra}/dashboard`

### Minhas turmas
- **Lista:** Cards com nome da turma, matéria, progresso semanal
- **Tap:** Navega para atividades da turma

### Atividades da semana
- **Lista:** Cards com título, tempo estimado, badge status (pendente/concluída)
- **Tap em pendente:** Navega para realizar questão
- **Tap em concluída:** Navega para feedback (somente leitura)

### Realizar questão
- **Conteúdo:** Título, enunciado, complemento (link/vídeo/imagem)
- **Timer:** Cronômetro visível, inicia ao montar tela (POST sessions/start)
- **Input:** Conforme `answer_type`:
  - `multiple_choice`: radio buttons
  - `true_false`: dois botões Verdadeiro/Falso
  - `short_text`: input texto
  - `essay`: textarea
- **Ações:** Botão "Enviar resposta"
- **Validação:** Resposta obrigatória

### Feedback pós-envio
- **Conteúdo:** Confirmação, tempo gasto formatado
- **Ações:** "Voltar ao dashboard" ou "Próxima atividade"

## Componentes reutilizáveis

- `StatCard` — card de métrica
- `HoursChart` — gráfico de horas (Victory Native)
- `ActivityCard` — card de atividade na lista
- `QuestionInput` — input dinâmico por tipo
- `StudyTimer` — display do cronômetro
