package main

import (
	"fmt"

	"github.com/avressatelier/lokan/internal/query"
	"github.com/avressatelier/lokan/internal/store"
	"github.com/spf13/cobra"
)

func newSubtasksCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "subtasks <id>",
		Short: "List direct children of a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// verify the task exists, then gather + sort its children
			id := args[0]
			root := cmdRoot(cmd)
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
				fmt.Fprintf(cmd.OutOrStdout(), "No subtasks for %s.\n", id)
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Subtasks of %s\n", id)
			for _, child := range children {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", rowLine(child))
			}
			return nil
		},
	}
}
