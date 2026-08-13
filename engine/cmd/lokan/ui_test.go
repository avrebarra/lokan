package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunUINotInProject(t *testing.T) {
	t.Chdir(t.TempDir())

	var out, errBuf bytes.Buffer
	code := run([]string{"ui"}, &out, &errBuf)
	if code == 0 {
		t.Fatal("expected non-zero exit for non-project directory")
	}
	if !strings.Contains(errBuf.String(), "not a lokan project") {
		t.Fatalf("unexpected error: %q", errBuf.String())
	}
}
