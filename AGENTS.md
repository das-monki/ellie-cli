# Development Guide for Coding Agents

## Project Overview

This is a Go CLI for the Ellie Daily Planner API, packaged with Nix flake.

## Development Commands

```bash
# Enter development shell
nix develop

# Build locally
go build -o ellie ./cmd/ellie

# Run tests
go test ./...

# Format code
gofmt -w .

# Lint code
staticcheck ./...

# Build Nix package
nix build

# Run built package
./result/bin/ellie --help
```

## Project Structure

```
ellie-cli/
├── cmd/ellie/main.go           # Entry point
├── internal/
│   ├── api/                   # API client and endpoint implementations
│   │   ├── client.go          # HTTP client with x-api-key auth
│   │   ├── tasks.go           # Task operations
│   │   ├── labels.go          # Label operations
│   │   ├── lists.go           # List operations
│   │   └── users.go           # User operations
│   ├── cmd/                   # Cobra commands
│   │   ├── root.go            # Root command + --json flag
│   │   ├── config.go          # Config subcommands
│   │   ├── tasks.go           # Task subcommands
│   │   ├── labels.go          # Label subcommands
│   │   ├── lists.go           # List subcommands
│   │   └── users.go           # User subcommands
│   ├── config/                # Configuration management
│   │   └── config.go          # API key + base URL handling
│   └── models/                # Data models
│       └── models.go          # API request/response types
├── flake.nix                  # Nix flake
└── go.mod                     # Go module
```

## Key Implementation Details

### API Authentication
- Uses `x-api-key` header (not Bearer token)
- Base URL: `https://api.ellieplanner.com`

### API Response Format
- Responses are returned directly (no `data` wrapper)
- Task fields use snake_case: `complete`, `estimated_time`, `created_at`
- Some fields like `date`, `start`, `created_at` may be null or objects

### Adding New Endpoints

1. Add model types in `internal/models/models.go`
2. Add API method in appropriate `internal/api/*.go` file
3. Add Cobra command in `internal/cmd/*.go`
4. Register command in `init()` function

### Testing

- Tests use `httptest` for mocking API responses
- Set `HOME` and `XDG_CONFIG_HOME` to temp dirs for sandboxed tests
- Run `go test ./...` before committing

### Code Style

- Always run `gofmt -w .` and `staticcheck ./...` before committing
- Use `json.RawMessage` for fields with variable types
- Return errors from `RunE` functions; don't print and exit
