package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/avrebarra/lokan/internal/server"
	"github.com/urfave/cli/v2"
)

const defaultPort = 17762

func newUICmd() *cli.Command {
	return &cli.Command{
		Name:         "ui",
		Usage:        "Start the kanban web UI",
		ArgsUsage:    "<board>",
		OnUsageError: quietUsageError,
		Flags: []cli.Flag{
			&cli.IntFlag{Name: "port", Aliases: []string{"p"}, Value: defaultPort, Usage: "port to listen on"},
			&cli.BoolFlag{Name: "no-browser", Usage: "serve without opening the browser"},
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

	// bind the port up front. An explicitly requested port is a hard
	// requirement (fail loudly when taken); the default falls back to an
	// OS-assigned free port so a second board can run alongside.
	var ln net.Listener
	if c.IsSet("port") {
		ln, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return fmt.Errorf("ui: port %d is already in use", port)
		}
	} else {
		ln, port, err = listenForUI(port)
		if err != nil {
			return err
		}
	}
	defer ln.Close()

	// build the http server over the embedded api
	url := fmt.Sprintf("http://localhost:%d", port)
	httpServer := &http.Server{Handler: server.New(board).Handler()}

	fmt.Printf("lokan ui — %s\n", url)
	fmt.Println("Press Ctrl+C to stop.")

	// open the browser unless the user opted out (listener already bound,
	// so the page is reachable as soon as it loads)
	if !c.Bool("no-browser") {
		openBrowser(url)
	}

	// serve until SIGINT/SIGTERM or the listener fails
	errCh := make(chan error, 1)
	go func() {
		errCh <- httpServer.Serve(ln)
	}()

	// wait for a shutdown signal or a fatal serve error
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

// openBrowser launches the default browser at the given URL, best-effort:
// a missing opener or spawn failure is reported to stderr, never fatal.
func openBrowser(url string) {
	// resolve the platform opener; unknown platforms just skip
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return
	}
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not open browser: %v\n", err)
	}
}

// listenForUI binds the default port, falling back to a free one when taken.
func listenForUI(port int) (net.Listener, int, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		ln, err = net.Listen("tcp", ":0")
		if err != nil {
			return nil, 0, fmt.Errorf("ui: no free port available (wanted %d): %w", port, err)
		}
	}
	return ln, ln.Addr().(*net.TCPAddr).Port, nil
}
