package query

import (
	"github.com/avrebarra/lokan/internal/types"
)

// FilterTasks applies QueryOptions with AND semantics; zero-valued options
// are ignored.
func FilterTasks(tasks []types.TaskSummary, opts types.QueryOptions) []types.TaskSummary {
	var out []types.TaskSummary
	for _, t := range tasks {
		// skip on any active filter mismatch (AND semantics)
		if opts.Status != "" && t.Status != opts.Status {
			continue
		}
		if len(opts.Tags) > 0 && !hasAllTags(t.Tags, opts.Tags) {
			continue
		}
		out = append(out, t)
	}
	return out
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
