# CLI Task Manager & Focus Timer

A lightweight command-line productivity tool for developers, built with Go using a Spec-Driven Development workflow.

The application allows users to manage tasks from the terminal and run fixed 25-minute Pomodoro focus sessions.

## Features

* Add a new task from the command line
* View all tasks with completion status
* Mark a task as completed using its ID
* Persist tasks across application restarts
* Run a fixed 25-minute Pomodoro focus timer
* Clear CLI validation and error messages
* Automated unit, integration, and architecture tests

## Architecture

The project follows a separation-of-concerns approach:

```text
cmd/taskman/
└── main.go              # Application entrypoint

internal/
├── cli/                 # CLI commands, routing, and output
├── task/                # Task domain and business logic
├── store/               # Local JSON persistence
└── timer/               # Pomodoro focus timer

architecture/
└── fitness_test.go      # ADR architecture fitness function

adr/
└── 0001-persist-tasks-to-local-json.md

specs/
└── 001-task-manager-focus-timer/
    ├── spec.md
    ├── plan.md
    ├── tasks.md
    ├── research.md
    ├── data-model.md
    └── quickstart.md

.specify/
└── memory/
    └── constitution.md

reflection.md
```

## Architecture Decision

The primary architectural decision is documented in:

```text
adr/0001-persist-tasks-to-local-json.md
```

### ADR 0001: Persist Tasks Using Local JSON Files

Task data is persisted using a local JSON file.

Local JSON was selected because the application is:

* a local single-user CLI tool
* small in scope
* intended to have minimal dependencies
* expected to preserve task data across process restarts

SQLite, SQL databases, ORMs, and alternative persistence mechanisms are intentionally not used unless a future ADR supersedes ADR 0001.

## Spec-Driven Development

The project was developed using Spec Kit and follows this workflow:

```text
Architecture Decision Record
        ↓
/speckit.constitution
        ↓
/speckit.specify
        ↓
/speckit.plan
        ↓
/speckit.tasks
        ↓
/speckit.implement
        ↓
/speckit.converge
```

The specification and architecture decisions act as governance constraints for implementation.

The final `/speckit.converge` assessment reported:

* 17 / 17 functional requirements satisfied
* 5 / 5 success criteria satisfied
* 12 / 12 acceptance scenarios satisfied
* 7 / 7 constitution principles satisfied
* ADR 0001 compliance verified
* 0 architectural gaps found

No corrective tasks were required after convergence.

## Requirements

* Go 1.22 or later

Check your Go installation:

```bash
go version
```

## Build

Clone the repository and enter the project directory:

```bash
git clone https://github.com/nisaphu/cli-task-manager.git
cd cli-task-manager
```

Build the CLI:

```bash
go build -o taskman ./cmd/taskman
```

## Usage

### Add a Task

```bash
./taskman add "Write unit tests"
```

Example output:

```text
Task 1 created: "Write unit tests"
```

### List Tasks

```bash
./taskman list
```

Example:

```text
ID   STATUS   DESCRIPTION
1    [ ]      Write unit tests
2    [x]      Review architecture
```

`[ ]` represents an incomplete task.

`[x]` represents a completed task.

### Complete a Task

```bash
./taskman complete 1
```

Example:

```text
Task 1 marked as completed.
```

### Start a Focus Session

```bash
./taskman focus
```

The production focus session duration is fixed at:

```text
25 minutes
```

The application notifies the user when the session starts and when the session is completed.

Use `Ctrl+C` to cancel an active session.

### Help

```bash
./taskman help
```

## Task Persistence

Task data is stored locally as JSON and remains available after the CLI exits and is started again.

The persistence layer uses Go's standard library:

```go
encoding/json
```

Writes use temporary-file replacement to reduce the risk of leaving partially written task data.

The application also handles:

* missing JSON files
* corrupted JSON data
* task reload after restart

## Testing

Run the complete automated test suite:

```bash
go test ./...
```

To force all tests to execute without using cached results:

```bash
go test -count=1 ./...
```

The test suite covers:

* task creation
* empty task descriptions
* unique task IDs
* task listing
* task completion
* non-existent task IDs
* already-completed tasks
* JSON save and reload
* missing JSON files
* corrupted JSON
* CLI routing and error handling
* Pomodoro timer behavior
* persistence architecture compliance

## Architecture Fitness Function

The project contains an automated architecture fitness function:

```text
architecture/fitness_test.go
```

It enforces ADR 0001 by failing the test suite if forbidden persistence technologies are introduced, including examples such as:

```text
database/sql
SQLite drivers
GORM
Ent
BoltDB
BadgerDB
```

It also verifies that the JSON persistence implementation continues to use:

```text
encoding/json
```

Because the fitness function is part of the Go test suite, architectural governance is automatically checked whenever:

```bash
go test ./...
```

is executed.

## Governance

Architectural decisions are documented as ADRs and treated as project-level guardrails.

The project constitution requires:

* simple and maintainable code
* compliance with accepted ADRs
* separation of concerns
* testable business logic
* specification-driven implementation
* minimal infrastructure complexity
* correction of architectural drift before completion

See:

```text
.specify/memory/constitution.md
```

## Reflection

The socio-technical reflection for the assignment is available in:

```text
reflection.md
```

It discusses:

* architectural guardrails
* closed-loop correction using `/speckit.converge`
* human-in-the-loop governance
* Decision Guardian
* automated fitness functions
* Decisions as Code

## Technology Stack

* Go
* Go Standard Library
* Local JSON persistence
* Spec Kit
* Architecture Decision Records
* Automated Go testing

No external application dependencies or database engines are required.

## Author

Nisa Phutthanawong
