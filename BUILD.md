# Build and Test Commands - Quick Reference

## Main Build Script

```powershell
.\build.ps1 [command]
```

## Commands

| Command | Description |
|---------|-------------|
| `.\build.ps1` | Build + Test + Start Docker (default) |
| `.\build.ps1 build` | Build binaries only |
| `.\build.ps1 test` | Run all tests |
| `.\build.ps1 docker` | Start Docker services |
| `.\build.ps1 run` | Build and run services |
| `.\build.ps1 clean` | Remove build artifacts |
| `.\build.ps1 help` | Show help message |

## Examples

### First Time Setup
```powershell
# 1. Copy and edit environment
copy .env.example .env
notepad .env

# 2. Build, test, and setup
.\build.ps1

# 3. Run the application
.\build.ps1 run
```

### Daily Development
```powershell
# After making changes, rebuild and test
.\build.ps1 build
.\build.ps1 test

# Run services
.\build.ps1 run
```

### Testing Only
```powershell
# Run tests
.\build.ps1 test

# Run specific package tests
go test -v ./internal/config
go test -v ./internal/webhook
go test -v ./pkg/llm
```

### Clean Build
```powershell
# Clean everything and rebuild
.\build.ps1 clean
.\build.ps1 build
```

## Docker Services

```powershell
# Start services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Restart services
docker-compose restart
```

## RabbitMQ Management

- URL: http://localhost:15672
- Username: `guest`
- Password: `guest`

## Test Results

Current status: **12/12 tests passing**

```
✓ Config validation (6 tests)
✓ Webhook handling (4 tests)
✓ LLM providers (6 tests)
```

## Build Output

Binaries are created in `bin/` directory:
- `bin/webhook-listener.exe` - HTTP webhook server
- `bin/worker.exe` - Background PR analyzer
