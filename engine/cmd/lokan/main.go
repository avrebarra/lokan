// Command lokan is the lokan markdown task manager CLI.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run executes the app with args, writing output to stdout and errors to
// stderr. It returns the process exit code.
func run(args []string, stdout, stderr io.Writer) int {
	app := newApp()
	app.Writer = stdout
	app.ErrWriter = stderr
	// silence urfave's own error printing — run() owns the single error line
	app.ExitErrHandler = func(_ *cli.Context, _ error) {}
	if err := app.Run(append([]string{app.Name}, reorderFlags(args)...)); err != nil {
		fmt.Fprintf(stderr, "error: %s\n", err)
		return 1
	}
	return 0
}
