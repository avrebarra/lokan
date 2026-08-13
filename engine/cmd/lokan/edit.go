package main

import (
	"fmt"

	"github.com/avressatelier/lokan/internal/store"
	"github.com/avressatelier/lokan/internal/types"
	"github.com/spf13/cobra"
)

func newEditCmd() *cobra.Command {
	var newStatus, newPriority, newTitle, newParent string

	cmd := &cobra.Command{
		Use:   "edit <id>",
		Short: "Update task fields",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]

			// validate flag enums before touching the store
			if newStatus != "" && !contains(types.Statuses, types.Status(newStatus)) {
				return cliErrorf("Invalid status %q. Must be one of: %s", newStatus, joinStatuses())
			}
			if newPriority != "" && !contains(types.Priorities, types.Priority(newPriority)) {
				return cliErrorf("Invalid priority %q. Must be one of: %s", newPriority, joinPriorities())
			}

			// load the target task
			root := cmdRoot(cmd)
			summary, err := store.FindByID(root, taskID)
			if err != nil {
				return notFoundError(taskID, err)
			}
			task, err := store.LoadTask(summary.FilePath)
			if err != nil {
				return err
			}
			changes := []string{}

			// apply status, priority, parent, title edits
			if newStatus != "" && types.Status(newStatus) != task.Status {
				changes = append(changes, fmt.Sprintf("status: %s → %s", task.Status, newStatus))
				task.Status = types.Status(newStatus)
			}
			if newPriority != "" && types.Priority(newPriority) != task.Priority {
				changes = append(changes, fmt.Sprintf("priority: %s → %s", task.Priority, newPriority))
				task.Priority = types.Priority(newPriority)
			}
			if cmd.Flags().Changed("parent") {
				old := task.Parent
				if old == "" {
					old = "none"
				}
				changes = append(changes, fmt.Sprintf("parent: %s → %s", old, parentLabel(newParent)))
				task.Parent = newParent
			}
			titleChanged := newTitle != "" && newTitle != task.Title
			if titleChanged {
				changes = append(changes, fmt.Sprintf("title: %q → %q", task.Title, newTitle))
				task.Title = newTitle
			}

			// nothing to do — report and bail early
			if len(changes) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No changes for %s.\n", taskID)
				return nil
			}

			if err := store.WriteTask(task); err != nil {
				return err
			}

			// print the change list
			fmt.Fprintf(cmd.OutOrStdout(), "Updated %s\n", taskID)
			for _, c := range changes {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", c)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&newStatus, "status", "", "New status: todo, in-progress, backlog, done, cancelled")
	cmd.Flags().StringVar(&newPriority, "priority", "", "New priority: critical, high, medium, low")
	cmd.Flags().StringVar(&newTitle, "title", "", "New title")
	cmd.Flags().StringVar(&newParent, "parent", "", "Set parent task ID (empty clears)")
	return cmd
}

func parentLabel(p string) string {
	if p == "" {
		return "none"
	}
	return p
}
