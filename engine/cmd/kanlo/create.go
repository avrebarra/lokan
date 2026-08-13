package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/avressatelier/kanlo/internal/id"
	"github.com/avressatelier/kanlo/internal/store"
	"github.com/avressatelier/kanlo/internal/types"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	var taskType, priority, parent string
	var tags []string

	cmd := &cobra.Command{
		Use:   "create <title>",
		Short: "Create a new task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			title := args[0]
			t := types.TaskType(taskType)
			if !contains(types.TaskTypes, t) {
				return cliErrorf("Invalid type %q. Must be one of: %s", taskType, joinTypes())
			}
			p := types.Priority(priority)
			if !contains(types.Priorities, p) {
				return cliErrorf("Invalid priority %q. Must be one of: %s", priority, joinPriorities())
			}

			root := cmdRoot(cmd)
			if parent != "" {
				parentSummary, err := store.FindByID(root, parent)
				if err != nil {
					return cliErrorf("Parent task not found: %s", parent)
				}
				if !contains(types.AllowedParents[t], parentSummary.Type) {
					return cliErrorf("Cannot create %s under %s (%s). Allowed parents: %s", t, parentSummary.Type, parent, allowedParents(t))
				}
			}

			counter, err := id.NextCounter(root)
			if err != nil {
				return err
			}
			taskID := id.GenerateID(t, counter)
			filename := id.GenerateFilename(t, counter, title)
			today := time.Now().UTC().Format("2006-01-02")

			task, err := store.CreateTask(root, types.TaskFrontmatter{
				ID:       taskID,
				Title:    title,
				Type:     t,
				Status:   types.StatusTodo,
				Priority: p,
				Created:  today,
				Updated:  today,
				Parent:   parent,
				Tags:     tags,
			}, filename, "")
			if err != nil {
				return err
			}
			rel, err := filepath.Rel(root, task.FilePath)
			if err != nil {
				rel = task.FilePath
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Created %s → %s\n", task.ID, rel)
			if parent != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  Parent: %s\n", parent)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&taskType, "type", "t", string(types.TypeTask), "Task type: epic, task, subtask, bug")
	cmd.Flags().StringVar(&priority, "priority", string(types.PriorityMedium), "Priority: critical, high, medium, low")
	cmd.Flags().StringVar(&parent, "parent", "", "Parent task ID")
	cmd.Flags().StringArrayVar(&tags, "tag", nil, "Tag to add (repeatable)")
	return cmd
}
