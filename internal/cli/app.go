package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"cli-task-manager/internal/task"
	"cli-task-manager/internal/timer"
)

// App manages command line input, routing, and output formatting.
type App struct {
	taskManager *task.TaskManager
	timerRunner timer.Runner
	out         io.Writer
	errOut      io.Writer
}

// NewApp creates a new CLI App instance with custom I/O streams for clean testing.
func NewApp(tm *task.TaskManager, tr timer.Runner, out io.Writer, errOut io.Writer) *App {
	if out == nil {
		out = os.Stdout
	}
	if errOut == nil {
		errOut = os.Stderr
	}
	if tr == nil {
		tr = timer.NewFocusTimer(timer.ProductionDuration, out)
	}
	return &App{
		taskManager: tm,
		timerRunner: tr,
		out:         out,
		errOut:      errOut,
	}
}

// Run parses command line arguments and routes execution to the corresponding handler.
func (a *App) Run(args []string) int {
	if len(args) < 1 {
		a.printHelp()
		return 0
	}

	command := args[0]
	cmdArgs := args[1:]

	switch command {
	case "add":
		return a.handleAdd(cmdArgs)
	case "list", "view":
		return a.handleList()
	case "complete":
		return a.handleComplete(cmdArgs)
	case "focus":
		return a.handleFocus()
	case "help", "-h", "--help":
		a.printHelp()
		return 0
	default:
		fmt.Fprintf(a.errOut, "Error: Unknown command %q.\n\n", command)
		a.printAvailableCommands()
		return 1
	}
}

func (a *App) handleAdd(args []string) int {
	if len(args) < 1 || args[0] == "" {
		fmt.Fprintln(a.errOut, "Error: Task description cannot be empty.")
		fmt.Fprintln(a.errOut, "Usage: taskman add <description>")
		return 1
	}

	description := args[0]
	t, err := a.taskManager.AddTask(description)
	if err != nil {
		if errors.Is(err, task.ErrEmptyDescription) {
			fmt.Fprintln(a.errOut, "Error: Task description cannot be empty.")
			fmt.Fprintln(a.errOut, "Usage: taskman add <description>")
			return 1
		}
		fmt.Fprintf(a.errOut, "Error: %v\n", err)
		return 1
	}

	fmt.Fprintf(a.out, "Task %d created: %q\n", t.ID, t.Description)
	return 0
}

func (a *App) handleList() int {
	tasks, err := a.taskManager.ListTasks()
	if err != nil {
		fmt.Fprintf(a.errOut, "Error loading tasks: %v\n", err)
		return 1
	}

	if len(tasks) == 0 {
		fmt.Fprintln(a.out, "No tasks found. Add a task using 'taskman add <description>'.")
		return 0
	}

	fmt.Fprintln(a.out, "ID   STATUS       DESCRIPTION")
	for _, t := range tasks {
		fmt.Fprintf(a.out, "%-4d %-12s %s\n", t.ID, t.StatusTag(), t.Description)
	}
	return 0
}

func (a *App) handleComplete(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(a.errOut, "Error: Missing task ID.")
		fmt.Fprintln(a.errOut, "Usage: taskman complete <id>")
		return 1
	}

	rawID := args[0]
	id, err := strconv.Atoi(rawID)
	if err != nil || id <= 0 {
		fmt.Fprintf(a.errOut, "Error: Invalid task ID %q. Task ID must be a positive integer.\n", rawID)
		fmt.Fprintln(a.errOut, "Usage: taskman complete <id>")
		return 1
	}

	t, err := a.taskManager.CompleteTask(id)
	if err != nil {
		if errors.Is(err, task.ErrTaskAlreadyCompleted) {
			fmt.Fprintf(a.out, "Task %d is already completed.\n", id)
			return 0
		}
		if errors.Is(err, task.ErrTaskNotFound) {
			fmt.Fprintf(a.errOut, "Error: Task with ID %d was not found.\n", id)
			return 1
		}
		fmt.Fprintf(a.errOut, "Error: %v\n", err)
		return 1
	}

	fmt.Fprintf(a.out, "Task %d marked as completed.\n", t.ID)
	return 0
}

func (a *App) handleFocus() int {
	if err := a.timerRunner.Start(); err != nil {
		if errors.Is(err, timer.ErrCancelled) {
			fmt.Fprintln(a.out, "Focus session cancelled.")
			return 0
		}
		fmt.Fprintf(a.errOut, "Error starting focus timer: %v\n", err)
		return 1
	}
	return 0
}

func (a *App) printHelp() {
	fmt.Fprintln(a.out, "CLI Task Manager & Focus Timer")
	fmt.Fprintln(a.out)
	a.printAvailableCommands()
}

func (a *App) printAvailableCommands() {
	fmt.Fprintln(a.out, "Available commands:")
	fmt.Fprintln(a.out, "  add <description>  Add a new task")
	fmt.Fprintln(a.out, "  list               List all tasks")
	fmt.Fprintln(a.out, "  complete <id>      Mark a task as completed")
	fmt.Fprintln(a.out, "  focus              Start a 25-minute Pomodoro focus session")
	fmt.Fprintln(a.out, "  help               Show help information")
}
