# Testing Guide

## Overview

The project includes comprehensive unit tests for all core components.

## Running Tests

### Quick Test

```powershell
# Run all tests
.\build.ps1 test
```

### Manual Test Execution

```powershell
# Run all tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test -v ./internal/config
go test -v ./internal/webhook
go test -v ./pkg/llm
```

## Test Coverage

### ✅ Configuration Tests (`internal/config/config_test.go`)

Tests for environment variable loading and validation:
- ✓ Loading configuration from environment
- ✓ Missing GitHub webhook secret validation
- ✓ Missing GitHub token validation
- ✓ Invalid LLM provider validation
- ✓ Missing OpenAI API key validation
- ✓ Valid configuration acceptance

**Coverage:** 6 tests

### ✅ Webhook Handler Tests (`internal/webhook/handler_test.go`)

Tests for GitHub webhook processing and security:
- ✓ Valid HMAC signature verification
- ✓ Invalid signature rejection
- ✓ Event type filtering (only pull_request)
- ✓ Action filtering (opened, synchronize, reopened)

**Coverage:** 4 tests

### ✅ LLM Provider Tests (`pkg/llm/provider_test.go`)

Tests for LLM provider factory and prompt building:
- ✓ Building comprehensive code review prompts
- ✓ Creating OpenAI provider
- ✓ Creating Anthropic provider
- ✓ Creating Ollama provider
- ✓ Unsupported provider error handling
- ✓ Missing API key error handling

**Coverage:** 6 tests

## Current Test Results

```
✓ 12/12 tests passing
✓ 0 test failures
✓ All packages tested
```

## Integration Testing

For integration testing with real services:

1. **Start Docker services:**
   ```powershell
   docker-compose up -d
   ```

2. **Run integration tests** (future enhancement):
   ```powershell
   go test -tags=integration ./...
   ```

## Manual Testing

### Test Webhook Locally

1. Start services:
   ```powershell
   .\build.ps1 run
   ```

2. Use a tool like curl or Postman to send a test webhook:
   ```powershell
   $payload = @{
       action = "opened"
       number = 1
       pull_request = @{
           number = 1
           title = "Test PR"
       }
   } | ConvertTo-Json

   # Calculate HMAC signature
   $secret = "your-webhook-secret"
   $hmac = New-Object System.Security.Cryptography.HMACSHA256
   $hmac.Key = [Text.Encoding]::UTF8.GetBytes($secret)
   $hash = $hmac.ComputeHash([Text.Encoding]::UTF8.GetBytes($payload))
   $signature = "sha256=" + [BitConverter]::ToString($hash).Replace("-", "").ToLower()

   # Send request
   Invoke-WebRequest -Uri http://localhost:8080/webhook/github `
       -Method POST `
       -Headers @{
           "X-GitHub-Event" = "pull_request"
           "X-Hub-Signature-256" = $signature
       } `
       -Body $payload `
       -ContentType "application/json"
   ```

### Test with ngrok

1. Install ngrok: https://ngrok.com/download
2. Expose local server:
   ```powershell
   ngrok http 8080
   ```
3. Configure GitHub webhook with ngrok URL
4. Create a test PR and verify review is posted

## Continuous Integration

To set up CI/CD (GitHub Actions example):

```yaml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run tests
        run: go test -v ./...
```

## Adding New Tests

### Test Structure

```go
package mypackage

import "testing"

func TestMyFunction(t *testing.T) {
    // Arrange
    input := "test input"
    expected := "expected output"

    // Act
    result := MyFunction(input)

    // Assert
    if result != expected {
        t.Errorf("Expected %s, got %s", expected, result)
    }
}
```

### Mock Example

```go
type mockProvider struct {
    response string
    err      error
}

func (m *mockProvider) Analyze(ctx context.Context, prompt string) (string, error) {
    if m.err != nil {
        return "", m.err
    }
    return m.response, nil
}
```

## Benchmarking

To benchmark performance-critical code:

```powershell
go test -bench=. -benchmem ./pkg/llm
```

## Test Best Practices

1. **Table-driven tests** for multiple scenarios
2. **Mock external dependencies** (APIs, databases)
3. **Test edge cases** (empty inputs, errors)
4. **Keep tests fast** (<100ms per test)
5. **Use descriptive test names**

## Future Test Additions

- [ ] Integration tests with real RabbitMQ
- [ ] Integration tests with real PostgreSQL
- [ ] E2E tests with test GitHub repository
- [ ] Performance benchmarks
- [ ] Load testing for webhook endpoint
