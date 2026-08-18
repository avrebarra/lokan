package main

import (
	"fmt"

	"github.com/avrebarra/lokan/internal/store"
	"github.com/urfave/cli/v2"
)

func newCreateCmd() *cli.Command {
	return &cli.Command{
		Name:         "create",
		Usage:        "Create a new task",
		ArgsUsage:    "<board> <title>",
		OnUsageError: quietUsageError,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{Name: "tag", Usage: "Tag to add (repeatable)"},
		},
		Action: func(c *cli.Context) error {
			// validate the positional board and title
			if err := requireArgs(c, 2); err != nil {
				return err
			}
			title := c.Args().Get(1)

			// resolve the board, then create via the shared flow — same
			// validation as the API
			board, err := requireBoard(c)
			if err != nil {
				return err
			}
			task, err := store.CreateTaskFromInput(board, title, c.StringSlice("tag"))
			if err != nil {
				return err
			}
			// print the created task with the targeted board path
			fmt.Fprintf(c.App.Writer, "Created %s → %s\n", task.ID, board)
			return nil
		},
	}
}
