package main

import (
	"os"
	"path/filepath"

	"cli-task-manager/internal/cli"
	"cli-task-manager/internal/store"
	"cli-task-manager/internal/task"
	"cli-task-manager/internal/timer"
)

func getStoreFilePath() string {
	if envPath := os.Getenv("TASKMAN_DATA_FILE"); envPath != "" {
		return envPath
	}
	homeDir, err := os.UserHomeDir()
	if err == nil && homeDir != "" {
		return filepath.Join(homeDir, ".taskman", "tasks.json")
	}
	return "tasks.json"
}

func main() {
	filePath := getStoreFilePath()
	jsonStore := store.NewJSONStore(filePath)
	taskManager := task.NewTaskManager(jsonStore)
	focusTimer := timer.NewFocusTimer(timer.ProductionDuration, os.Stdout)

	app := cli.NewApp(taskManager, focusTimer, os.Stdout, os.Stderr)
	exitCode := app.Run(os.Args[1:])
	os.Exit(exitCode)
}
