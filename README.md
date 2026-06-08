# AdaptaLearn

Plataforma educacional adaptativa para apoio a alunos com dificuldades e exercícios extras, conectando professores (web) e alunos (mobile) via API Go.

## Estrutura do monorepo

```
AdaptaLearn/
├── specs/v0.1.0/     # Especificações versionadas
├── apps/
│   ├── api/          # Backend Go (Gin + SQLite)
│   ├── web/          # Portal do professor (React + Vite)
│   └── mobile/       # App do aluno (Expo SDK 54)
└── packages/
    ├── shared-types/ # Tipos e schemas Zod
    └── api-client/   # Cliente HTTP tipado
```

## Pré-requisitos

- Node.js 20+ (recomendado via [nvm](https://github.com/nvm-sh/nvm))
- pnpm 9+ (ou use `npx pnpm@9.15.0` — ver abaixo)
- Go 1.22+

## Instalar o pnpm

O projeto declara `packageManager: pnpm@9.15.0`. Se `pnpm` não for reconhecido no terminal, escolha **uma** opção:

### Opção 1 — Corepack (recomendado, vem com o Node)

```bash
corepack enable
corepack prepare pnpm@9.15.0 --activate
pnpm --version   # deve mostrar 9.15.0
```

Se der erro de permissão com nvm, use a opção 2 ou 3.

### Opção 2 — npm global

```bash
npm install -g pnpm@9.15.0
pnpm --version
```

### Opção 3 — Sem instalar (sempre funciona)

Use `npx` no lugar de `pnpm`:

```bash
npx pnpm@9.15.0 install
npx pnpm@9.15.0 --filter @adaptalearn/web dev
```

Ou os scripts npm do projeto (já configurados com `npx`):

```bash
npm run install:deps
npm run dev:api
npm run dev:web
npm run dev:mobile
```

## Setup

```bash
# Instalar dependências
pnpm install
# ou: npm run install:deps

# Build dos pacotes compartilhados
pnpm --filter @adaptalearn/shared-types build
pnpm --filter @adaptalearn/api-client build

# Subir API (seed automático na primeira execução)
npm run dev:api
# ou: cd apps/api && go run ./cmd/server
```

Em terminais separados:

```bash
pnpm dev:web      # http://localhost:3000
pnpm dev:mobile   # Expo — escaneie QR ou pressione 'w' para web
```

## Credenciais de demonstração (seed v0.1.0)

### Professor
- **Email:** maurilio@adaptalearn.local
- **Senha:** senha123
- **Turmas:** 12 turmas TeoCom (S4–S7, A/B/C)

### Aluno
- **RA:** 23159293-2 (Vinicius — matriculado em TeoCom-S7-B)
- **RA:** 23034350-2 (João Pedro — TeoCom-S7-A)
- **RA:** 23220783-2 (Michel — TeoCom-S6-C)

## Variáveis de ambiente

| App | Variável | Default |
|-----|----------|---------|
| API | `PORT` | 8080 |
| API | `JWT_SECRET` | adaptalearn-dev-secret |
| API | `DATABASE_PATH` | apps/api/data/adaptalearn.db |
| Web | `VITE_API_URL` | http://localhost:8080 |
| Mobile | `EXPO_PUBLIC_API_URL` | http://localhost:8080 |

## Testes

```bash
# Todos (via Turborepo)
pnpm test

# Backend Go
cd apps/api && go test ./...

# Pacotes individuais
pnpm --filter @adaptalearn/shared-types test
pnpm --filter @adaptalearn/api-client test
pnpm --filter @adaptalearn/web test
pnpm --filter @adaptalearn/mobile test
```

## Fluxo de teste manual

1. Inicie a API e acesse o portal web em http://localhost:3000
2. Login como Maurilio → Turmas → TeoCom-S7-B → Nova atividade
3. Abra o app mobile, login com RA `23159293-2`
4. Dashboard → Turmas → TeoCom-S7-B → Atividade → Responder
5. Verifique métricas atualizadas no dashboard e no portal (aba Métricas)

## Specs

Documentação completa da v0.1.0 em [`specs/v0.1.0/`](specs/v0.1.0/README.md).

## Equipe

- João Pedro Souza Peixoto Saraiva — 23034350-2
- Michel Bocchi Junior — 23220783-2
- Vinicius Reginaldo Ferrarini — 23159293-2
