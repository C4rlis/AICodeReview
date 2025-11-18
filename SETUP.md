# Setup Instructions

## First Time Setup

### Option 1: Using PowerShell (Recommended)

If you get a "scripts is disabled" error, run this **once** as Administrator:

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

Then use the PowerShell build script:

```powershell
.\build.ps1
```

### Option 2: Using Batch File (No Admin Required)

Use the batch file alternative that works without changing execution policies:

```batch
build.bat
```

Both scripts support the same commands:
- `build` - Build binaries only
- `test` - Run tests only  
- `docker` - Start Docker services
- `run` - Build and run services
- `clean` - Clean build artifacts
- `all` - Do everything (default)

## Quick Commands

```powershell
# PowerShell
.\build.ps1           # Build + Test + Docker
.\build.ps1 build     # Just build
.\build.ps1 test      # Just test
.\build.ps1 run       # Build and run

# Batch (alternative)
build.bat            # Build + Test + Docker
build.bat build      # Just build
build.bat test       # Just test
build.bat run        # Build and run
```

## Environment Setup

1. **Copy environment template:**
   ```powershell
   copy .env.example .env
   ```

2. **Edit `.env` with your credentials:**
   - Add GitHub token
   - Add LLM API key (OpenAI/Anthropic) or use Ollama (free)
   - Keep other defaults or customize

3. **Run the build:**
   ```powershell
   .\build.ps1
   # or
   build.bat
   ```

## Running Services

**Automatic (Recommended):**
```powershell
.\build.ps1 run
# or
build.bat run
```

This will:
1. Build both binaries
2. Start Docker services
3. Launch webhook listener and worker

**Manual:**
```powershell
# Terminal 1 - Webhook Listener
go run .\cmd\webhook-listener\main.go

# Terminal 2 - Worker  
go run .\cmd\worker\main.go
```

## Verifying Setup

1. **Check build status:**
   ```powershell
   .\build.ps1 build
   ```
   Should create `bin/webhook-listener.exe` and `bin/worker.exe`

2. **Run tests:**
   ```powershell
   .\build.ps1 test
   ```
   Should show: 12/12 tests passing

3. **Check Docker services:**
   ```powershell
   docker ps
   ```
   Should show RabbitMQ and PostgreSQL running

4. **Access management UIs:**
   - RabbitMQ: http://localhost:15672 (guest/guest)
   - Webhook endpoint: http://localhost:8080/webhook/github

## Troubleshooting

### PowerShell Execution Policy Error

**Error:** `running scripts is disabled on this system`

**Solution 1 (Recommended):**
Run as Administrator:
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

**Solution 2:**
Use the batch file instead:
```batch
build.bat
```

**Solution 3 (One-time bypass):**
```powershell
powershell -ExecutionPolicy Bypass -File .\build.ps1
```

### Docker Not Running

**Error:** `Docker is not running`

**Solution:**
1. Start Docker Desktop
2. Wait for it to fully initialize
3. Run the build script again

### Tests Failing

**If tests fail:**
1. Check Go version: `go version` (need 1.21+)
2. Run `go mod tidy` to fix dependencies
3. Check specific failing tests for details

### Services Won't Start

**Check:**
1. `.env` file exists and has correct values
2. Docker services are running: `docker ps`
3. Ports 8080, 5672, 5432, 15672 are not in use
4. API keys are valid

## Next Steps

Once everything is running:

1. **Configure GitHub webhook** (see [QUICKSTART.md](QUICKSTART.md))
2. **Create a test PR** in your repository
3. **Watch the logs** for AI analysis
4. **See the review** posted on your PR!

## Getting Help

- Build issues: Check [BUILD.md](BUILD.md)
- Testing: See [TESTING.md](TESTING.md)
- General setup: Read [QUICKSTART.md](QUICKSTART.md)
- Full docs: See [README.md](README.md)
