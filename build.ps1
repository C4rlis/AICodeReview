# Build and Test Script for Code Review AI Tool
# Usage: .\build.ps1 [command]
# Commands: build, test, run, docker, clean, all

param(
    [Parameter(Position=0)]
    [string]$Command = "all"
)

$ErrorActionPreference = "Stop"

# Colors for output
function Write-Info { Write-Host $args -ForegroundColor Cyan }
function Write-Success { Write-Host $args -ForegroundColor Green }
function Write-Error { Write-Host $args -ForegroundColor Red }
function Write-Warning { Write-Host $args -ForegroundColor Yellow }

# Banner
Write-Host "`n==================================================" -ForegroundColor Magenta
Write-Host "  Code Review AI - Build & Test System" -ForegroundColor Magenta
Write-Host "==================================================" -ForegroundColor Magenta

function Build-Project {
    Write-Info "`n[1/3] Building binaries..."
    
    try {
        # Create bin directory
        if (-not (Test-Path "bin")) {
            New-Item -ItemType Directory -Path "bin" | Out-Null
        }

        # Build webhook listener
        Write-Info "  Building webhook-listener..."
        go build -o bin/webhook-listener.exe ./cmd/webhook-listener
        Write-Success "  ✓ webhook-listener.exe"

        # Build worker
        Write-Info "  Building worker..."
        go build -o bin/worker.exe ./cmd/worker
        Write-Success "  ✓ worker.exe"

        Write-Success "`n✓ Build completed successfully!`n"
        return $true
    }
    catch {
        Write-Error "`n✗ Build failed: $_`n"
        return $false
    }
}

function Run-Tests {
    Write-Info "`n[2/3] Running tests..."
    
    try {
        # Run all tests with coverage
        Write-Info "  Running unit tests..."
        go test -v -cover ./...
        
        if ($LASTEXITCODE -eq 0) {
            Write-Success "`n✓ All tests passed!`n"
            return $true
        } else {
            Write-Error "`n✗ Tests failed!`n"
            return $false
        }
    }
    catch {
        Write-Error "`n✗ Test execution failed: $_`n"
        return $false
    }
}

function Start-Docker {
    Write-Info "`n[3/3] Starting Docker services..."
    
    try {
        # Check if Docker is running
        docker info 2>&1 | Out-Null
        if ($LASTEXITCODE -ne 0) {
            Write-Warning "  Docker is not running. Please start Docker Desktop."
            return $false
        }

        Write-Info "  Starting RabbitMQ and PostgreSQL..."
        docker-compose up -d
        
        Write-Info "  Waiting for services to be ready..."
        Start-Sleep -Seconds 5
        
        Write-Success "  ✓ RabbitMQ: http://localhost:15672 (guest/guest)"
        Write-Success "  ✓ PostgreSQL: localhost:5432"
        Write-Success "`n✓ Docker services started!`n"
        return $true
    }
    catch {
        Write-Error "`n✗ Docker startup failed: $_`n"
        return $false
    }
}

function Run-Services {
    Write-Info "`n[4/4] Starting services..."
    Write-Warning "`nThis will start both services in the foreground."
    Write-Warning "Press Ctrl+C to stop.`n"
    
    # Check if .env exists
    if (-not (Test-Path ".env")) {
        Write-Warning ".env file not found. Creating from template..."
        Copy-Item ".env.example" ".env"
        Write-Warning "Please edit .env with your API keys before running!"
        return
    }

    Write-Info "Starting webhook listener and worker..."
    Write-Info "Webhook: http://localhost:8080/webhook/github`n"

    # Start both services in background jobs
    $webhook = Start-Job -ScriptBlock { 
        Set-Location $args[0]
        go run ./cmd/webhook-listener/main.go 
    } -ArgumentList $PWD

    Start-Sleep -Seconds 2

    $worker = Start-Job -ScriptBlock { 
        Set-Location $args[0]
        go run ./cmd/worker/main.go 
    } -ArgumentList $PWD

    Write-Success "✓ Services started!`n"
    Write-Info "Logs:`n"

    try {
        while ($true) {
            Receive-Job -Job $webhook
            Receive-Job -Job $worker
            Start-Sleep -Milliseconds 500
        }
    }
    finally {
        Write-Info "`nStopping services..."
        Stop-Job -Job $webhook, $worker
        Remove-Job -Job $webhook, $worker
        Write-Success "✓ Services stopped"
    }
}

function Clean-Project {
    Write-Info "`nCleaning project..."
    
    if (Test-Path "bin") {
        Remove-Item -Recurse -Force "bin"
        Write-Success "  ✓ Removed bin/ directory"
    }

    # Clean Go cache
    go clean -cache
    Write-Success "  ✓ Cleaned Go cache"
    
    Write-Success "`n✓ Clean completed!`n"
}

function Show-Help {
    Write-Host @"

Usage: .\build.ps1 [command]

Commands:
  build    - Build binaries only
  test     - Run tests only
  docker   - Start Docker services only
  run      - Run the application (builds first)
  clean    - Remove build artifacts
  all      - Build + Test + Docker (default)
  help     - Show this help

Examples:
  .\build.ps1           # Build, test, and start Docker
  .\build.ps1 build     # Just build binaries
  .\build.ps1 test      # Just run tests
  .\build.ps1 run       # Build and run services

"@ -ForegroundColor White
}

# Main execution
switch ($Command.ToLower()) {
    "build" {
        Build-Project
    }
    "test" {
        Run-Tests
    }
    "docker" {
        Start-Docker
    }
    "run" {
        if (Build-Project) {
            Start-Docker
            Start-Sleep -Seconds 3
            Run-Services
        }
    }
    "clean" {
        Clean-Project
    }
    "all" {
        $success = Build-Project
        if ($success) {
            $success = Run-Tests
        }
        if ($success) {
            Start-Docker
        }
        
        if ($success) {
            Write-Host "`n" + "=" * 50 -ForegroundColor Green
            Write-Success "  🎉 ALL STEPS COMPLETED SUCCESSFULLY!"
            Write-Host "=" * 50 -ForegroundColor Green
            Write-Host @"

Next steps:
  1. Edit .env with your API keys
  2. Run: .\build.ps1 run
  3. Configure GitHub webhook
  4. Create a test PR!

"@ -ForegroundColor Cyan
        } else {
            Write-Host "`n" + "=" * 50 -ForegroundColor Red
            Write-Error "  ✗ BUILD FAILED"
            Write-Host "=" * 50 -ForegroundColor Red
            exit 1
        }
    }
    "help" {
        Show-Help
    }
    default {
        Write-Warning "Unknown command: $Command`n"
        Show-Help
        exit 1
    }
}
