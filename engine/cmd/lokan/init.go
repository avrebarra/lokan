package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/avressatelier/lokan/internal/id"
	"github.com/avressatelier/lokan/internal/store"
	"github.com/avressatelier/lokan/internal/types"
	"github.com/urfave/cli/v2"
)

func newInitCmd() *cli.Command {
	return &cli.Command{
		Name:         "init",
		Usage:        "Initialize a lokan project in the current directory",
		OnUsageError: quietUsageError,
		Action: func(c *cli.Context) error {
			// only a bare invocation is valid
			if err := requireArgs(c, 0); err != nil {
				return err
			}
			root, err := os.Getwd()
			if err != nil {
				return err
			}
			out := c.App.Writer

			// already initialized — report and bail
			configPath := id.ConfigPath(root)
			if _, err := os.Stat(configPath); err == nil {
				fmt.Fprintf(out, "Already a lokan project. Config: %s\n", configPath)
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
			fmt.Fprintf(out, "Initialized lokan project.\n")
			fmt.Fprintf(out, "  Config: %s\n", configPath)
			fmt.Fprintf(out, "  Board:  %s\n", board)
			return nil
		},
	}
}
