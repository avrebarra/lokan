package main

import (
	"fmt"

	"github.com/avressatelier/lokan/internal/query"
	"github.com/avressatelier/lokan/internal/store"
	"github.com/avressatelier/lokan/internal/types"
	"github.com/urfave/cli/v2"
)

func newListCmd() *cli.Command {
	return &cli.Command{
		Name:         "list",
		Usage:        "List tasks with optional filters",
		OnUsageError: quietUsageError,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "type", Usage: "Filter by type: epic, task, subtask, bug"},
			&cli.StringFlag{Name: "status", Usage: "Filter by status: todo, in-progress, backlog, done, cancelled"},
			&cli.StringFlag{Name: "priority", Usage: "Filter by priority: critical, high, medium, low"},
		},
		Action: func(c *cli.Context) error {
			// only a bare invocation is valid
			if err := requireArgs(c, 0); err != nil {
				return err
			}
			root, err := requireProject(c)
			if err != nil {
				return err
			}

			// load, filter, sort the board
			all, err := store.LoadAllSummaries(root)
			if err != nil {
				return err
			}
			filtered := query.FilterTasks(all, types.QueryOptions{
				Type:     types.TaskType(c.String("type")),
				Status:   types.Status(c.String("status")),
				Priority: types.Priority(c.String("priority")),
			})
			sorted := query.SortByPriority(filtered)

			fmt.Fprintln(c.App.Writer, renderTable(sorted))
			return nil
		},
	}
}
