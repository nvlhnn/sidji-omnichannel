# 🗨️ Sidji-Omnichannel - Omnichannel Chat Platform

A unified inbox platform for managing customer conversations across WhatsApp and Instagram.

## ✨ Features

- **Unified Inbox**: View and respond to messages from multiple channels in one place
- **WhatsApp Business Integration**: Connect your WhatsApp Business account via Meta Cloud API
- **Instagram DM Integration**: Manage Instagram Direct Messages
- **Team Collaboration**: Multiple agents with role-based access (Admin, Supervisor, Agent)
- **Real-time Updates**: WebSocket-powered live message updates
- **Conversation Management**: Assign, label, and track conversation status
- **Canned Responses**: Quick reply templates for common questions
- **AI Assistant**: Smart auto-replies with RAG (Retrieval-Augmented Generation)
- **AI Modes**: Support for Manual, Fully Automatic, and Hybrid (Human-First Transition) modes
- **Knowledge Base**: Train your AI with business-specific data for accurate responses

## 🏗️ Tech Stack

- **Backend**: Go 1.21+ with Gin framework
- **Architecture**: **Hexagonal Architecture (Ports and Adapters)**
- **Frontend**: Next.js 14 with TypeScript (separate project)
- **Database**: PostgreSQL 15 + pgvector (for knowledge base)
- **Cache/Realtime**: Redis
- **WebSocket**: Gorilla WebSocket
- **Deployment**: Docker + AWS

## 📁 Project Structure

```
Sidji-Omnichannel/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── domain/
│   │   └── ports/               # Port definitions (Interfaces)
│   │       ├── repository/      # Outbound ports (Driven Ports)
│   │       └── service/         # Inbound ports (Driving Ports)
│   ├── adapters/                # Concrete implementations (Driven Adapters)
│   │   └── db/
│   │       └── postgres/        # PostgreSQL repository implementations
│   ├── services/                # Business logic & Use Cases (Implementing Inbound Ports)
│   ├── handlers/                # HTTP handlers (Driving Adapters)
│   ├── models/                  # Domain entities and data models
│   ├── middleware/              # Auth, CORS, etc.
│   ├── integrations/            # External API integrations (Meta, etc.)
│   ├── ai/                      # AI Provider implementations (Gemini, OpenAI)
│   ├── websocket/               # Real-time hub
│   ├── config/                  # Configuration
│   └── testutil/                # Testing utilities
├── migrations/                  # SQL migrations
├── web/                         # Frontend (Next.js)
├── docker-compose.yml
├── Dockerfile
└── go.mod
```

## 📐 Architecture

This project follows the **Hexagonal Architecture** (also known as Ports and Adapters) to ensure high maintainability, testability, and decoupling from external dependencies:

1.  **Domain (Core)**: Contains the business entities (`models/`) and Port definitions (`domain/ports/`).
2.  **Ports**: Interfaces that define how the core interacts with the outside world.
    -   **Inbound Ports (Service)**: Define the use cases available to be triggered (e.g., `TeamService` interface).
    -   **Outbound Ports (Repository)**: Define how the core fetches/saves data (e.g., `TeamRepository` interface).
3.  **Adapters**:
    -   **Driving Adapters (Handlers)**: The HTTP handlers that call the Service ports.
    -   **Driven Adapters (Postgres)**: The concrete implementations of Repository ports.

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Node.js 18+ (for frontend)

### 1. Clone and Setup

```bash
git clone https://github.com/yourusername/sidji-omnichannel.git
cd sidji-omnichannel
cp .env.example .env
```

### 2. Start Dependencies

```bash
# Start PostgreSQL and Redis
docker-compose up -d postgres redis
```

### 3. Run Migrations

```bash
# Using psql or your preferred tool
psql -h localhost -p 5433 -U sidji -d sidji_test -f migrations/001_initial_schema.up.sql
```

### 4. Run the Server

```bash
# Install dependencies
go mod tidy

# Run the server
go run cmd/server/main.go
```

The API will be available at `http://localhost:8080`

### 5. Run Frontend (separate terminal)

```bash
cd web
npm install
npm run dev
```

The frontend will be available at `http://localhost:3000`

## 📡 API Endpoints

### Authentication
- `POST /api/auth/register` - Register new organization & admin
- `POST /api/auth/login` - Login
- `GET /api/auth/me` - Get current user

### Conversations
- `GET /api/conversations` - List conversations
- `GET /api/conversations/:id` - Get conversation details
- `POST /api/conversations/:id/assign` - Assign to agent
- `PATCH /api/conversations/:id/status` - Update status
- `GET /api/conversations/:id/messages` - Get messages
- `POST /api/conversations/:id/messages` - Send message

### Channels
- `GET /api/channels` - List connected channels
- `POST /api/channels/discover/meta` - Discover available Meta accounts
- `POST /api/channels/whatsapp/connect` - Connect selected WhatsApp account
- `POST /api/channels/instagram/connect` - Connect selected Instagram account
- `POST /api/channels/facebook/connect` - Connect selected Facebook page
- `DELETE /api/channels/:id` - Disconnect channel

### Webhooks
- `GET /api/webhooks` - Meta webhook verification
- `POST /api/webhooks` - Receive Meta (WhatsApp/Instagram) webhooks

### WebSocket
- `GET /api/ws` - WebSocket connection for real-time updates

### Team Management
- `GET /api/team/members` - List team members
- `POST /api/team/members` - Invite new member
- `PATCH /api/team/members/:id` - Update member

### Canned Responses
- `GET /api/canned-responses` - List canned responses
- `POST /api/canned-responses` - Create new response
- `PUT /api/canned-responses/:id` - Update response
- `DELETE /api/canned-responses/:id` - Delete response

### Labels
- `GET /api/labels` - List all labels
- `POST /api/labels` - Create new label
- `PUT /api/labels/:id` - Update label
- `DELETE /api/labels/:id` - Delete label

### AI Configuration
- `GET /api/channels/:id/ai` - Get AI configuration for a channel
- `PUT /api/channels/:id/ai` - Update AI configuration
- `GET /api/channels/:id/ai/knowledge` - List knowledge base items
- `POST /api/channels/:id/ai/knowledge` - Add new knowledge item
- `PUT /api/channels/:id/ai/knowledge/:kid` - Update knowledge item
- `DELETE /api/channels/:id/ai/knowledge/:kid` - Delete knowledge item
- `POST /api/channels/:id/ai/test` - Test AI reply generation

## 🔧 Meta API Setup

To receive messages from WhatsApp and Instagram, you need to:

1. Create a Meta Business Account
2. Create a Facebook App with WhatsApp and Instagram products
3. Configure webhooks to point to your server
4. Set the verify token in your `.env` file

See [META_SETUP.md](docs/META_SETUP.md) for detailed instructions.

## 📦 Deployment

### Docker

```bash
# Build and run everything
docker-compose up -d
```

### AWS

See [DEPLOYMENT.md](docs/DEPLOYMENT.md) for AWS deployment guide.

## 🔒 Environment Variables

| Variable | Description |
|----------|-------------|
| `APP_ENV` | Environment (development/production) |
| `APP_PORT` | Server port (default: 8080) |
| `APP_SECRET` | JWT secret key |
| `DB_HOST` | PostgreSQL host |
| `DB_PORT` | PostgreSQL port |
| `DB_USER` | PostgreSQL user |
| `DB_PASSWORD` | PostgreSQL password |
| `DB_NAME` | PostgreSQL database name |
| `REDIS_HOST` | Redis host |
| `REDIS_PORT` | Redis port |
| `META_APP_ID` | Meta App ID (from App Dashboard) |
| `META_APP_SECRET` | Meta App Secret (from App Dashboard) |
| `META_VERIFY_TOKEN` | Webhook verification token |
| `GEMINI_API_KEY` | Google Gemini API Key |
| `AI_PROVIDER` | AI Provider (gemini or openai) |
| `OPENAI_API_KEY` | OpenAI API Key |

## 📝 License

MIT License
