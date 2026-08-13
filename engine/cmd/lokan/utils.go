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

// joinStatuses renders the allowed statuses as a comma-separated list.
func joinStatuses() string {
	parts := make([]string, len(types.Statuses))
	for i, s := range types.Statuses {
		parts[i] = string(s)
	}
	return strings.Join(parts, ", ")
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
