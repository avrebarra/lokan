package main

import (
	"fmt"

	"github.com/avressatelier/lokan/internal/query"
	"github.com/avressatelier/lokan/internal/store"
	"github.com/avressatelier/lokan/internal/types"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	var filterType, filterStatus, filterPriority string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tasks with optional filters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			root := cmdRoot(cmd)
			all, err := store.LoadAllSummaries(root)
			if err != nil {
				return err
			}
			filtered := query.FilterTasks(all, types.QueryOptions{
				Type:     types.TaskType(filterType),
				Status:   types.Status(filterStatus),
				Priority: types.Priority(filterPriority),
			})
			sorted := query.SortByPriority(filtered)
			fmt.Fprintln(cmd.OutOrStdout(), renderTable(sorted))
			return nil
		},
	}

	cmd.Flags().StringVar(&filterType, "type", "", "Filter by type: epic, task, subtask, bug")
	cmd.Flags().StringVar(&filterStatus, "status", "", "Filter by status: todo, in-progress, backlog, done, cancelled")
	cmd.Flags().StringVar(&filterPriority, "priority", "", "Filter by priority: critical, high, medium, low")
	return cmd
}
