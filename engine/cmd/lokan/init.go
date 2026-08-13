package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/avressatelier/lokan/internal/id"
	"github.com/avressatelier/lokan/internal/store"
	"github.com/avressatelier/lokan/internal/types"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a lokan project in the current directory",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := os.Getwd()
			if err != nil {
				return err
			}

			// already initialized — report and bail
			configPath := id.ConfigPath(root)
			if _, err := os.Stat(configPath); err == nil {
				fmt.Fprintf(cmd.OutOrStdout(), "Already a lokan project. Config: %s\n", configPath)
				return nil
			}

			// write config with a fresh counter, scaffold the tasks dir
			if err := id.WriteConfig(root, types.LokanConfig{Counter: 0, Version: "1"}); err != nil {
				return err
			}
			// scaffold the single board file to match the frozen layout
			board := store.BoardPath(root)
			if _, err := os.Stat(board); os.IsNotExist(err) {
				if err := os.MkdirAll(filepath.Dir(board), 0o755); err != nil {
					return err
				}
				initial := "# Kanlo Board\n\n## Active\n\n## Archive\n"
				if err := os.WriteFile(board, []byte(initial), 0o644); err != nil {
					return err
				}
			}

			// report the created paths
			fmt.Fprintf(cmd.OutOrStdout(), "Initialized lokan project.\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  Config: %s\n", configPath)
			fmt.Fprintf(cmd.OutOrStdout(), "  Board:  %s\n", board)
			return nil
		},
	}
}
