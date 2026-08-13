# CLI Contract Specification

**Feature**: CLI Task Manager & Focus Timer  
**Spec**: [`specs/001-task-manager-focus-timer/spec.md`](../spec.md)  
**Date**: 2026-08-14  

---

## 1. Command Syntax & Dispatch Matrix

Binary executable: `taskman` (or `go run ./cmd/taskman`)

| Subcommand | Syntax | Description | Required Arguments |
|------------|--------|-------------|--------------------|
| `add` | `taskman add <description>` | Adds a new task | `<description>` (string) |
| `list` / `view` | `taskman list` | Displays all tasks | None |
| `complete` | `taskman complete <id>` | Marks task as completed | `<id>` (integer) |
| `focus` | `taskman focus` | Starts a 25-min focus session | None |
| `help` / `-h` | `taskman help` | Prints CLI usage guide | None |

---

## 2. Standard Output & Error Formatting Contracts

### 2.1 `add` Command

#### Success (Exit Code 0)
```text
Task 1 created: "Implement user authentication"
```

#### Failure: Missing Description (Exit Code 1)
```text
Error: Task description cannot be empty.
Usage: taskman add <description>
```

---

### 2.2 `list` Command

#### Success with Tasks (Exit Code 0)
```text
ID   STATUS       DESCRIPTION
1    [ ]          Implement user authentication
2    [x]          Write unit tests for storage
```

#### Success with No Tasks (Exit Code 0)
```text
No tasks found. Add a task using 'taskman add <description>'.
```

---

### 2.3 `complete` Command

#### Success (Exit Code 0)
```text
Task 1 marked as completed.
```

#### Idempotent / Already Completed (Exit Code 0)
```text
Task 1 is already completed.
```

#### Failure: Non-Existent Task ID (Exit Code 1)
```text
Error: Task with ID 99 was not found.
```

#### Failure: Invalid ID Format (Exit Code 1)
```text
Error: Invalid task ID "abc". Task ID must be a positive integer.
Usage: taskman complete <id>
```

---

### 2.4 `focus` Command

#### Start Output (Exit Code 0)
```text
Focus session started. Duration: 25 minutes.
Press Ctrl+C to cancel session.
```

#### Completion Output (upon timer finish)
```text
Focus session complete! 25 minutes elapsed. Great work.
```

#### Cancellation Output (upon SIGINT / Ctrl+C)
```text
Focus session cancelled.
```

---

### 2.5 Unknown Command / Global Usage

#### Failure: Invalid Subcommand (Exit Code 1)
```text
Error: Unknown command "foo".

Available commands:
  add <description>  Add a new task
  list               List all tasks
  complete <id>      Mark a task as completed
  focus              Start a 25-minute Pomodoro focus session
  help               Show help information
```
