@echo off
REM Build and Test Script for Code Review AI Tool
REM Usage: build.bat [command]
REM Commands: build, test, run, docker, clean, all

setlocal enabledelayedexpansion

set COMMAND=%1
if "%COMMAND%"=="" set COMMAND=all

echo.
echo ==================================================
echo   Code Review AI - Build ^& Test System
echo ==================================================
echo.

if "%COMMAND%"=="build" goto BUILD
if "%COMMAND%"=="test" goto TEST
if "%COMMAND%"=="docker" goto DOCKER
if "%COMMAND%"=="run" goto RUN
if "%COMMAND%"=="clean" goto CLEAN
if "%COMMAND%"=="all" goto ALL
if "%COMMAND%"=="help" goto HELP

echo Unknown command: %COMMAND%
goto HELP

:BUILD
echo [1/3] Building binaries...
if not exist bin mkdir bin

echo   Building webhook-listener...
go build -o bin\webhook-listener.exe .\cmd\webhook-listener
if errorlevel 1 goto BUILD_FAILED
echo   [OK] webhook-listener.exe

echo   Building worker...
go build -o bin\worker.exe .\cmd\worker
if errorlevel 1 goto BUILD_FAILED
echo   [OK] worker.exe

echo.
echo [OK] Build completed successfully!
echo.
goto :EOF

:BUILD_FAILED
echo.
echo [ERROR] Build failed!
echo.
exit /b 1

:TEST
echo [2/3] Running tests...
echo   Running unit tests...
go test -v .\...
if errorlevel 1 goto TEST_FAILED

echo.
echo [OK] All tests passed!
echo.
goto :EOF

:TEST_FAILED
echo.
echo [ERROR] Tests failed!
echo.
exit /b 1

:DOCKER
echo [3/3] Starting Docker services...
docker info >nul 2>&1
if errorlevel 1 (
    echo [WARNING] Docker is not running. Please start Docker Desktop.
    exit /b 1
)

echo   Starting RabbitMQ and PostgreSQL...
docker-compose up -d

echo   Waiting for services to be ready...
timeout /t 5 /nobreak >nul

echo   [OK] RabbitMQ: http://localhost:15672 ^(guest/guest^)
echo   [OK] PostgreSQL: localhost:5432
echo.
echo [OK] Docker services started!
echo.
goto :EOF

:RUN
call :BUILD
if errorlevel 1 exit /b 1

call :DOCKER
timeout /t 3 /nobreak >nul

if not exist .env (
    echo [WARNING] .env file not found. Creating from template...
    copy .env.example .env
    echo [WARNING] Please edit .env with your API keys before running!
    pause
    exit /b 0
)

echo.
echo [4/4] Starting services...
echo   Webhook: http://localhost:8080/webhook/github
echo.
echo Press Ctrl+C to stop services
echo.

start "Webhook Listener" cmd /k go run .\cmd\webhook-listener\main.go
timeout /t 2 /nobreak >nul
start "Worker" cmd /k go run .\cmd\worker\main.go

echo [OK] Services started in separate windows!
echo.
goto :EOF

:CLEAN
echo Cleaning project...

if exist bin (
    rmdir /s /q bin
    echo   [OK] Removed bin\ directory
)

go clean -cache
echo   [OK] Cleaned Go cache

echo.
echo [OK] Clean completed!
echo.
goto :EOF

:ALL
call :BUILD
if errorlevel 1 exit /b 1

call :TEST
if errorlevel 1 exit /b 1

call :DOCKER

echo.
echo ==================================================
echo   [OK] ALL STEPS COMPLETED SUCCESSFULLY!
echo ==================================================
echo.
echo Next steps:
echo   1. Edit .env with your API keys
echo   2. Run: build.bat run
echo   3. Configure GitHub webhook
echo   4. Create a test PR!
echo.
goto :EOF

:HELP
echo.
echo Usage: build.bat [command]
echo.
echo Commands:
echo   build    - Build binaries only
echo   test     - Run tests only
echo   docker   - Start Docker services only
echo   run      - Run the application ^(builds first^)
echo   clean    - Remove build artifacts
echo   all      - Build + Test + Docker ^(default^)
echo   help     - Show this help
echo.
echo Examples:
echo   build.bat           # Build, test, and start Docker
echo   build.bat build     # Just build binaries
echo   build.bat test      # Just run tests
echo   build.bat run       # Build and run services
echo.
goto :EOF
