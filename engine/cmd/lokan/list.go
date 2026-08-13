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
		ArgsUsage:    "<board>",
		OnUsageError: quietUsageError,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "type", Usage: "Filter by type: epic, task, subtask, bug"},
			&cli.StringFlag{Name: "status", Usage: "Filter by status: todo, in-progress, backlog, done, cancelled"},
			&cli.StringFlag{Name: "priority", Usage: "Filter by priority: critical, high, medium, low"},
			&cli.BoolFlag{Name: "md", Usage: "Output compact markdown (agent-friendly)"},
		},
		Action: func(c *cli.Context) error {
			// the board path is the required positional argument
			if err := requireArgs(c, 1); err != nil {
				return err
			}
			board, err := requireBoard(c)
			if err != nil {
				return err
			}

			// validate the status filter against the configured lanes
			if st := c.String("status"); st != "" && !contains(statusIDs(board), types.Status(st)) {
				return cliErrorf("Invalid status %q. Must be one of: %s", st, joinStatuses(board))
			}
			cfg, err := store.ReadConfig(board)
			if err != nil {
				return err
			}

			// load, filter, sort the board
			all, err := store.LoadAllSummaries(board)
			if err != nil {
				return err
			}
			filtered := query.FilterTasks(all, types.QueryOptions{
				Type:     types.TaskType(c.String("type")),
				Status:   types.Status(c.String("status")),
				Priority: types.Priority(c.String("priority")),
			})
			sorted := query.SortByPriority(filtered)

			// render as markdown or the aligned table
			if c.Bool("md") {
				fmt.Fprintln(c.App.Writer, renderMarkdownBoard(sorted, cfg.Statuses))
			} else {
				fmt.Fprintln(c.App.Writer, renderTable(sorted))
			}
			return nil
		},
	}
}
