package main

import (
	"bytes"
	"net"
	"path/filepath"
	"strconv"
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

func TestListenForUIPortInUse(t *testing.T) {
	occupied, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("reserve a port: %v", err)
	}
	defer occupied.Close()
	want := occupied.Addr().(*net.TCPAddr).Port

	ln, port, err := listenForUI(want)
	if err != nil {
		t.Fatalf("listenForUI: %v", err)
	}
	defer ln.Close()

	if port == want {
		t.Fatalf("expected a fallback port, got %d", want)
	}
	if got := ln.Addr().(*net.TCPAddr).Port; got != port {
		t.Fatalf("reported port %d does not match bound port %d", port, got)
	}
}

func TestListenForUIPicksRequestedPort(t *testing.T) {
	probe, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("probe a port: %v", err)
	}
	want := probe.Addr().(*net.TCPAddr).Port
	probe.Close()

	ln, port, err := listenForUI(want)
	if err != nil {
		t.Fatalf("listenForUI: %v", err)
	}
	defer ln.Close()

	if port != want {
		t.Fatalf("expected port %d, got %d", want, port)
	}
}

func TestRunUIFailsOnExplicitPortClash(t *testing.T) {
	// occupy a port so the explicit --port cannot bind
	occupied, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("reserve a port: %v", err)
	}
	defer occupied.Close()
	taken := occupied.Addr().(*net.TCPAddr).Port

	// run ui with an explicit port that is taken — must fail, not fall back
	board := filepath.Join(t.TempDir(), "board.md")
	if code := run([]string{"init", board}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatal("failed to init board")
	}
	var out, errBuf bytes.Buffer
	code := run([]string{"ui", "--no-browser", "--port", strconv.Itoa(taken), board}, &out, &errBuf)
	if code == 0 {
		t.Fatal("expected non-zero exit when explicit port is taken")
	}
	if !strings.Contains(errBuf.String(), "already in use") {
		t.Fatalf("expected 'already in use' error, got %q", errBuf.String())
	}
}
