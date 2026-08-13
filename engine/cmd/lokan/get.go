package main

import (
	"fmt"

	"github.com/avressatelier/lokan/internal/store"
	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Show a task by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := cmdRoot(cmd)
			summary, err := store.FindByID(root, args[0])
			if err != nil {
				return notFoundError(args[0], err)
			}
			task, err := store.LoadTask(summary.FilePath)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), renderTaskDetail(task))
			return nil
		},
	}
}
