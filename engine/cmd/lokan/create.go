package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/avressatelier/lokan/internal/id"
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
			&cli.StringFlag{Name: "type", Aliases: []string{"t"}, Value: string(types.TypeTask), Usage: "Task type: epic, task, subtask, bug"},
			&cli.StringFlag{Name: "priority", Value: string(types.PriorityMedium), Usage: "Priority: critical, high, medium, low"},
			&cli.StringFlag{Name: "parent", Usage: "Parent task ID"},
			&cli.StringSliceFlag{Name: "tag", Usage: "Tag to add (repeatable)"},
		},
		Action: func(c *cli.Context) error {
			// validate the positional title
			if err := requireArgs(c, 1); err != nil {
				return err
			}
			title := c.Args().First()
			taskType := c.String("type")
			priority := c.String("priority")
			parent := c.String("parent")
			tags := c.StringSlice("tag")

			// every command except init/help needs a project
			root, err := requireProject(c)
			if err != nil {
				return err
			}

			// validate type and priority enums
			t := types.TaskType(taskType)
			if !contains(types.TaskTypes, t) {
				return cliErrorf("Invalid type %q. Must be one of: %s", taskType, joinTypes())
			}
			p := types.Priority(priority)
			if !contains(types.Priorities, p) {
				return cliErrorf("Invalid priority %q. Must be one of: %s", priority, joinPriorities())
			}

			// resolve and validate the parent, if given
			if parent != "" {
				parentSummary, err := store.FindByID(root, parent)
				if err != nil {
					return cliErrorf("Parent task not found: %s", parent)
				}
				if !contains(types.AllowedParents[t], parentSummary.Type) {
					return cliErrorf("Cannot create %s under %s (%s). Allowed parents: %s", t, parentSummary.Type, parent, allowedParents(t))
				}
			}

			// allocate id/filename from the counter and write the task file
			counter, err := id.NextCounter(root)
			if err != nil {
				return err
			}
			taskID := id.GenerateID(counter)
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
			}, "")
			if err != nil {
				return err
			}
			// print the created task with a project-relative path
			rel, err := filepath.Rel(root, store.BoardPath(root))
			if err != nil {
				rel = store.BoardPath(root)
			}
			fmt.Fprintf(c.App.Writer, "Created %s → %s\n", task.ID, rel)
			if parent != "" {
				fmt.Fprintf(c.App.Writer, "  Parent: %s\n", parent)
			}
			return nil
		},
	}
}
