package main

import (
	"fmt"

	"github.com/avrebarra/lokan/internal/store"
	"github.com/urfave/cli/v2"
)

func newGetCmd() *cli.Command {
	return &cli.Command{
		Name:         "get",
		Usage:        "Show a task by ID",
		ArgsUsage:    "<board> <id>",
		OnUsageError: quietUsageError,
		Action: func(c *cli.Context) error {
			// validate the positional board and id
			if err := requireArgs(c, 2); err != nil {
				return err
			}
			board, err := requireBoard(c)
			if err != nil {
				return err
			}
			id := c.Args().Get(1)

			// resolve the task then print its full detail
			summary, err := store.FindByID(board, id)
			if err != nil {
				return notFoundError(id, err)
			}
			task, err := store.LoadTask(summary.FilePath)
			if err != nil {
				return err
			}
			fmt.Fprintln(c.App.Writer, renderTaskDetail(task))
			return nil
		},
	}
}
