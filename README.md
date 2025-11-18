# Code Review AI Tool

An AI-powered code review tool that integrates with GitHub to provide intelligent feedback on pull requests.

## Architecture

- **Webhook Listener**: HTTP server receiving GitHub webhook events
- **Message Queue**: RabbitMQ for async processing
- **Worker Pool**: Go workers analyzing code with AI
- **LLM Integration**: Configurable AI providers (OpenAI, Anthropic, etc.)

## Prerequisites

- Go 1.21+
- Docker & Docker Compose
- GitHub account and repository

## Quick Start

### 1. Setup Environment

```powershell
# Copy environment template
copy .env.example .env
# Edit .env with your API keys and configuration
```

### 2. Build and Test Everything

```powershell
# Build, test, and start Docker services
.\build.ps1

# Or use individual commands:
.\build.ps1 build   # Just build binaries
.\build.ps1 test    # Just run tests
.\build.ps1 docker  # Just start Docker services
.\build.ps1 run     # Build and run services
.\build.ps1 clean   # Clean build artifacts
```

### 3. Run Services

```powershell
# Start all services (webhook listener + worker)
.\build.ps1 run
```

### 4. Configure GitHub Webhook

1. Go to your repository Settings → Webhooks → Add webhook
2. Payload URL: `http://your-server:8080/webhook/github`
3. Content type: `application/json`
4. Secret: Use the value from `GITHUB_WEBHOOK_SECRET` in your `.env`
5. Events: Select "Pull requests" and "Pull request reviews"

## Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `GITHUB_WEBHOOK_SECRET` | Secret for validating GitHub webhooks | Yes |
| `GITHUB_TOKEN` | GitHub personal access token | Yes |
| `LLM_PROVIDER` | AI provider: `openai`, `anthropic`, `ollama` | Yes |
| `OPENAI_API_KEY` | OpenAI API key (if using OpenAI) | Conditional |
| `ANTHROPIC_API_KEY` | Anthropic API key (if using Anthropic) | Conditional |
| `OLLAMA_URL` | Ollama server URL (if using Ollama) | Conditional |
| `RABBITMQ_URL` | RabbitMQ connection URL | Yes |
| `POSTGRES_URL` | PostgreSQL connection URL | Yes |

## Development

### Building

```powershell
# Build both binaries
.\build.ps1 build

# Or manually:
go build -o bin/webhook-listener.exe ./cmd/webhook-listener
go build -o bin/worker.exe ./cmd/worker
```

### Testing

```powershell
# Run all tests
.\build.ps1 test

# Or manually:
go test -v ./...
go test -cover ./...
```

See [`TESTING.md`](TESTING.md) for detailed testing guide.

**Current Test Coverage:**
- ✅ 12/12 tests passing
- ✅ Config validation tests
- ✅ Webhook signature verification tests
- ✅ LLM provider factory tests
```

## Deployment

### Docker

```bash
# Build images
docker build -t codereview-webhook -f Dockerfile.webhook .
docker build -t codereview-worker -f Dockerfile.worker .

# Run with docker-compose
docker-compose up
```

## License

MIT
