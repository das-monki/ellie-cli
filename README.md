# Ellie CLI

A command-line interface for the [Ellie Daily Planner](https://ellieplanner.com/) REST API.

## Installation

### Using Nix (recommended)

```bash
nix build
./result/bin/ellie --help
```

### From source

```bash
go build -o ellie ./cmd/ellie
```

## Configuration

The CLI needs an API key to authenticate with the Ellie API. Configure it using one of these methods (in priority order):

1. **Environment variable**: `export ELLIE_API_KEY=<your-key>`
2. **File-based** (for secrets management like agenix): `export ELLIE_API_KEY_FILE=/path/to/key`
3. **Config file**: `ellie config set-api-key <your-key>`

Get your API key from the Ellie app settings.

## Usage

```bash
# User info
ellie users me
ellie users usage

# Labels
ellie labels list
ellie labels create --name "Work" --color "#FF5733"

# Lists
ellie lists list

# Tasks
ellie tasks list --date 2024-01-15
ellie tasks by-list --list-id <id>
ellie tasks braindump
ellie tasks get <id>
ellie tasks search "meeting"
ellie tasks create --desc "New task" --date 2024-01-15
ellie tasks update <id> --desc "Updated task"
ellie tasks complete <id>
ellie tasks delete <id>

# JSON output (for scripting)
ellie users me --json
```

## API Documentation

- [Ellie API Documentation](https://ellieplanner.com/api-documentation)

## License

MIT
