# Quick Start Guide

## Prerequisites

1. **Go 1.21+** - [Download Go](https://go.dev/dl/)
2. **Docker Desktop** - [Download Docker](https://www.docker.com/products/docker-desktop/)
3. **GitHub Account** - For testing

## Step 1: Clone and Setup

```powershell
cd C:\Users\carlr\Desktop\Work\Spacework\codereviewtool

# Copy environment template
copy .env.example .env
```

## Step 2: Configure Environment

Edit `.env` and fill in your configuration:

```env
# Required: GitHub Configuration
GITHUB_WEBHOOK_SECRET=mysecretkey123
GITHUB_TOKEN=ghp_your_github_personal_access_token_here

# Required: Choose LLM Provider
LLM_PROVIDER=openai  # or anthropic, or ollama

# If using OpenAI:
OPENAI_API_KEY=sk-your-api-key-here
OPENAI_MODEL=gpt-4-turbo-preview

# If using Anthropic:
# ANTHROPIC_API_KEY=sk-ant-your-api-key-here
# ANTHROPIC_MODEL=claude-3-opus-20240229

# If using Ollama (FREE local option):
# OLLAMA_URL=http://localhost:11434
# OLLAMA_MODEL=codellama
```

### Getting API Keys

**GitHub Token:**
1. Go to https://github.com/settings/tokens
2. Click "Generate new token (classic)"
3. Select scopes: `repo` (all), `write:discussion`
4. Copy the token

**OpenAI API Key:**
1. Go to https://platform.openai.com/api-keys
2. Create new secret key
3. Copy the key

**Anthropic API Key:**
1. Go to https://console.anthropic.com/settings/keys
2. Create API key
3. Copy the key

**Ollama (FREE local option):**
1. Install Ollama: https://ollama.ai
2. Run: `ollama pull codellama`
3. The server runs on http://localhost:11434

## Step 3: Start Services

```powershell
# Start RabbitMQ and PostgreSQL
docker-compose up -d

# Wait a few seconds for services to start
Start-Sleep -Seconds 5
```

## Step 4: Run the Application

Open **two terminal windows**:

**Terminal 1 - Webhook Listener:**
```powershell
go run cmd/webhook-listener/main.go
```

**Terminal 2 - Worker:**
```powershell
go run cmd/worker/main.go
```

You should see:
- Webhook listener: `Webhook listener starting on :8080`
- Worker: `Worker is ready and waiting for pull requests...`

## Step 5: Test Locally (Optional)

To test locally before deploying:

1. **Install ngrok**: https://ngrok.com/download
2. **Expose local port**:
   ```powershell
   ngrok http 8080
   ```
3. **Copy the forwarding URL** (e.g., `https://abc123.ngrok.io`)

## Step 6: Configure GitHub Webhook

1. Go to your test repository on GitHub
2. Click **Settings** → **Webhooks** → **Add webhook**
3. Configure:
   - **Payload URL**: `https://your-server.com/webhook/github` (or ngrok URL for testing)
   - **Content type**: `application/json`
   - **Secret**: Same as `GITHUB_WEBHOOK_SECRET` in `.env`
   - **Events**: Select "Pull requests"
4. Click **Add webhook**

## Step 7: Test with a Pull Request

1. Create a new branch in your repository
2. Make some changes (add a file, modify code)
3. Open a pull request
4. Watch the logs in your terminal windows
5. Within seconds, you should see an AI review posted on your PR!

## Troubleshooting

### "Failed to connect to RabbitMQ"
- Ensure Docker is running: `docker ps`
- Check if services are up: `docker-compose ps`
- Restart services: `docker-compose restart`

### "Invalid signature" error
- Make sure `GITHUB_WEBHOOK_SECRET` matches the secret configured in GitHub webhook settings

### "OpenAI API error" or "Anthropic API error"
- Verify your API key is correct
- Check you have credits/quota available
- Ensure the model name is correct

### No review posted
- Check worker logs for errors
- Verify GitHub token has correct permissions
- Check RabbitMQ management UI: http://localhost:15672 (guest/guest)

## Management Tools

**RabbitMQ Management:**
- URL: http://localhost:15672
- Username: `guest`
- Password: `guest`

**View Queue Messages:** Check if events are being queued properly

## Next Steps

- **Deploy to production**: Use a cloud provider (AWS, GCP, Azure)
- **Add more features**: Custom rules, code formatting checks
- **Fine-tune prompts**: Edit `pkg/llm/provider.go` to customize AI feedback
- **Monitor performance**: Add logging, metrics, alerts

## Architecture

```
GitHub PR → Webhook → Queue → Worker → LLM Analysis → Post Review → GitHub
```

Enjoy your AI code reviewer! 🤖✨
