package query

import (
	"reflect"
	"testing"

	"github.com/avressatelier/lokan/internal/types"
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
// getDescendants
// ---------------------------------------------------------------------------

func TestGetDescendantsRecursesAllLevels(t *testing.T) {
	got := ids(GetDescendants(fixtureTasks(), "1"))
	want := []string{"2", "4", "5", "3"}
	if len(got) != 4 {
		t.Fatalf("got %v", got)
	}
	for _, w := range want {
		found := false
		for _, g := range got {
			if g == w {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("missing %s in %v", w, got)
		}
	}
}

func TestGetDescendantsDirectChildren(t *testing.T) {
	got := ids(GetDescendants(fixtureTasks(), "2"))
	if !reflect.DeepEqual(got, []string{"4", "5"}) {
		t.Fatalf("got %v", got)
	}
}

// ---------------------------------------------------------------------------
// getAvailable / getBacklog
// ---------------------------------------------------------------------------

func TestGetAvailableOnlyTodoAndInProgress(t *testing.T) {
	got := GetAvailable(fixtureTasks(), nil)
	statuses := map[types.Status]bool{}
	for _, t := range got {
		statuses[t.Status] = true
	}
	if !statuses[types.StatusTodo] || !statuses[types.StatusInProgress] || len(statuses) != 2 {
		t.Fatalf("statuses = %v", statuses)
	}
}

func TestGetAvailableExcludesBacklog(t *testing.T) {
	for _, task := range GetAvailable(fixtureTasks(), nil) {
		if task.Status == types.StatusBacklog {
			t.Fatal("backlog task leaked into available")
		}
	}
}

func TestGetAvailableExcludesDone(t *testing.T) {
	for _, task := range GetAvailable(fixtureTasks(), nil) {
		if task.Status == types.StatusDone {
			t.Fatal("done task leaked into available")
		}
	}
}

func TestGetAvailableExcludesCancelled(t *testing.T) {
	for _, task := range GetAvailable(fixtureTasks(), nil) {
		if task.Status == types.StatusCancelled {
			t.Fatal("cancelled task leaked into available")
		}
	}
}

func TestGetAvailableWithOptions(t *testing.T) {
	for _, task := range GetAvailable(fixtureTasks(), &types.QueryOptions{Type: types.TypeTask}) {
		if task.Type != types.TypeTask {
			t.Fatalf("type = %v", task.Type)
		}
	}
}

func TestGetBacklogOnlyBacklog(t *testing.T) {
	got := GetBacklog(fixtureTasks(), nil)
	if len(got) != 1 || got[0].ID != "6" {
		t.Fatalf("got %v", ids(got))
	}
}

// ---------------------------------------------------------------------------
// buildTree
// ---------------------------------------------------------------------------

func TestBuildTreeUnknownRoot(t *testing.T) {
	if tree := BuildTree(fixtureTasks(), "99"); tree != nil {
		t.Fatalf("expected nil, got %+v", tree)
	}
}

func TestBuildTreeLeaf(t *testing.T) {
	tasks := []types.TaskSummary{makeTask("1", func(t *types.TaskSummary) { t.Type = types.TypeEpic })}
	tree := BuildTree(tasks, "1")
	if tree == nil || tree.Task.ID != "1" || len(tree.Children) != 0 {
		t.Fatalf("tree = %+v", tree)
	}
}

func TestBuildTreeTwoLevel(t *testing.T) {
	tasks := []types.TaskSummary{
		makeTask("1", func(t *types.TaskSummary) { t.Type = types.TypeEpic }),
		makeTask("2", func(t *types.TaskSummary) { t.Parent = "1" }),
		makeTask("3", func(t *types.TaskSummary) { t.Parent = "1" }),
	}
	tree := BuildTree(tasks, "1")
	if len(tree.Children) != 2 {
		t.Fatalf("children = %d, want 2", len(tree.Children))
	}
	childIDs := []string{tree.Children[0].Task.ID, tree.Children[1].Task.ID}
	if !contains(childIDs, "2") || !contains(childIDs, "3") {
		t.Fatalf("child ids = %v", childIDs)
	}
}

func TestBuildTreeThreeLevel(t *testing.T) {
	tasks := []types.TaskSummary{
		makeTask("1", func(t *types.TaskSummary) { t.Type = types.TypeEpic }),
		makeTask("2", func(t *types.TaskSummary) { t.Parent = "1" }),
		makeTask("3", func(t *types.TaskSummary) { t.Type = types.TypeSubtask; t.Parent = "2" }),
	}
	tree := BuildTree(tasks, "1")
	if len(tree.Children) != 1 || tree.Children[0].Task.ID != "2" {
		t.Fatalf("tree = %+v", tree)
	}
	grand := tree.Children[0].Children
	if len(grand) != 1 || grand[0].Task.ID != "3" {
		t.Fatalf("grandchildren = %+v", grand)
	}
}

func TestBuildTreeCycleSafe(t *testing.T) {
	tasks := []types.TaskSummary{
		makeTask("2", func(t *types.TaskSummary) { t.Parent = "3" }),
		makeTask("3", func(t *types.TaskSummary) { t.Parent = "2" }),
	}
	tree := BuildTree(tasks, "2")
	if tree == nil {
		t.Fatal("expected non-nil tree")
	}
	if len(tree.Children) != 1 || tree.Children[0].Task.ID != "3" {
		t.Fatalf("tree = %+v", tree)
	}
	if len(tree.Children[0].Children) != 0 {
		t.Fatalf("cycle not broken: %+v", tree)
	}
}

func TestBuildTreeSelfCycleSafe(t *testing.T) {
	tasks := []types.TaskSummary{
		makeTask("1", func(t *types.TaskSummary) { t.Parent = "1" }),
	}
	tree := BuildTree(tasks, "1")
	if tree == nil || len(tree.Children) != 0 {
		t.Fatalf("tree = %+v", tree)
	}
}

func TestBuildTreeDiamondDoesNotLoseNodes(t *testing.T) {
	tasks := []types.TaskSummary{
		makeTask("1", func(t *types.TaskSummary) { t.Type = types.TypeEpic }),
		makeTask("2", func(t *types.TaskSummary) { t.Parent = "1" }),
		makeTask("3", func(t *types.TaskSummary) { t.Parent = "1" }),
		makeTask("4", func(t *types.TaskSummary) { t.Type = types.TypeSubtask; t.Parent = "2" }),
		makeTask("5", func(t *types.TaskSummary) { t.Type = types.TypeSubtask; t.Parent = "3" }),
	}
	tree := BuildTree(tasks, "1")
	seen := map[string]bool{}
	var walk func(n *types.TreeNode)
	walk = func(n *types.TreeNode) {
		seen[n.Task.ID] = true
		for i := range n.Children {
			walk(&n.Children[i])
		}
	}
	walk(tree)
	for _, id := range []string{"1", "2", "3", "4", "5"} {
		if !seen[id] {
			t.Fatalf("lost node %s in tree", id)
		}
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

func contains(hay []string, needle string) bool {
	for _, s := range hay {
		if s == needle {
			return true
		}
	}
	return false
}
