package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type rootKey struct{}

// findRoot walks up from dir looking for a directory with .lokan/config.json.
func findRoot(dir string) (string, error) {
	// walk up until the project config is found or the fs root is hit
	for {
		if _, err := os.Stat(filepath.Join(dir, ".lokan", "config.json")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", ErrNotInProject
		}
		dir = parent
	}
}

// newRootCmd assembles the command tree. The persistent pre-run enforces the
// project requirement for every command except init and help.
func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "lokan",
		Short:         "Markdown task manager for Claude Code",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// init and help work without a project; resolve root for the rest
			switch cmd.Name() {
			case "init", "help":
				return nil
			}
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			projectRoot, err := findRoot(cwd)
			if err != nil {
				return err
			}
			cmd.SetContext(context.WithValue(cmd.Context(), rootKey{}, projectRoot))
			return nil
		},
	}
	root.CompletionOptions.DisableDefaultCmd = true

	// register the command tree
	root.AddCommand(newInitCmd())
	root.AddCommand(newCreateCmd())
	root.AddCommand(newGetCmd())
	root.AddCommand(newListCmd())
	root.AddCommand(newEditCmd())
	root.AddCommand(newSubtasksCmd())
	root.AddCommand(newUICmd())
	return root
}

// cmdRoot returns the project root stored by the persistent pre-run.
func cmdRoot(cmd *cobra.Command) string {
	return cmd.Context().Value(rootKey{}).(string)
}
