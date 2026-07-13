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

# Reschedule a task. --start needs --date in the SAME command, even if the
# date is not changing -- otherwise the command errors out.
ellie tasks update <task-id> --start "16:00" --date 2025-01-28

# Mark task complete
ellie tasks complete <task-id>

# Delete a task
ellie tasks delete <task-id>
```

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

## Gotchas

Read these before scheduling anything -- each one has already cost a debugging session.

### `tasks agenda` is broken; use `tasks list --date`

`ellie tasks agenda` calls `/v1/tasks/forDate`, which is **not deployed**. It always fails:

```
Error: API error (status 403): {"message":"Missing Authentication Token"}
```

That message is a lie. It is what AWS API Gateway returns for a route that does not
exist, and it has nothing to do with the API key -- the same key works on every other
command. **Do not go debugging the API key or the secret wiring when you see it.**
Use `ellie tasks list --date <date>` instead; it is the working daily view.

### `--start` is UTC, not local time

`--start "16:00"` is sent as **16:00 UTC**, not 16:00 in the user's timezone. When
planning someone's day, either convert their local times to UTC first, or state the
schedule back to them in UTC so the offset is visible. Silently treating it as local
time will place every task wrong by the UTC offset.

`--date`, by contrast, is stored as local midnight, so a task created for `2026-12-31`
reads back as `2026-12-30T23:00:00.000Z` in a UTC+1 timezone. That is expected, not a bug.

### `--start` requires `--date`

`ellie tasks update <id> --start "16:00"` fails with
`--date is required when using --start with HH:MM format`. Always pass `--date` too,
even when the date is unchanged. (An older workaround was to delete and recreate a task
in order to change its time. That is no longer necessary -- `update --start --date`
works.)

### No `--version` flag

`ellie --version` errors with `unknown flag: --version`. Use `ellie --help`.

## When to Use

- When the user wants to manage their daily tasks or schedule
- When the user asks about their tasks for today or a specific date
- When the user wants to create, update, or complete tasks
- When the user wants to see unscheduled tasks (braindump)
- When the user mentions ELLI or daily planner
