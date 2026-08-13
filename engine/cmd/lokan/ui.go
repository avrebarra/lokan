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
	"github.com/spf13/cobra"
)

const defaultPort = 7777

func newUICmd() *cobra.Command {
	uiCmd := &cobra.Command{
		Use:   "ui",
		Short: "Start the kanban web UI",
		RunE:  runUI,
	}
	uiCmd.Flags().IntP("port", "p", defaultPort, "port to listen on")
	return uiCmd
}

func runUI(cmd *cobra.Command, args []string) error {
	port, err := cmd.Flags().GetInt("port")
	if err != nil {
		return err
	}
	root := cmdRoot(cmd)

	url := fmt.Sprintf("http://localhost:%d", port)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: server.New(root).Handler(),
	}

	fmt.Printf("lokan ui — %s\n", url)
	fmt.Println("Press Ctrl+C to stop.")
	openBrowser(url)

	// serve until SIGINT/SIGTERM or the listener fails
	errCh := make(chan error, 1)
	go func() {
		errCh <- httpServer.ListenAndServe()
	}()

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
