# Quickstart & Verification Guide: CLI Task Manager & Focus Timer

**Feature**: CLI Task Manager & Focus Timer  
**Spec**: [`specs/001-task-manager-focus-timer/spec.md`](spec.md)  
**Contract**: [`specs/001-task-manager-focus-timer/contracts/cli-contract.md`](contracts/cli-contract.md)  
**Date**: 2026-08-14  

---

## 1. Environment & Setup

- **Language**: Go 1.22+ installed
- **Dependencies**: Go Standard Library only (no `go get` or external packages required)

---

## 2. Verification Workflows

### 2.1 Run Automated Tests

Execute the automated test suite covering domain logic, persistence reload, corrupted JSON, and focus timer:

```bash
go test -v ./...
```

Expected output: `PASS` across all package tests (`internal/task`, `internal/store`, `internal/timer`, `internal/cli`).

---

### 2.2 Build Binary

Build the local executable:

```bash
go build -o taskman ./cmd/taskman
```

---

### 2.3 Manual End-to-End Validation Scenarios

#### Scenario 1: Add Tasks
```bash
./taskman add "Write unit tests"
# Expected: Task 1 created: "Write unit tests"

./taskman add "Implement focus timer"
# Expected: Task 2 created: "Implement focus timer"
```

#### Scenario 2: List Tasks
```bash
./taskman list
# Expected:
# ID   STATUS       DESCRIPTION
# 1    [ ]          Write unit tests
# 2    [ ]          Implement focus timer
```

#### Scenario 3: Complete Task & State Persistence
```bash
./taskman complete 1
# Expected: Task 1 marked as completed.

# Verify state after complete:
./taskman list
# Expected: Task 1 displays [x], Task 2 displays [ ]

# Test Persistence Across Process Restarts:
# (Data saved to local JSON store file)
./taskman list
# Expected: Completed status for Task 1 remains intact.
```

#### Scenario 4: Error Handling & Idempotency
```bash
# Complete already completed task:
./taskman complete 1
# Expected: Task 1 is already completed.

# Complete non-existent task ID:
./taskman complete 999
# Expected: Error: Task with ID 999 was not found.

# Add task with empty description:
./taskman add ""
# Expected: Error: Task description cannot be empty.
```

#### Scenario 5: Focus Timer
```bash
./taskman focus
# Expected:
# Focus session started. Duration: 25 minutes.
# Press Ctrl+C to cancel session.
```
