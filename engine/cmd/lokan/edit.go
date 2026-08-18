package main

import (
	"fmt"

	"github.com/avrebarra/lokan/internal/store"
	"github.com/avrebarra/lokan/internal/types"
	"github.com/urfave/cli/v2"
)

func newEditCmd() *cli.Command {
	return &cli.Command{
		Name:         "edit",
		Usage:        "Update task fields",
		ArgsUsage:    "<board> <id>",
		OnUsageError: quietUsageError,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "status", Usage: "New status: todo, in-progress, backlog, done, cancelled"},
			&cli.StringFlag{Name: "title", Usage: "New title"},
		},
		Action: func(c *cli.Context) error {
			// validate the positional board and id
			if err := requireArgs(c, 2); err != nil {
				return err
			}
			board, err := requireBoard(c)
			if err != nil {
				return err
			}
			taskID := c.Args().Get(1)
			newStatus := c.String("status")
			newTitle := c.String("title")

			// validate flag enums before touching the store
			if newStatus != "" && !types.Contains(statusIDs(board), types.Status(newStatus)) {
				return cliErrorf("Invalid status %q. Must be one of: %s", newStatus, joinStatuses(board))
			}

			// load the target task
			summary, err := store.FindByID(board, taskID)
			if err != nil {
				return notFoundError(taskID, err)
			}
			task, err := store.LoadTask(summary.FilePath)
			if err != nil {
				return err
			}
			changes := []string{}

			// apply status and title edits
			if newStatus != "" && types.Status(newStatus) != task.Status {
				changes = append(changes, fmt.Sprintf("status: %s → %s", task.Status, newStatus))
				task.Status = types.Status(newStatus)
			}
			titleChanged := newTitle != "" && newTitle != task.Title
			if titleChanged {
				changes = append(changes, fmt.Sprintf("title: %q → %q", task.Title, newTitle))
				task.Title = newTitle
			}

			// nothing to do — report and bail early
			if len(changes) == 0 {
				fmt.Fprintf(c.App.Writer, "No changes for %s.\n", taskID)
				return nil
			}

			if err := store.WriteTask(task); err != nil {
				return err
			}

			// print the change list
			fmt.Fprintf(c.App.Writer, "Updated %s\n", taskID)
			for _, ch := range changes {
				fmt.Fprintf(c.App.Writer, "  %s\n", ch)
			}
			return nil
		},
	}
}
