package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/avrebarra/lokan/internal/store"
	"github.com/avrebarra/lokan/internal/types"
	"github.com/urfave/cli/v2"
)

func newInitCmd() *cli.Command {
	return &cli.Command{
		Name:         "init",
		Usage:        "Initialize a lokan board at an explicit path",
		ArgsUsage:    "<board>",
		OnUsageError: quietUsageError,
		Action: func(c *cli.Context) error {
			// the board path is the required positional argument
			if err := requireArgs(c, 1); err != nil {
				return err
			}
			board := c.Args().First()
			abs, err := filepath.Abs(board)
			if err != nil {
				return err
			}
			out := c.App.Writer

			// already a board at that path — report and bail
			if store.IsBoard(abs) {
				fmt.Fprintf(out, "Already a lokan board: %s\n", abs)
				return nil
			}
			if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
				return err
			}

			// write the self-contained board: config block + empty sections
			raw, err := store.InitialBoard(types.LokanConfig{
				Counter:  0,
				Version:  "1",
				Statuses: types.DefaultStatusDefs(),
			})
			if err != nil {
				return err
			}
			if err := os.WriteFile(abs, []byte(raw), 0o644); err != nil {
				return err
			}

			// report the created path
			fmt.Fprintf(out, "Initialized lokan project.\n")
			fmt.Fprintf(out, "  Board:  %s\n", abs)
			return nil
		},
	}
}
