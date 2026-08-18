package main

import (
	"fmt"

	"github.com/avrebarra/lokan/internal/query"
	"github.com/avrebarra/lokan/internal/store"
	"github.com/avrebarra/lokan/internal/types"
	"github.com/urfave/cli/v2"
)

func newListCmd() *cli.Command {
	return &cli.Command{
		Name:         "list",
		Usage:        "List tasks with optional filters",
		ArgsUsage:    "<board>",
		OnUsageError: quietUsageError,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "status", Usage: "Filter by status: todo, in-progress, backlog, done, cancelled"},
			&cli.StringSliceFlag{Name: "tag", Usage: "Filter by tag (comma-separated or repeatable, AND semantics)"},
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
			if st := c.String("status"); st != "" && !types.Contains(statusIDs(board), types.Status(st)) {
				return cliErrorf("Invalid status %q. Must be one of: %s", st, joinStatuses(board))
			}
			cfg, err := store.ReadConfig(board)
			if err != nil {
				return err
			}

			// load and filter the board
			all, err := store.LoadAllSummaries(board)
			if err != nil {
				return err
			}
			filtered := query.FilterTasks(all, types.QueryOptions{
				Status: types.Status(c.String("status")),
				Tags:   c.StringSlice("tag"),
			})

			// render as markdown or the aligned table
			if c.Bool("md") {
				fmt.Fprintln(c.App.Writer, renderMarkdownBoard(filtered, cfg.Statuses))
			} else {
				fmt.Fprintln(c.App.Writer, renderTable(filtered))
			}
			return nil
		},
	}
}
