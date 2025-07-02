package core_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/shellvon/go-sender/core"
)

func TestLevelStringAndKey(t *testing.T) {
	levels := []core.Level{core.LevelDebug, core.LevelInfo, core.LevelWarn, core.LevelError, core.LevelFatal, 99}
	expects := []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL", ""}
	for i, lv := range levels {
		if lv.String() != expects[i] {
			t.Errorf("Level %v string got %q, want %q", lv, lv.String(), expects[i])
		}
		if lv.Key() != core.LevelKey {
			t.Errorf("Level.Key() should be %q", core.LevelKey)
		}
	}
}

func TestStdLogger_LogAndWith(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	std := core.NewStdLogger(logger)
	if err := std.Log(core.LevelInfo, "msg", "test"); err != nil {
		t.Errorf("StdLogger.Log error: %v", err)
	}
	l2 := std.With("foo", "bar")
	if l2 == nil {
		t.Error("With should return a logger")
	}
	_ = l2.Log(core.LevelWarn, "msg", "with test")
}

func TestNoOpLogger(t *testing.T) {
	var l core.Logger = &core.NoOpLogger{}
	if err := l.Log(core.LevelInfo, "msg", "noop"); err != nil {
		t.Errorf("NoOpLogger.Log should not error: %v", err)
	}
	if l2 := l.With("foo", "bar"); l2 == nil {
		t.Error("NoOpLogger.With should return self")
	}
}
