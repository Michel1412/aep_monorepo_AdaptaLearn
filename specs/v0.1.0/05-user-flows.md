# Fluxos de Usuário — v0.1.0

## Fluxo 1: Professor publica atividade

1. Maurilio acessa `/login`
2. Informa email `maurilio@adaptalearn.local` e senha
3. Redirecionado ao dashboard
4. Navega para Turmas → `TeoCom-S7-B`
5. Aba Atividades → "Nova atividade"
6. Preenche template (pergunta, 30 min, essay, vídeo complementar)
7. Salva → POST `/teachers/me/classes/{id}/activities`
8. Atividade aparece na lista da turma

## Fluxo 2: Professor matricula aluno

1. Maurilio em `TeoCom-S7-B` → aba Alunos
2. Informa RA `23159293-2`
3. Clica "Matricular" → POST `/teachers/me/classes/{id}/students`
4. Vinicius aparece na lista

## Fluxo 3: Aluno realiza exercício

1. Vinicius abre app → tela Login
2. Informa RA `23159293-2`
3. Dashboard carrega métricas
4. Navega Minhas turmas → `TeoCom-S7-B`
5. Vê atividades da semana (pendentes e concluídas)
6. Abre atividade pendente
7. Timer inicia automaticamente
8. Lê enunciado, assiste vídeo complementar (link externo)
9. Escreve resposta dissertativa
10. Envia → POST submit
11. Tela de feedback com tempo gasto
12. Dashboard atualizado (horas realizadas incrementadas)

## Fluxo 4: Professor consulta métricas

1. Maurilio em `TeoCom-S7-B` → aba Métricas
2. Vê tabela: Vinicius, 2 atividades, 3600s total
3. Dados vindos de GET `/teachers/me/classes/{id}/metrics`

## Fluxo 5: Aluno vê progresso semanal

1. Vinicius no Dashboard
2. Gráfico mostra TeoCom: 1.5h previstas, 0.5h realizadas
3. Breakdown lista cada exercício com tempo individual
