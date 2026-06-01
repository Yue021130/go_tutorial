package logprocessor

import (
	"strings"
	"testing"
)

func TestParseLogLine(t *testing.T) {
	cases := []struct {
		line      string
		wantLevel string
		wantMsg   string
		wantOK    bool
	}{
		{"INFO hello world", "INFO", "hello world", true},
		{"ERROR something went wrong", "ERROR", "something went wrong", true},
		{"WARN low disk", "WARN", "low disk", true},
		{"INVALID", "", "", false},
		{"", "", "", false},
	}

	for _, c := range cases {
		entry, ok := ParseLogLine(c.line)
		if ok != c.wantOK {
			t.Fatalf("ParseLogLine(%q) ok=%v, want %v", c.line, ok, c.wantOK)
		}
		if !ok {
			continue
		}
		if entry.Level != c.wantLevel || entry.Message != c.wantMsg {
			t.Fatalf("ParseLogLine(%q) = %+v, want level=%s msg=%s",
				c.line, entry, c.wantLevel, c.wantMsg)
		}
	}
}

func TestProcessLogs(t *testing.T) {
	input := `INFO request received
ERROR database timeout
WARN slow query
INFO response sent
ERROR connection refused
INFO request received
`
	stats := ProcessLogs(strings.NewReader(input), 2)
	counts := stats.Counts()

	if counts["INFO"] != 3 {
		t.Errorf("INFO count = %d, want 3", counts["INFO"])
	}
	if counts["ERROR"] != 2 {
		t.Errorf("ERROR count = %d, want 2", counts["ERROR"])
	}
	if counts["WARN"] != 1 {
		t.Errorf("WARN count = %d, want 1", counts["WARN"])
	}
}

func BenchmarkProcessLogs(b *testing.B) {
	input := "INFO benchmark log message\n"
	for i := 0; i < b.N; i++ {
		ProcessLogs(strings.NewReader(input), 4)
	}
}
