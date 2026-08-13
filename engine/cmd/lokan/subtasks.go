package main

import (
	"fmt"

	"github.com/avressatelier/lokan/internal/query"
	"github.com/avressatelier/lokan/internal/store"
	"github.com/urfave/cli/v2"
)

func newSubtasksCmd() *cli.Command {
	return &cli.Command{
		Name:         "subtasks",
		Usage:        "List direct children of a task",
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

			// verify the task exists, then gather + sort its children
			if _, err := store.FindByID(root, id); err != nil {
				return notFoundError(id, err)
			}
			all, err := store.LoadAllSummaries(root)
			if err != nil {
				return err
			}
			children := query.SortByPriority(query.GetChildren(all, id))

			// report children, or a friendly empty message
			if len(children) == 0 {
				fmt.Fprintf(c.App.Writer, "No subtasks for %s.\n", id)
				return nil
			}
			fmt.Fprintf(c.App.Writer, "Subtasks of %s\n", id)
			for _, child := range children {
				fmt.Fprintf(c.App.Writer, "  %s\n", rowLine(child))
			}
			return nil
		},
	}
}
