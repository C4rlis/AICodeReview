# Git Setup and Initial Commit

## Git Not Found

Git is not currently available in your PATH. You have two options:

### Option 1: Install Git (Recommended)

1. Download Git for Windows: https://git-scm.com/download/win
2. Install with default settings
3. Restart your terminal/PowerShell
4. Run the commands below

### Option 2: Use GitHub Desktop

1. Download GitHub Desktop: https://desktop.github.com/
2. Open the project folder in GitHub Desktop
3. It will detect changes and help you commit

## Manual Commit Commands

Once git is installed, run these commands in PowerShell from the project directory:

```powershell
# Navigate to project
cd C:\Users\carlr\Desktop\Work\Spacework\codereviewtool

# Initialize git repository (if not already done)
git init

# Add all files
git add .

# Create initial commit
git commit -m "Initial commit: AI Code Review Tool

- Complete webhook listener and worker implementation
- Configurable LLM providers (OpenAI/Anthropic/Ollama)
- GitHub webhook integration with HMAC verification
- RabbitMQ message queue for async processing
- Comprehensive unit tests (12 passing)
- Build scripts (PowerShell and batch)
- Docker Compose setup for local development
- Full documentation (README, QUICKSTART, SETUP, TESTING, BUILD)"

# View commit history
git log --oneline
```

## Quick Commit Script

Save this as `commit.bat` for quick commits:

```batch
@echo off
git add .
git status
echo.
set /p message="Commit message: "
git commit -m "%message%"
echo.
echo Commit complete!
```

Usage:
```powershell
.\commit.bat
```

## Pushing to GitHub

Once you create a repository on GitHub:

```powershell
# Add remote repository
git remote add origin https://github.com/yourusername/codereviewtool.git

# Push to GitHub
git push -u origin main
```

## .gitignore Already Created

The `.gitignore` file is already set up to exclude:
- Binary files (`bin/`, `*.exe`)
- Environment files (`.env`)
- IDE files (`.vscode/`, `.idea/`)
- OS files (`.DS_Store`, `Thumbs.db`)
- Build artifacts

## Current Project Status

All files are ready to commit:
- ✅ 18 source files created
- ✅ 12 unit tests passing
- ✅ Build system working
- ✅ Documentation complete
- ✅ `.gitignore` configured

Total files to commit: ~30 files

## Recommended First Commit

Once git is available:

```powershell
git init
git add .
git commit -m "Initial commit: Complete AI code review tool implementation"
```

This will commit everything in one go with a clean starting point.
