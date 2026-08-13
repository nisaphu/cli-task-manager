# 0001. Persist Tasks to Local JSON Files

* Status: Accepted
* Date: 2026-08-13

## Context

We are building a CLI Task Manager and Focus Timer tool for developers. The system needs to persist user tasks across process restarts.

As a greenfield CLI application built by a single architect/developer, we need a data storage solution that is lightweight, human-readable, requires minimal setup, and fits naturally into a developer's local terminal environment without introducing unnecessary dependencies.

We considered three persistence options:

1. **Local JSON files** — simple, human-readable, lightweight, and persistent across application restarts.
2. **SQLite database** — provides relational structure, indexing, and stronger querying capabilities, but introduces additional complexity that is not currently required for this small single-user CLI application.
3. **In-memory storage** — simple and fast, but task data would be lost whenever the process exits, which does not satisfy the persistence requirement.

Because the current system is intended for a single local user with a relatively small number of tasks, simplicity and low operational overhead are higher priorities than database scalability or advanced querying.

## Decision

**Persist tasks using local JSON files.**

Task data will be serialized into JSON and stored locally. The application will load the file when it starts and update it whenever tasks are added or marked as completed.

## Consequences

### Positive

* **Zero Configuration:** Users do not need to install or configure an external database.
* **Human-Readable and Debuggable:** Tasks can be inspected using standard terminal tools such as `cat` and `jq`.
* **Low Overhead:** The solution requires minimal infrastructure and dependencies.
* **Fast Initial Development:** The persistence mechanism is simple to implement and suitable for the current project scope.

### Negative / Trade-offs

* **Concurrency Limits:** JSON file writes are not optimized for highly concurrent access, which is acceptable for a single-user CLI tool.
* **Query Performance:** Searching and filtering may require loading and parsing the file rather than using indexed database queries.
* **Scalability Limits:** If the application grows to support large datasets, multiple users, concurrent processes, or more complex relationships, JSON may no longer be suitable.

### Future Considerations

If future requirements introduce more complex querying, concurrency, transactional guarantees, or larger data volumes, the persistence strategy should be revisited.

A future ADR may supersede this decision and introduce a database such as SQLite.
