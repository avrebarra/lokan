package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
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
		ArgsUsage:    "<board>",
		OnUsageError: quietUsageError,
		Flags: []cli.Flag{
			&cli.IntFlag{Name: "port", Aliases: []string{"p"}, Value: defaultPort, Usage: "port to listen on"},
		},
		Action: runUI,
	}
}

func runUI(c *cli.Context) error {
	// the board path is the required positional argument
	if err := requireArgs(c, 1); err != nil {
		return err
	}
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
