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

// GetDescendants returns all tasks nested under rootID, recursively.
func GetDescendants(tasks []types.TaskSummary, rootID string) []types.TaskSummary {
	var out []types.TaskSummary
	for _, child := range GetChildren(tasks, rootID) {
		out = append(out, child)
		out = append(out, GetDescendants(tasks, child.ID)...)
	}
	return out
}

// GetAvailable returns todo and in-progress tasks, optionally filtered.
func GetAvailable(tasks []types.TaskSummary, opts *types.QueryOptions) []types.TaskSummary {
	return statusSlice(tasks, opts, types.StatusTodo, types.StatusInProgress)
}

// GetBacklog returns backlog tasks, optionally filtered.
func GetBacklog(tasks []types.TaskSummary, opts *types.QueryOptions) []types.TaskSummary {
	return statusSlice(tasks, opts, types.StatusBacklog)
}

func statusSlice(tasks []types.TaskSummary, opts *types.QueryOptions, statuses ...types.Status) []types.TaskSummary {
	included := make(map[types.Status]bool, len(statuses))
	for _, s := range statuses {
		included[s] = true
	}
	base := make([]types.TaskSummary, 0, len(tasks))
	for _, t := range tasks {
		if included[t.Status] {
			base = append(base, t)
		}
	}
	if opts == nil {
		return base
	}
	return FilterTasks(base, *opts)
}

// BuildTree assembles the task hierarchy under rootID, or nil if rootID is
// unknown. Cycle-safe: any id already on the current path is skipped, so
// corrupt parent graphs cannot infinite-loop (Issue 4).
func BuildTree(tasks []types.TaskSummary, rootID string) *types.TreeNode {
	var root *types.TaskSummary
	for i := range tasks {
		if tasks[i].ID == rootID {
			root = &tasks[i]
			break
		}
	}
	if root == nil {
		return nil
	}
	node := buildNode(tasks, *root, map[string]bool{rootID: true})
	return &node
}

func buildNode(tasks []types.TaskSummary, task types.TaskSummary, seen map[string]bool) types.TreeNode {
	node := types.TreeNode{Task: task}
	for _, child := range GetChildren(tasks, task.ID) {
		if seen[child.ID] {
			continue
		}
		seen[child.ID] = true
		node.Children = append(node.Children, buildNode(tasks, child, seen))
		delete(seen, child.ID)
	}
	return node
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
	out := make([]types.TaskSummary, len(tasks))
	copy(out, tasks)
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
