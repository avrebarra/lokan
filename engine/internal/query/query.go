package query

import (
	"sort"

	"github.com/avressatelier/lokan/internal/types"
)

// FilterTasks applies QueryOptions with AND semantics; zero-valued options
// are ignored.
func FilterTasks(tasks []types.TaskSummary, opts types.QueryOptions) []types.TaskSummary {
	var out []types.TaskSummary
	for _, t := range tasks {
		// skip on any active filter mismatch (AND semantics)
		if opts.Type != "" && t.Type != opts.Type {
			continue
		}
		if opts.Status != "" && t.Status != opts.Status {
			continue
		}
		if opts.Priority != "" && t.Priority != opts.Priority {
			continue
		}
		if opts.Parent != "" && t.Parent != opts.Parent {
			continue
		}
		if len(opts.Tags) > 0 && !hasAllTags(t.Tags, opts.Tags) {
			continue
		}
		out = append(out, t)
	}
	return out
}

// GetChildren returns direct children (tasks whose parent is parentID).
func GetChildren(tasks []types.TaskSummary, parentID string) []types.TaskSummary {
	var out []types.TaskSummary
	for _, t := range tasks {
		if t.Parent == parentID {
			out = append(out, t)
		}
	}
	return out
}

var priorityOrder = map[types.Priority]int{
	types.PriorityCritical: 0,
	types.PriorityHigh:     1,
	types.PriorityMedium:   2,
	types.PriorityLow:      3,
}

// SortByPriority sorts by priority (critical first), tie-breaking on created.
// The input slice is not mutated.
func SortByPriority(tasks []types.TaskSummary) []types.TaskSummary {
	// copy to avoid mutating the caller's slice
	out := make([]types.TaskSummary, len(tasks))
	copy(out, tasks)

	// sort by priority rank, then by created date
	sort.Slice(out, func(i, j int) bool {
		pi := priorityOf(out[i].Priority)
		pj := priorityOf(out[j].Priority)
		if pi != pj {
			return pi < pj
		}
		return out[i].Created < out[j].Created
	})
	return out
}

func priorityOf(p types.Priority) int {
	if order, ok := priorityOrder[p]; ok {
		return order
	}
	return 99
}

func hasAllTags(taskTags []string, want []string) bool {
	for _, tag := range want {
		found := false
		for _, have := range taskTags {
			if have == tag {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
