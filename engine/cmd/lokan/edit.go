package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/avressatelier/lokan/internal/id"
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

			// rename the file when the title changed, then persist
			if titleChanged {
				// Issue 2: derive counter from the existing filename basename,
				// not from the user-supplied id
				counter, err := counterFromFilename(task.FilePath)
				if err != nil {
					return err
				}
				newFilename := id.GenerateFilename(task.Type, counter, task.Title)
				task, err = store.RenameTask(task, newFilename)
				if err != nil {
					return err
				}
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

// counterFromFilename parses the numeric counter out of a task filename
// basename, e.g. "task-12-old-title.md" → 12.
func counterFromFilename(filePath string) (int, error) {
	base := filepath.Base(filePath)
	base = strings.TrimSuffix(base, filepath.Ext(base))
	parts := strings.Split(base, "-")
	if len(parts) < 2 {
		return 0, cliErrorf("Cannot rename: could not parse counter from filename %q", base)
	}
	counter, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, cliErrorf("Cannot rename: could not parse counter from filename %q", base)
	}
	return counter, nil
}

func parentLabel(p string) string {
	if p == "" {
		return "none"
	}
	return p
}
