package main

import (
	"fmt"

	"github.com/avressatelier/lokan/internal/store"
	"github.com/urfave/cli/v2"
)

func newClearCmd() *cli.Command {
	return &cli.Command{
		Name:         "clear",
		Usage:        "Delete tasks in bulk",
		OnUsageError: quietUsageError,
		Flags: []cli.Flag{
			boardFlag(),
			&cli.BoolFlag{Name: "archived", Usage: "Delete all tasks in archived lanes"},
			&cli.BoolFlag{Name: "all", Usage: "Delete every task on the board"},
		},
		Action: func(c *cli.Context) error {
			// only a bare invocation is valid
			if err := requireArgs(c, 0); err != nil {
				return err
			}
			board, err := requireBoard(c)
			if err != nil {
				return err
			}

			// require exactly one scope flag
			archived, all := c.Bool("archived"), c.Bool("all")
			if archived == all {
				return cliErrorf("Specify exactly one of --archived or --all")
			}

			// run the requested scope and report the count
			var deleted int
			if archived {
				deleted, err = store.ClearArchived(board)
			} else {
				deleted, err = store.ClearAll(board)
			}
			if err != nil {
				return err
			}
			fmt.Fprintf(c.App.Writer, "Cleared %d task(s).\n", deleted)
			return nil
		},
	}
}
