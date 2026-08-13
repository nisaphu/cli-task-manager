package timer

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ErrCancelled indicates that a focus session was interrupted by the user.
var ErrCancelled = errors.New("focus session cancelled")

// ProductionDuration is fixed at 25 minutes per the Pomodoro technique.
const ProductionDuration = 25 * time.Minute

// Runner defines the interface for running a focus timer session.
type Runner interface {
	Start() error
}

// FocusTimer manages Pomodoro timing sessions.
type FocusTimer struct {
	duration time.Duration
	out      io.Writer
}

// NewFocusTimer creates a FocusTimer with a custom duration and output writer.
func NewFocusTimer(duration time.Duration, out io.Writer) *FocusTimer {
	if out == nil {
		out = os.Stdout
	}
	return &FocusTimer{
		duration: duration,
		out:      out,
	}
}

// Start begins the focus session, listening for termination signals or completion.
func (ft *FocusTimer) Start() error {
	minutes := int(ft.duration.Minutes())
	if minutes < 1 {
		fmt.Fprintf(ft.out, "Focus session started. Duration: %v.\n", ft.duration)
	} else {
		fmt.Fprintf(ft.out, "Focus session started. Duration: %d minutes.\n", minutes)
	}
	fmt.Fprintln(ft.out, "Press Ctrl+C to cancel session.")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	timer := time.NewTimer(ft.duration)
	defer timer.Stop()

	select {
	case <-timer.C:
		if minutes < 1 {
			fmt.Fprintf(ft.out, "Focus session complete! %v elapsed. Great work.\n", ft.duration)
		} else {
			fmt.Fprintf(ft.out, "Focus session complete! %d minutes elapsed. Great work.\n", minutes)
		}
		return nil
	case <-sigChan:
		return ErrCancelled
	}
}
