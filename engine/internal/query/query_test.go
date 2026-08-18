package query

import (
	"reflect"
	"testing"

	"github.com/avrebarra/lokan/internal/types"
)

func makeTask(id string, mut func(*types.TaskSummary)) types.TaskSummary {
	t := types.TaskSummary{
		TaskFrontmatter: types.TaskFrontmatter{
			ID:      id,
			Title:   "Task " + id,
			Status:  types.StatusTodo,
			Created: "2024-01-01",
			Updated: "2024-01-01",
		},
		FilePath: "/tasks/" + id + ".md",
	}
	if mut != nil {
		mut(&t)
	}
	return t
}

func fixtureTasks() []types.TaskSummary {
	return []types.TaskSummary{
		makeTask("1", func(t *types.TaskSummary) {
			t.Status = types.StatusInProgress
			t.Tags = []string{"frontend"}
		}),
		makeTask("2", func(t *types.TaskSummary) {
			t.Status = types.StatusTodo
			t.Tags = []string{"frontend", "auth"}
		}),
		makeTask("3", func(t *types.TaskSummary) {
			t.Status = types.StatusInProgress
		}),
		makeTask("4", func(t *types.TaskSummary) {
			t.Status = types.StatusTodo
		}),
		makeTask("5", func(t *types.TaskSummary) {
			t.Status = types.StatusDone
		}),
		makeTask("6", func(t *types.TaskSummary) {
			t.Status = types.StatusBacklog
			t.Tags = []string{"auth"}
		}),
		makeTask("7", func(t *types.TaskSummary) {
			t.Status = types.StatusCancelled
		}),
	}
}

func ids(tasks []types.TaskSummary) []string {
	out := make([]string, len(tasks))
	for i, t := range tasks {
		out[i] = t.ID
	}
	return out
}

// ---------------------------------------------------------------------------
// filterTasks
// ---------------------------------------------------------------------------

func TestFilterTasksByStatus(t *testing.T) {
	got := ids(FilterTasks(fixtureTasks(), types.QueryOptions{Status: types.StatusTodo}))
	if !reflect.DeepEqual(got, []string{"2", "4"}) {
		t.Fatalf("got %v", got)
	}
}

func TestFilterTasksBySingleTag(t *testing.T) {
	got := ids(FilterTasks(fixtureTasks(), types.QueryOptions{Tags: []string{"auth"}}))
	if !reflect.DeepEqual(got, []string{"2", "6"}) {
		t.Fatalf("got %v", got)
	}
}

func TestFilterTasksByMultipleTagsAND(t *testing.T) {
	got := ids(FilterTasks(fixtureTasks(), types.QueryOptions{Tags: []string{"frontend", "auth"}}))
	if !reflect.DeepEqual(got, []string{"2"}) {
		t.Fatalf("got %v", got)
	}
}

func TestFilterTasksMultipleFilters(t *testing.T) {
	got := ids(FilterTasks(fixtureTasks(), types.QueryOptions{Status: types.StatusInProgress, Tags: []string{"frontend"}}))
	if !reflect.DeepEqual(got, []string{"1"}) {
		t.Fatalf("got %v", got)
	}
}

func TestFilterTasksNoMatch(t *testing.T) {
	if got := FilterTasks(fixtureTasks(), types.QueryOptions{Status: types.StatusDone, Tags: []string{"auth"}}); len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}

func TestFilterTasksNoOptions(t *testing.T) {
	if got := FilterTasks(fixtureTasks(), types.QueryOptions{}); len(got) != 7 {
		t.Fatalf("got %d", len(got))
	}
}
