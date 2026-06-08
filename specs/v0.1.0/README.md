# AdaptaLearn — Spec v0.1.0

## Escopo

Primeira projeção do ecossistema AdaptaLearn: protótipo funcional com backend Go, portal web para professores (React) e app mobile para alunos (Expo). Todos os clientes consomem a API real com dados seed; mocks são usados apenas em testes unitários.

## Objetivo

Permitir que professores publiquem atividades extras por turma e que alunos realizem exercícios com rastreamento de tempo de estudo, visualizando métricas no dashboard.

## Premissas

- Login de aluno apenas por RA (Registro Acadêmico)
- Login de professor por email e senha
- Nomenclatura de turma: `{MateriaCode}-S{serie}-{Turma}` (ex.: `TeoCom-S7-B`)
- Complementos de questão via URL (vídeo, imagem ou link) — sem upload na v0.1.0
- Sem algoritmo de IA / trilhas adaptativas nesta versão

## Critérios de aceite

- [ ] Professor Maurilio consegue logar e ver suas 12 turmas de TeoCom
- [ ] Maurilio consegue criar atividade para `TeoCom-S7-B`
- [ ] Aluno Vinicius (RA `23159293-2`) consegue logar e ver dashboard
- [ ] Vinicius vê atividades da semana em `TeoCom-S7-B`
- [ ] Vinicius inicia sessão, responde questão e tempo é contabilizado
- [ ] Dashboard do aluno reflete horas previstas vs realizadas
- [ ] Testes passam em backend, api-client, web e mobile

## Documentos desta versão

| Arquivo | Descrição |
|---------|-----------|
| [01-domain-model.md](./01-domain-model.md) | Entidades e relacionamentos |
| [02-api-contracts.md](./02-api-contracts.md) | Contratos REST |
| [03-mobile-screens.md](./03-mobile-screens.md) | Telas do app aluno |
| [04-web-screens.md](./04-web-screens.md) | Telas do portal professor |
| [05-user-flows.md](./05-user-flows.md) | Fluxos ponta a ponta |
| [06-seed-data.md](./06-seed-data.md) | Dados de demonstração |
| [07-test-matrix.md](./07-test-matrix.md) | Matriz de testes |
