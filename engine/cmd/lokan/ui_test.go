package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunUIRequiresBoard(t *testing.T) {
	t.Chdir(t.TempDir())

	var out, errBuf bytes.Buffer
	code := run([]string{"ui"}, &out, &errBuf)
	if code == 0 {
		t.Fatal("expected non-zero exit without a board path")
	}
	if !strings.Contains(errBuf.String(), "accepts 1 arg(s), received 0") {
		t.Fatalf("unexpected error: %q", errBuf.String())
	}
}
