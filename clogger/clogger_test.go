package clogger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestLevelMethods(t *testing.T) {
	var tests = []struct {
		name string
		call func(clog *SClog)
		want string
	}{
		{name: "Info", call: func(clog *SClog) { clog.Info("msg") }, want: "Info"},
		{name: "Debug", call: func(clog *SClog) { clog.Debug("msg") }, want: "Debug"},
		{name: "Error", call: func(clog *SClog) { clog.Error("msg") }, want: "Error"},
		{name: "Fatal", call: func(clog *SClog) { clog.Fatal("msg") }, want: "Fatal"},
		{name: "Infof", call: func(clog *SClog) { clog.Infof("msg") }, want: "Info"},
		{name: "Debugf", call: func(clog *SClog) { clog.Debugf("msg") }, want: "Debug"},
		{name: "Errorf", call: func(clog *SClog) { clog.Errorf("msg") }, want: "Error"},
		{name: "Fatalf", call: func(clog *SClog) { clog.Fatalf("msg") }, want: "Fatal"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			var clog SClog
			clog.logger = log.New(&buf, "", 0)
			clog.Name = "test"
			clog.PrefixMask = PrefixLevel
			tt.call(&clog)
			var output = buf.String()
			if !strings.Contains(output, tt.want) {
				t.Errorf("output = %q, missing %q", output, tt.want)
			}
			if !strings.Contains(output, "msg") {
				t.Errorf("output = %q, missing \"msg\"", output)
			}
		})
	}
}

func TestStateFilter(t *testing.T) {
	var stateTrue = func() bool { return true }
	var stateFalse = func() bool { return false }
	var clog = SClog{
		Name:         "test",
		PrefixMask:   PrefixLevel,
		DefaultState: stateFalse,
		LogLevelState: MLogLevelState{
			LogInfo: stateTrue,
		},
	}

	// LogInfo is overridden to true — should produce output
	var bufEnabled bytes.Buffer
	clog.logger = log.New(&bufEnabled, "", 0)
	clog.Info("enabled")
	if bufEnabled.Len() == 0 {
		t.Error("expected output for LogInfo (state overridden to true)")
	}

	// LogDebug uses DefaultState=false — should be suppressed
	var bufDisabled bytes.Buffer
	clog.logger = log.New(&bufDisabled, "", 0)
	clog.Debug("disabled")
	if bufDisabled.Len() > 0 {
		t.Errorf("expected no output for LogDebug, got %q", bufDisabled.String())
	}
}

func TestFatalStackTrace(t *testing.T) {
	var buf bytes.Buffer
	var clog SClog
	clog.logger = log.New(&buf, "", 0)
	clog.Name = "test"
	clog.PrefixMask = PrefixLevel

	clog.Fatal("crash")
	var output = buf.String()
	if !strings.Contains(output, "crash") {
		t.Errorf("output = %q, missing \"crash\"", output)
	}
	if !strings.Contains(output, "goroutine") && !strings.Contains(output, ".go:") {
		t.Errorf("output = %q, missing stack trace (goroutine or .go:)", output)
	}
}
