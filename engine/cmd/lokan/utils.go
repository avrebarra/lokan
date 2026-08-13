package main

import (
	"errors"
	"strings"

	"github.com/avressatelier/lokan/internal/store"
	"github.com/avressatelier/lokan/internal/types"
)

// contains reports whether v is in list.
func contains[T comparable](list []T, v T) bool {
	for _, item := range list {
		if item == v {
			return true
		}
	}
	return false
}

// joinTypes renders the allowed task types as a comma-separated list.
func joinTypes() string {
	parts := make([]string, len(types.TaskTypes))
	for i, t := range types.TaskTypes {
		parts[i] = string(t)
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

// joinStatuses renders the configured lane ids as a comma-separated list.
func joinStatuses(board string) string {
	parts := make([]string, 0, len(types.Statuses))
	for _, s := range statusIDs(board) {
		parts = append(parts, string(s))
	}
	return strings.Join(parts, ", ")
}

// defaultStatus returns the first non-archived lane — the landing status for
// new tasks. Falls back to the first lane when every lane is archived.
func defaultStatus(board string) types.Status {
	cfg, err := store.ReadConfig(board)
	if err != nil {
		return types.StatusTodo
	}
	for _, s := range cfg.Statuses {
		if !s.Archived {
			return s.ID
		}
	}
	return cfg.Statuses[0].ID
}

// isArchivedStatus reports whether a status belongs in the Archive section,
// per the configured lane set. Unknown statuses default to active.
func isArchivedStatus(status types.Status, statuses []types.StatusDef) bool {
	for _, s := range statuses {
		if s.ID == status {
			return s.Archived
		}
	}
	return false
}

// joinPriorities renders the allowed priorities as a comma-separated list.
func joinPriorities() string {
	parts := make([]string, len(types.Priorities))
	for i, p := range types.Priorities {
		parts[i] = string(p)
	}
	return strings.Join(parts, ", ")
}

// allowedParents renders the allowed parent types for t, or "none".
func allowedParents(t types.TaskType) string {
	allowed := types.AllowedParents[t]
	if len(allowed) == 0 {
		return "none"
	}
	parts := make([]string, len(allowed))
	for i, a := range allowed {
		parts[i] = string(a)
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
