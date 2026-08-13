package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/avressatelier/kanlo/internal/id"
	"github.com/avressatelier/kanlo/internal/store"
	"github.com/avressatelier/kanlo/internal/types"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a kanlo project in the current directory",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := os.Getwd()
			if err != nil {
				return err
			}
			configPath := id.ConfigPath(root)
			if _, err := os.Stat(configPath); err == nil {
				fmt.Fprintf(cmd.OutOrStdout(), "Already a kanlo project. Config: %s\n", configPath)
				return nil
			}
			if err := id.WriteConfig(root, types.NodConfig{Counter: 0, Version: "1"}); err != nil {
				return err
			}
			// scaffold the tasks dir to match the frozen layout
			tasksDir := store.TasksDir(root)
			if err := os.MkdirAll(tasksDir, 0o755); err != nil {
				return err
			}
			gitkeep := filepath.Join(tasksDir, ".gitkeep")
			if _, err := os.Stat(gitkeep); os.IsNotExist(err) {
				if err := os.WriteFile(gitkeep, []byte{}, 0o644); err != nil {
					return err
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Initialized kanlo project.\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  Config: %s\n", configPath)
			fmt.Fprintf(cmd.OutOrStdout(), "  Tasks:  %s\n", tasksDir)
			return nil
		},
	}
}
