package clogger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestPrefixShown(t *testing.T) {
	var buf bytes.Buffer
	var clog SClog
	clog.logger = log.New(&buf, "", 0)
	clog.Name = "test"
	clog.PrefixMask = PrefixName | PrefixLevel

	clog.Info("msg")
	var output = buf.String()
	if !strings.Contains(output, "test") {
		t.Errorf("output = %q, missing name \"test\"", output)
	}
	if !strings.Contains(output, "Info") {
		t.Errorf("output = %q, missing level \"Info\"", output)
	}
	if !strings.Contains(output, "msg") {
		t.Errorf("output = %q, missing \"msg\"", output)
	}
}

func TestPrefixSuppressed(t *testing.T) {
	var buf bytes.Buffer
	var clog SClog
	clog.logger = log.New(&buf, "", 0)
	clog.Name = "test"
	clog.PrefixMask = 0

	clog.Info("msg")
	var output = buf.String()
	if !strings.Contains(output, "msg") {
		t.Errorf("output = %q, missing \"msg\"", output)
	}
	if strings.Contains(output, ": ") {
		t.Errorf("output = %q, contains \": \" but PrefixMask is 0", output)
	}
}
