package main

import (
	"errors"
	"strings"

	"github.com/avrebarra/lokan/internal/store"
	"github.com/avrebarra/lokan/internal/types"
)

// joinStatuses renders the configured lane ids as a comma-separated list.
func joinStatuses(board string) string {
	parts := make([]string, 0, len(types.Statuses))
	for _, s := range statusIDs(board) {
		parts = append(parts, string(s))
	}
	return strings.Join(parts, ", ")
}

// statusIDs resolves the configured lane ids for a board, falling back to
// the built-in enum when the config is unreadable.
func statusIDs(board string) []types.Status {
	cfg, err := store.ReadConfig(board)
	if err != nil {
		return types.Statuses
	}
	ids := make([]types.Status, len(cfg.Statuses))
	for i, s := range cfg.Statuses {
		ids[i] = s.ID
	}
	return ids
}

// joinPriorities renders the allowed priorities as a comma-separated list.
func joinPriorities() string {
	parts := make([]string, len(types.Priorities))
	for i, p := range types.Priorities {
		parts[i] = string(p)
	}
	return strings.Join(parts, ", ")
}

// notFoundError converts a store not-found error into a NotFoundError.
func notFoundError(id string, err error) error {
	if errors.Is(err, store.ErrNotFound) {
		return &NotFoundError{ID: id}
	}
	return err
}
