package main

import (
	"fmt"

	"github.com/avressatelier/lokan/internal/store"
	"github.com/avressatelier/lokan/internal/types"
	"github.com/urfave/cli/v2"
)

func newCreateCmd() *cli.Command {
	return &cli.Command{
		Name:         "create",
		Usage:        "Create a new task",
		ArgsUsage:    "<title>",
		OnUsageError: quietUsageError,
		Flags: []cli.Flag{
			boardFlag(),
			&cli.StringFlag{Name: "type", Aliases: []string{"t"}, Value: string(types.TypeTask), Usage: "Task type: epic, task, subtask, bug"},
			&cli.StringFlag{Name: "priority", Value: string(types.PriorityMedium), Usage: "Priority: critical, high, medium, low"},
			&cli.StringFlag{Name: "parent", Usage: "Parent task ID"},
			&cli.StringSliceFlag{Name: "tag", Usage: "Tag to add (repeatable)"},
		},
		Action: func(c *cli.Context) error {
			// validate the positional title, resolve the board
			if err := requireArgs(c, 1); err != nil {
				return err
			}
			board, err := requireBoard(c)
			if err != nil {
				return err
			}

			// create via the shared flow — same validation as the API
			task, err := store.CreateTaskFromInput(board, c.Args().First(),
				types.TaskType(c.String("type")), types.Priority(c.String("priority")),
				c.String("parent"), c.StringSlice("tag"))
			if err != nil {
				return err
			}
			// print the created task with the targeted board path
			fmt.Fprintf(c.App.Writer, "Created %s → %s\n", task.ID, c.String("board"))
			if parent := c.String("parent"); parent != "" {
				fmt.Fprintf(c.App.Writer, "  Parent: %s\n", parent)
			}
			return nil
		},
	}
}
