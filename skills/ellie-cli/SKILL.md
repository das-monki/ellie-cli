---
name: ellie-cli
description: CLI for ELLI Daily Planner API
metadata: {"clawdbot":{"requires":{"bins":["ellie"],"env":["ELLIE_API_KEY_FILE"]}}}
---

Use `ellie` to interact with the ELLI Daily Planner API for task management.

## Commands

### Tasks

```bash
# List tasks for a specific date (--date is required)
ellie tasks list --date 2025-01-28

# List tasks with timezone
ellie tasks list --date 2025-01-28 --timezone America/New_York

# Daily agenda: all tasks for a date, including recurring ones (--date is required)
ellie tasks agenda --date 2025-01-28

# Get unscheduled tasks (braindump)
ellie tasks braindump

# Get a specific task
ellie tasks get <task-id>

# Search tasks
ellie tasks search "meeting"

# Create a task
ellie tasks create --desc "Review pull request"

# Create a scheduled task
ellie tasks create --desc "Team standup" --date 2025-01-28 --start "09:00"

# Create with estimated time (in seconds)
ellie tasks create --desc "Code review" --estimated-time 1800

# Create with priority (1=Low, 2=Medium, 3=High, 4=Urgent)
ellie tasks create --desc "Fix critical bug" --priority 4

# Create in a list, with a label
ellie tasks create --desc "Draft spec" --list-id <list-id> --label <label>

# Update a task (also accepts --date, --start, --estimated-time,
# --list-id, --label, --priority, --complete)
ellie tasks update <task-id> --desc "Updated description"

# Mark task complete
ellie tasks complete <task-id>

# Delete a task
ellie tasks delete <task-id>
```

`list` vs `agenda`: `list` returns the tasks scheduled on a date; `agenda` returns
the full daily view for that date, including recurring tasks.

### Lists

```bash
# List all lists
ellie lists list

# Get tasks by list
ellie tasks by-list --list-id <list-id>
```

### Labels

```bash
# List all labels
ellie labels list

# Create a label
ellie labels create --name "Work" --color "#FF5733"
```

### Users

```bash
# Get current user info
ellie users me

# Get API usage
ellie users usage
```

### Configuration

```bash
# Set API key (writes a config file; see "Configuration" below for its location)
ellie config set-api-key <your-api-key>

# Show current configuration
ellie config show

# Point the CLI at a different API base URL
ellie config set-base-url <url>
```

## JSON Output

`--json` is a global flag and works on any command:

```bash
ellie --json tasks list --date 2025-01-28
ellie --json tasks braindump
```

## Configuration

The API key is resolved in this order:
1. `ELLIE_API_KEY` environment variable (the key itself)
2. `ELLIE_API_KEY_FILE` environment variable (path to a file containing the key)
3. Config file, via `ellie config set-api-key`

`ELLIE_BASE_URL` overrides the API base URL.

The config file lives under the OS config directory (`~/Library/Application Support/ellie`
on macOS, `~/.config/ellie` on Linux). It is only created when a command writes config, so
when the key comes from the environment no config directory is needed at all -- which lets
`ellie` run in a sandbox with no write access to that directory.

## When to Use

- When the user wants to manage their daily tasks or schedule
- When the user asks about their tasks for today or a specific date
- When the user wants to create, update, or complete tasks
- When the user wants to see unscheduled tasks (braindump)
- When the user mentions ELLI or daily planner
