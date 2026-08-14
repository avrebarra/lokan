package query

import (
	"reflect"
	"testing"

	"github.com/avrebarra/lokan/internal/types"
)

func makeTask(id string, mut func(*types.TaskSummary)) types.TaskSummary {
	t := types.TaskSummary{
		TaskFrontmatter: types.TaskFrontmatter{
			ID:       id,
			Title:    "Task " + id,
			Type:     types.TypeTask,
			Status:   types.StatusTodo,
			Priority: types.PriorityMedium,
			Created:  "2024-01-01",
			Updated:  "2024-01-01",
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
			t.Type = types.TypeEpic
			t.Status = types.StatusInProgress
			t.Priority = types.PriorityCritical
			t.Tags = []string{"frontend"}
		}),
		makeTask("2", func(t *types.TaskSummary) {
			t.Status = types.StatusTodo
			t.Priority = types.PriorityHigh
			t.Parent = "1"
			t.Tags = []string{"frontend", "auth"}
		}),
		makeTask("3", func(t *types.TaskSummary) {
			t.Status = types.StatusInProgress
			t.Parent = "1"
		}),
		makeTask("4", func(t *types.TaskSummary) {
			t.Type = types.TypeSubtask
			t.Priority = types.PriorityLow
			t.Parent = "2"
		}),
		makeTask("5", func(t *types.TaskSummary) {
			t.Type = types.TypeSubtask
			t.Status = types.StatusDone
			t.Priority = types.PriorityLow
			t.Parent = "2"
		}),
		makeTask("6", func(t *types.TaskSummary) {
			t.Type = types.TypeBug
			t.Status = types.StatusBacklog
			t.Priority = types.PriorityHigh
			t.Tags = []string{"auth"}
		}),
		makeTask("7", func(t *types.TaskSummary) {
			t.Status = types.StatusCancelled
			t.Priority = types.PriorityLow
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

func TestFilterTasksByType(t *testing.T) {
	got := ids(FilterTasks(fixtureTasks(), types.QueryOptions{Type: types.TypeSubtask}))
	if !reflect.DeepEqual(got, []string{"4", "5"}) {
		t.Fatalf("got %v", got)
	}
}

func TestFilterTasksByStatus(t *testing.T) {
	got := ids(FilterTasks(fixtureTasks(), types.QueryOptions{Status: types.StatusTodo}))
	if !reflect.DeepEqual(got, []string{"2", "4"}) {
		t.Fatalf("got %v", got)
	}
}

func TestFilterTasksByPriority(t *testing.T) {
	got := ids(FilterTasks(fixtureTasks(), types.QueryOptions{Priority: types.PriorityHigh}))
	if !reflect.DeepEqual(got, []string{"2", "6"}) {
		t.Fatalf("got %v", got)
	}
}

func TestFilterTasksByParent(t *testing.T) {
	got := ids(FilterTasks(fixtureTasks(), types.QueryOptions{Parent: "2"}))
	if !reflect.DeepEqual(got, []string{"4", "5"}) {
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
	got := ids(FilterTasks(fixtureTasks(), types.QueryOptions{Type: types.TypeTask, Status: types.StatusInProgress}))
	if !reflect.DeepEqual(got, []string{"3"}) {
		t.Fatalf("got %v", got)
	}
}

func TestFilterTasksNoMatch(t *testing.T) {
	if got := FilterTasks(fixtureTasks(), types.QueryOptions{Type: types.TypeBug, Status: types.StatusDone}); len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}

func TestFilterTasksNoOptions(t *testing.T) {
	if got := FilterTasks(fixtureTasks(), types.QueryOptions{}); len(got) != 7 {
		t.Fatalf("got %d", len(got))
	}
}

// ---------------------------------------------------------------------------
// getChildren
// ---------------------------------------------------------------------------

func TestGetChildrenDirectOnly(t *testing.T) {
	got := ids(GetChildren(fixtureTasks(), "2"))
	if !reflect.DeepEqual(got, []string{"4", "5"}) {
		t.Fatalf("got %v", got)
	}
}

func TestGetChildrenExcludesGrandchildren(t *testing.T) {
	got := ids(GetChildren(fixtureTasks(), "1"))
	if !reflect.DeepEqual(got, []string{"2", "3"}) {
		t.Fatalf("got %v", got)
	}
}

func TestGetChildrenLeaf(t *testing.T) {
	if got := GetChildren(fixtureTasks(), "4"); len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}

// ---------------------------------------------------------------------------
// sortByPriority
// ---------------------------------------------------------------------------

func TestSortByPriorityOrder(t *testing.T) {
	input := []types.TaskSummary{
		makeTask("a", func(t *types.TaskSummary) { t.Priority = types.PriorityLow }),
		makeTask("b", func(t *types.TaskSummary) { t.Priority = types.PriorityCritical }),
		makeTask("c", func(t *types.TaskSummary) { t.Priority = types.PriorityMedium }),
		makeTask("d", func(t *types.TaskSummary) { t.Priority = types.PriorityHigh }),
	}
	got := ids(SortByPriority(input))
	if !reflect.DeepEqual(got, []string{"b", "d", "c", "a"}) {
		t.Fatalf("got %v", got)
	}
}

func TestSortByPriorityDoesNotMutate(t *testing.T) {
	input := []types.TaskSummary{
		makeTask("x", func(t *types.TaskSummary) { t.Priority = types.PriorityLow }),
		makeTask("y", func(t *types.TaskSummary) { t.Priority = types.PriorityCritical }),
	}
	original := append([]types.TaskSummary(nil), input...)
	SortByPriority(input)
	if !reflect.DeepEqual(ids(input), ids(original)) {
		t.Fatalf("input mutated: %v", ids(input))
	}
}

func TestSortByPriorityTieBreaksByCreated(t *testing.T) {
	input := []types.TaskSummary{
		makeTask("task-b", func(t *types.TaskSummary) { t.Priority = types.PriorityHigh; t.Created = "2024-03-01" }),
		makeTask("task-a", func(t *types.TaskSummary) { t.Priority = types.PriorityHigh; t.Created = "2024-01-01" }),
	}
	got := ids(SortByPriority(input))
	if !reflect.DeepEqual(got, []string{"task-a", "task-b"}) {
		t.Fatalf("got %v", got)
	}
}
