package main

import (
	"fmt"

	"github.com/avressatelier/lokan/internal/store"
	"github.com/urfave/cli/v2"
)

func newGetCmd() *cli.Command {
	return &cli.Command{
		Name:         "get",
		Usage:        "Show a task by ID",
		ArgsUsage:    "<id>",
		OnUsageError: quietUsageError,
		Action: func(c *cli.Context) error {
			// validate the positional id, resolve the project
			if err := requireArgs(c, 1); err != nil {
				return err
			}
			root, err := requireProject(c)
			if err != nil {
				return err
			}
			id := c.Args().First()

			// resolve the task then print its full detail
			summary, err := store.FindByID(root, id)
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
