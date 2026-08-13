package main

import (
	"path/filepath"
	"strings"

	"github.com/avressatelier/lokan/internal/store"
	"github.com/urfave/cli/v2"
)

// boardFlag is the explicit board file flag shared by every command.
func boardFlag() cli.Flag {
	return &cli.StringFlag{Name: "board", Usage: "Board file (markdown with a lokan config block); required"}
}

// requireBoard resolves and validates the explicit --board file. Every
// command except init and help needs a board.
func requireBoard(c *cli.Context) (string, error) {
	board := c.String("board")
	if board == "" {
		return "", cliErrorf("missing required flag: --board <file>")
	}
	abs, err := filepath.Abs(board)
	if err != nil {
		return "", err
	}
	if !store.IsBoard(abs) {
		return "", cliErrorf("not a lokan board: %s (run lokan init --board <file>)", abs)
	}
	return abs, nil
}

// requireArgs validates that the command got exactly n positional args.
func requireArgs(c *cli.Context, n int) error {
	if got := c.Args().Len(); got != n {
		return cliErrorf("accepts %d arg(s), received %d", n, got)
	}
	return nil
}

// quietUsageError suppresses urfave's "Incorrect Usage" banner so errors only
// surface through run()'s single "error: " line, as before.
func quietUsageError(_ *cli.Context, err error, _ bool) error {
	return err
}

// reorderFlags moves interspersed flags ahead of positional args, replicating
// the pflag (interspersed) parsing urfave lacks. The command name is kept in
// place; value flags consume the next token; "--" and boolean help flags pass
// through untouched.
func reorderFlags(args []string) []string {
	if len(args) < 2 {
		return args
	}
	var flags, positional []string
	afterDash := false
	for i := 1; i < len(args); i++ {
		a := args[i]
		switch {
		case afterDash:
			positional = append(positional, a)
		case a == "--":
			afterDash = true
			positional = append(positional, a)
		case isValueFlag(a):
			flags = append(flags, a)
			if !strings.Contains(a, "=") && i+1 < len(args) {
				flags = append(flags, args[i+1])
				i++
			}
		default:
			positional = append(positional, a)
		}
	}
	return append([]string{args[0]}, append(flags, positional...)...)
}

// isValueFlag reports whether a token is a non-help flag that takes a value.
func isValueFlag(a string) bool {
	if a == "-" || a == "--help" || a == "-h" {
		return false
	}
	return strings.HasPrefix(a, "-")
}

// newApp assembles the command tree.
func newApp() *cli.App {
	app := &cli.App{
		Name:  "lokan",
		Usage: "Markdown task manager for Claude Code",
		Commands: []*cli.Command{
			newInitCmd(),
			newCreateCmd(),
			newGetCmd(),
			newListCmd(),
			newEditCmd(),
			newSubtasksCmd(),
			newClearCmd(),
			newUICmd(),
		},
		OnUsageError: quietUsageError,
		// bare invocation shows help; unknown commands are errors
		Action: func(c *cli.Context) error {
			if c.Args().Present() {
				return cliErrorf("unknown command %q for \"lokan\"", c.Args().First())
			}
			return cli.ShowAppHelp(c)
		},
	}
	return app
}
