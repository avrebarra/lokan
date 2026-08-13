package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/avressatelier/lokan/internal/server"
	"github.com/urfave/cli/v2"
)

const defaultPort = 7777

func newUICmd() *cli.Command {
	return &cli.Command{
		Name:         "ui",
		Usage:        "Start the kanban web UI",
		OnUsageError: quietUsageError,
		Flags: []cli.Flag{
			boardFlag(),
			&cli.IntFlag{Name: "port", Aliases: []string{"p"}, Value: defaultPort, Usage: "port to listen on"},
		},
		Action: runUI,
	}
}

func runUI(c *cli.Context) error {
	port := c.Int("port")
	board, err := requireBoard(c)
	if err != nil {
		return err
	}

	// build the http server over the embedded api
	url := fmt.Sprintf("http://localhost:%d", port)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: server.New(board).Handler(),
	}

	fmt.Printf("lokan ui — %s\n", url)
	fmt.Println("Press Ctrl+C to stop.")
	openBrowser(url)

	// serve until SIGINT/SIGTERM or the listener fails
	errCh := make(chan error, 1)
	go func() {
		errCh <- httpServer.ListenAndServe()
	}()

	// wait for a shutdown signal or a fatal listen error
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	select {
	case <-stop:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return httpServer.Shutdown(ctx)
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

// openBrowser launches the default browser on the current platform.
func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return
	}
	_ = cmd.Start()
}
