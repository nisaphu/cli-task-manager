package timer

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestFocusTimer_Execution(t *testing.T) {
	t.Run("Fast sub-second completion test", func(t *testing.T) {
		var buf bytes.Buffer
		shortDuration := 10 * time.Millisecond
		ft := NewFocusTimer(shortDuration, &buf)

		start := time.Now()
		err := ft.Start()
		elapsed := time.Since(start)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if elapsed < shortDuration {
			t.Errorf("expected elapsed time to be at least %v, got %v", shortDuration, elapsed)
		}

		output := buf.String()
		if !strings.Contains(output, "Focus session started") {
			t.Errorf("expected start notice, got: %s", output)
		}
		if !strings.Contains(output, "Focus session complete") {
			t.Errorf("expected completion notice, got: %s", output)
		}
	})

	t.Run("Production duration is 25 minutes", func(t *testing.T) {
		if ProductionDuration != 25*time.Minute {
			t.Errorf("expected ProductionDuration to be 25 minutes, got %v", ProductionDuration)
		}
	})
}
