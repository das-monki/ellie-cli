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

# Create a scheduled task. --start is local time; pass --timezone for another zone.
ellie tasks create --desc "Team standup" --date 2025-01-28 --start "09:00"
ellie tasks create --desc "US sync" --date 2025-01-28 --start "09:00" --timezone America/New_York

# Create with estimated time (in seconds)
ellie tasks create --desc "Code review" --estimated-time 1800

# Create with priority (inverted scale: 0=High, 1=Medium, 2=Low)
ellie tasks create --desc "Fix critical bug" --priority 0

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

### `--start` is local time -- do not convert it yourself

`--start "16:00"` means 16:00 on the **machine's local clock**, and `tasks list` /
`tasks get` print start times back in local time with the zone named (`Start: 16:00 CEST`).
So schedule in the user's own hours and say them back plainly. Do **not** pre-convert to
UTC: that was needed before this was fixed, and doing it now shifts every task by the
offset in the wrong direction.

Pass `--timezone` to write a time in some other zone (`--timezone UTC`,
`--timezone America/New_York`). A `--start` given as a full ISO datetime carries its own
offset and is sent unchanged.

Over the wire the API stores instants, so a 16:00 CEST task reads back as
`14:00Z` in `--json` output (raw epoch seconds). That is the same moment, not a bug --
`--json` is for machines, and the human-readable output is already localized.

### `--start` requires `--date`

`ellie tasks update <id> --start "16:00"` fails with
`--date is required when using --start with HH:MM format`. A clock time alone does not say
which day it lands on, so pass `--date` too, even when the date is unchanged.

### `404 Task not found` on update means the ID is wrong, not the route

`tasks update` works -- on tasks made in the CLI and in the app alike. A
`API error (status 404): {"error":"Task not found"}` means the ID string you sent does not
exist, so re-read it from `tasks list --json` before suspecting the endpoint.

The classic way to send a bad ID is shell word-splitting. This bash idiom silently breaks
under zsh, which does not split unquoted expansions, so `$1` becomes the whole `"<id> <time>"`
pair and the ID travels with the time glued onto it:

```bash
# BROKEN in zsh -- $1 is "<id> 12:05", $2 is empty
for pair in "$id1 12:05" "$id2 12:35"; do set -- $pair; ellie tasks update "$1" ...; done
```

Loop over IDs directly instead, and quote them.

### Priority is now 0-2 on an inverted scale

The API used to take priorities 1-4 (Low/Medium/High/Urgent). It now accepts only
**0-2, and the scale is inverted**: `0 = High`, `1 = Medium`, `2 = Low`. So the
highest priority is `--priority 0`, not the biggest number. `--priority 3` or `4` is
rejected client-side with `invalid priority N: must be 0 (High), 1 (Medium), or 2 (Low)`,
and the API itself returns `400 invalid_task_priority` for them.

Priority is **off by default** in the Ellie app -- the API stores the value regardless,
but the UI shows no indicator until you turn it on under *Power features -> Enable task
priority*. So a task can have a priority the UI does not render.

The API also echoes priority back as a JSON **string** (`"0"`), so `--json` output shows
`"priority": "0"`, not a number. That is the API's shape, not a bug.

### No `--version` flag

`ellie --version` errors with `unknown flag: --version`. Use `ellie --help`.

## When to Use

- When the user wants to manage their daily tasks or schedule
- When the user asks about their tasks for today or a specific date
- When the user wants to create, update, or complete tasks
- When the user wants to see unscheduled tasks (braindump)
- When the user mentions ELLI or daily planner
