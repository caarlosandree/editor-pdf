# Editor PDF

Aplicação web full-stack para edição e processamento de documentos PDF, desenvolvida com Go (backend) e Next.js (frontend).

## 📋 Sobre o Projeto

Editor PDF é uma aplicação moderna para edição e manipulação de documentos PDF, oferecendo uma interface intuitiva e recursos avançados de edição. A aplicação permite:

- 📤 **Upload de documentos PDF** - Faça upload de arquivos PDF para o sistema
- 📄 **Visualização de documentos** - Visualize documentos PDF com preview de páginas
- ✏️ **Edição de documentos** - Edite documentos PDF com ferramentas de desenho, texto e imagens
- 🔄 **Processamento de documentos** - Processe documentos PDF com instruções de edição
- 📊 **Gerenciamento de documentos** - Liste, visualize e gerencie seus documentos PDF
- 📝 **Auditoria** - Sistema de logs de auditoria para rastreamento de ações

## 🏗️ Arquitetura

O projeto é dividido em duas partes principais:

- **Backend** (`/backend`) - API REST desenvolvida em Go com Echo framework, seguindo arquitetura limpa (Clean Architecture) com separação em handlers, use cases, repositories e domain
- **Frontend** (`/frontend`) - Interface web desenvolvida em Next.js com React 19, utilizando App Router e Server Components

## 🚀 Tecnologias

### Backend
- **Go 1.25.4** - Linguagem de programação
- **Echo v4** - Framework HTTP
- **sqlx** - Acesso ao banco de dados
- **PostgreSQL 18** - Banco de dados relacional
- **golang-migrate** - Migrations
- **zap** - Logging estruturado
- **viper** - Gerenciamento de configurações
- **validator** (go-playground/validator) - Validação de dados
- **golang-jwt** - Autenticação JWT
- **swaggo/swag** - Documentação Swagger/OpenAPI
- **pdfcpu** - Processamento de PDFs
- **unipdf** - Manipulação avançada de PDFs

### Frontend
- **Next.js 16.1.1** - Framework React full-stack
- **React 19.2.3** - Framework UI
- **TypeScript 5** - Tipagem estática
- **Tailwind CSS 4** - Estilização utilitária
- **shadcn/ui** - Componentes de UI (baseado em Radix UI)
- **React Hook Form** - Gerenciamento de formulários
- **Zod 4** - Validação de schemas
- **TanStack Query 5** - Gerenciamento de estado servidor e cache
- **Axios** - Cliente HTTP
- **Zustand** - Gerenciamento de estado global
- **Nuqs** - Sincronização de estado com URL (Type-safe Search Params)
- **Recharts** - Visualização de dados e gráficos
- **date-fns** - Manipulação de datas
- **Lucide React** - Ícones
- **@uidotdev/usehooks** - Hooks customizados

## 📁 Estrutura do Projeto

```
editor-pdf/
├── backend/                    # Backend Go
│   ├── cmd/
│   │   └── server/            # Ponto de entrada da aplicação
│   │       ├── main.go        # Arquivo principal
│   │       └── docs/          # Documentação Swagger gerada
│   ├── internal/              # Código interno da aplicação
│   │   ├── config/           # Configurações
│   │   ├── domain/           # Entidades e interfaces de domínio
│   │   ├── dto/              # Data Transfer Objects
│   │   ├── handler/          # Handlers HTTP (controllers)
│   │   ├── infrastructure/   # Implementações de infraestrutura
│   │   │   ├── pdf/          # Processadores de PDF
│   │   │   └── storage/       # Armazenamento de arquivos
│   │   ├── middleware/        # Middlewares HTTP
│   │   ├── model/            # Modelos de dados
│   │   ├── repository/       # Repositórios (acesso ao banco)
│   │   ├── usecase/          # Casos de uso (lógica de negócio)
│   │   ├── util/             # Utilitários
│   │   └── validator/        # Validadores customizados
│   ├── pkg/                   # Código reutilizável
│   │   ├── logger/           # Logger customizado
│   │   └── response/         # Helpers de resposta HTTP
│   ├── migrations/            # Scripts de migration do banco
│   ├── storage/               # Armazenamento de arquivos PDF
│   ├── tests/                 # Testes
│   ├── go.mod                 # Dependências Go
│   └── Makefile               # Comandos úteis
├── frontend/                   # Frontend Next.js
│   ├── app/                   # App Router (Next.js 13+)
│   │   ├── documents/         # Rotas de documentos
│   │   │   └── [id]/          # Rota dinâmica por ID
│   │   │       └── edit/      # Página de edição
│   │   ├── layout.tsx         # Layout raiz
│   │   ├── page.tsx           # Página inicial (Dashboard)
│   │   └── globals.css        # Estilos globais
│   ├── components/            # Componentes React
│   │   ├── documents/        # Componentes de documentos
│   │   ├── layout/           # Componentes de layout
│   │   ├── pdf-editor/       # Componentes do editor PDF
│   │   └── ui/               # Componentes shadcn/ui
│   ├── hooks/                 # Custom hooks
│   │   ├── mutations/        # Hooks de mutations (TanStack Query)
│   │   └── queries/          # Hooks de queries (TanStack Query)
│   ├── lib/                   # Bibliotecas e configurações
│   │   ├── axios.ts          # Cliente Axios configurado
│   │   └── utils.ts          # Utilitários (cn, etc.)
│   ├── providers/             # Providers React
│   │   ├── QueryProvider.tsx # Provider do TanStack Query
│   │   └── NuqsAdapter.tsx   # Adapter do Nuqs
│   ├── schemas/               # Schemas Zod
│   ├── services/              # Serviços e APIs
│   ├── stores/                # Stores Zustand
│   ├── types/                 # Tipos TypeScript
│   ├── utils/                 # Utilitários
│   ├── package.json          # Dependências Node.js
│   └── tsconfig.json         # Configuração TypeScript
└── .cursor/                   # Regras de desenvolvimento
    └── rules/                 # Regras organizadas por módulos
```

## 🛠️ Pré-requisitos

- **Go 1.25.4** ou superior
- **Node.js 20** ou superior
- **PostgreSQL 18** ou superior
- **npm** ou **yarn**

## ⚙️ Configuração

### Backend

1. Navegue até a pasta do backend:
```bash
cd backend
```

2. Instale as dependências:
```bash
go mod download
```

3. Configure as variáveis de ambiente criando um arquivo `.env.local` ou `.env`:
```env
SERVER_PORT=8080
SERVER_HOST=localhost

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=editor_pdf
DB_SSLMODE=disable

JWT_SECRET=your-secret-key-here-change-in-production
JWT_EXPIRATION=24h

CORS_ALLOWED_ORIGINS=http://localhost:3000

STORAGE_PATH=./storage
STORAGE_MAX_UPLOAD_SIZE=104857600

ENV=development
```

**Nota**: O arquivo `.env.local` tem prioridade sobre `.env`. Variáveis de ambiente também podem sobrescrever valores dos arquivos.

4. Execute as migrations:
```bash
make migrate-up
```

5. Execute o servidor:
```bash
make run
```

O servidor estará disponível em `http://localhost:8080`

### Frontend

1. Navegue até a pasta do frontend:
```bash
cd frontend
```

2. Instale as dependências:
```bash
npm install
```

3. Configure as variáveis de ambiente criando um arquivo `.env.local`:
```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NODE_ENV=development
```

4. Execute o servidor de desenvolvimento:
```bash
npm run dev
```

A aplicação estará disponível em `http://localhost:3000`

## 📚 Documentação

### API REST

A documentação Swagger da API está disponível em:
- **Desenvolvimento**: `http://localhost:8080/swagger/index.html`

### Endpoints Disponíveis

#### Documentos
- `POST /api/v1/documents` - Upload de documento PDF
- `GET /api/v1/documents` - Lista todos os documentos
- `GET /api/v1/documents/:id` - Obtém um documento específico
- `POST /api/v1/documents/:id/process` - Processa um documento com instruções de edição
- `GET /api/v1/documents/:id/preview/:page` - Gera preview de uma página do documento
- `DELETE /api/v1/documents/:id` - Remove um documento

#### Health Check
- `GET /health` - Verifica o status do servidor

### Regras de Desenvolvimento

O projeto possui regras de desenvolvimento organizadas em módulos:

- **Backend**: Regras em `.cursor/rules/backend/` e `.cursor/rules/postgresql.mdc`
- **Frontend**: Regras em `.cursor/rules/frontend/`
- **Commits**: Padrões em `.cursor/rules/commit.mdc`

## 🧪 Testes

### Backend
```bash
cd backend

# Executar todos os testes
make test

# Executar testes com cobertura
make test-coverage
```

### Frontend
```bash
cd frontend

# Executar testes (quando implementado)
npm run test
```

## 📝 Scripts Úteis

### Backend
```bash
cd backend

# Ver todos os comandos disponíveis
make help

# Instalar dependências
make deps

# Compilar
make build

# Executar (instala deps, executa migrations e compila antes)
make run

# Testes
make test
make test-coverage

# Linter
make lint

# Formatar código
make format

# Migrations
make migrate-up              # Executa todas as migrations pendentes
make migrate-down            # Reverte a última migration
make migrate-create NAME=nome_da_migration  # Cria nova migration

# Swagger
make swagger                 # Gera documentação Swagger
make swagger-clean           # Limpa arquivos gerados do Swagger

# Limpeza
make clean                   # Remove arquivos gerados
```

### Frontend
```bash
cd frontend

# Instalar dependências
npm install

# Desenvolvimento
npm run dev

# Build para produção
npm run build

# Executar build de produção
npm start

# Linter
npm run lint

# Typecheck (quando implementado)
npm run typecheck
```

## 🤝 Contribuindo

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças seguindo os padrões de commit
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

### Padrões de Commit

Este projeto segue o padrão [Conventional Commits](https://www.conventionalcommits.org/) com categorias:

**Tipos de Commit:**
- `feat:` - Nova funcionalidade
- `fix:` - Correção de bug
- `docs:` - Documentação
- `style:` - Formatação (não altera lógica)
- `refactor:` - Refatoração (sem mudança de funcionalidade)
- `test:` - Testes
- `chore:` - Tarefas de manutenção

**Organização por Categoria:**
- `Backend:` - Mudanças no backend
- `Frontend:` - Mudanças no frontend
- `Database:` - Mudanças no banco de dados

**Exemplos:**
```
feat(Backend): adiciona endpoint de upload de documentos
fix(Frontend): corrige validação de formulário de upload
refactor(Backend): reorganiza estrutura de use cases
docs: atualiza README com novas funcionalidades
```

Consulte `.cursor/rules/commit.mdc` para mais detalhes sobre os padrões de commit.

## 📄 Licença

Este projeto está sob a licença MIT.

## 👥 Autores

- **Carlos André Sabino** - Desenvolvimento inicial

## 🔒 Segurança

- **JWT**: Autenticação baseada em tokens JWT
- **CORS**: Configuração de CORS para controle de origens permitidas
- **Validação**: Validação de dados de entrada no backend e frontend
- **SQL Injection**: Proteção através de prepared statements (sqlx)
- **Headers de Segurança**: Middleware de segurança com headers HTTP apropriados
- **Auditoria**: Sistema de logs de auditoria para rastreamento de ações

## 📊 Funcionalidades

### Backend
- ✅ Upload e armazenamento de documentos PDF
- ✅ Processamento de PDFs com pdfcpu e unipdf
- ✅ Geração de preview de páginas PDF
- ✅ Sistema de auditoria (audit logs)
- ✅ API REST versionada (`/api/v1/`)
- ✅ Documentação Swagger/OpenAPI
- ✅ Logging estruturado com zap
- ✅ Migrations com golang-migrate
- ✅ Graceful shutdown

### Frontend
- ✅ Dashboard para gerenciamento de documentos
- ✅ Upload de documentos PDF
- ✅ Listagem de documentos
- ✅ Editor de PDF com ferramentas de desenho, texto e imagens
- ✅ Preview de documentos PDF
- ✅ Interface responsiva com Tailwind CSS
- ✅ Componentes acessíveis com shadcn/ui
- ✅ Gerenciamento de estado com TanStack Query e Zustand
- ✅ Validação de formulários com React Hook Form e Zod

## 🗄️ Banco de Dados

O projeto utiliza PostgreSQL 18 com as seguintes tabelas:

- **users** - Usuários do sistema
- **documents** - Documentos PDF armazenados
- **audit_logs** - Logs de auditoria

As migrations estão em `backend/migrations/` e podem ser executadas com `make migrate-up`.

## 🙏 Agradecimentos

- Comunidade Go
- Comunidade React/Next.js
- shadcn/ui por componentes incríveis e acessíveis
- pdfcpu e unipdf por bibliotecas de processamento de PDF
