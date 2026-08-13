package store

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/avressatelier/lokan/internal/types"
)

func makeFrontmatter(overrides map[string]interface{}) types.TaskFrontmatter {
	fm := types.TaskFrontmatter{
		ID:       "task-1",
		Title:    "Test Task",
		Type:     types.TypeTask,
		Status:   types.StatusTodo,
		Priority: types.PriorityMedium,
		Created:  "2024-01-01",
		Updated:  "2024-01-01",
	}
	if v, ok := overrides["id"].(string); ok {
		fm.ID = v
	}
	if v, ok := overrides["title"].(string); ok {
		fm.Title = v
	}
	if v, ok := overrides["type"].(types.TaskType); ok {
		fm.Type = v
	}
	if v, ok := overrides["status"].(types.Status); ok {
		fm.Status = v
	}
	if v, ok := overrides["priority"].(types.Priority); ok {
		fm.Priority = v
	}
	if v, ok := overrides["parent"].(string); ok {
		fm.Parent = v
	}
	if v, ok := overrides["tags"].([]string); ok {
		fm.Tags = v
	}
	if v, ok := overrides["updated"].(string); ok {
		fm.Updated = v
	}
	return fm
}

func newTempRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(TasksDir(root), 0o755); err != nil {
		t.Fatal(err)
	}
	return root
}

// ---------------------------------------------------------------------------
// createTask
// ---------------------------------------------------------------------------

func TestCreateTaskCreatesFile(t *testing.T) {
	root := newTempRoot(t)
	fm := makeFrontmatter(nil)
	if _, err := CreateTask(root, fm, "task-1-test-task.md", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(TasksDir(root), "task-1-test-task.md")); err != nil {
		t.Fatalf("task file missing: %v", err)
	}
}

func TestCreateTaskReturnsFields(t *testing.T) {
	root := newTempRoot(t)
	fm := makeFrontmatter(map[string]interface{}{"id": "epic-5", "title": "Big Epic", "type": types.TypeEpic})
	task, err := CreateTask(root, fm, "epic-5-big-epic.md", "")
	if err != nil {
		t.Fatal(err)
	}
	if task.ID != "epic-5" || task.Title != "Big Epic" {
		t.Fatalf("task = %+v", task)
	}
	if task.FilePath != filepath.Join(TasksDir(root), "epic-5-big-epic.md") {
		t.Fatalf("filePath = %q", task.FilePath)
	}
}

// ---------------------------------------------------------------------------
// loadTask round-trip
// ---------------------------------------------------------------------------

func TestLoadTaskRoundTrip(t *testing.T) {
	root := newTempRoot(t)
	fm := makeFrontmatter(map[string]interface{}{"id": "task-42", "title": "Round-trip Task", "priority": types.PriorityHigh})
	created, err := CreateTask(root, fm, "task-42-round-trip-task.md", "")
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadTask(created.FilePath)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.ID != fm.ID || loaded.Title != fm.Title || loaded.Type != fm.Type ||
		loaded.Status != fm.Status || loaded.Priority != fm.Priority {
		t.Fatalf("loaded = %+v, want frontmatter %+v", loaded, fm)
	}
	if !strings.Contains(loaded.Body, "## Notes") || !strings.Contains(loaded.Body, "## Work Log") {
		t.Fatalf("body sections missing: %q", loaded.Body)
	}
	if !strings.HasPrefix(loaded.Body, "# Round-trip Task\n") {
		t.Fatalf("body missing title heading: %q", loaded.Body)
	}
}

// ---------------------------------------------------------------------------
// loadAllSummaries
// ---------------------------------------------------------------------------

func TestLoadAllSummariesReturnsAll(t *testing.T) {
	root := newTempRoot(t)
	mustCreate(t, root, "task-1-first.md", makeFrontmatter(map[string]interface{}{"id": "task-1", "title": "First"}))
	mustCreate(t, root, "task-2-second.md", makeFrontmatter(map[string]interface{}{"id": "task-2", "title": "Second"}))

	summaries, err := LoadAllSummaries(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 2 {
		t.Fatalf("len = %d, want 2", len(summaries))
	}
	ids := []string{summaries[0].ID, summaries[1].ID}
	if (ids[0] != "task-1" || ids[1] != "task-2") && (ids[0] != "task-2" || ids[1] != "task-1") {
		t.Fatalf("ids = %v", ids)
	}
}

func TestLoadAllSummariesSkipsInvalid(t *testing.T) {
	root := newTempRoot(t)
	mustCreate(t, root, "task-1-valid.md", makeFrontmatter(map[string]interface{}{"id": "task-1", "title": "Valid"}))
	if err := os.WriteFile(filepath.Join(TasksDir(root), "garbage.md"), []byte("not frontmatter\njust text\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	summaries, err := LoadAllSummaries(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 1 || summaries[0].ID != "task-1" {
		t.Fatalf("summaries = %+v, want only task-1", summaries)
	}
}

func TestLoadAllSummariesEmptyDir(t *testing.T) {
	root := newTempRoot(t)
	summaries, err := LoadAllSummaries(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 0 {
		t.Fatalf("len = %d, want 0", len(summaries))
	}
}

func TestLoadAllSummariesMissingDir(t *testing.T) {
	root := t.TempDir()
	summaries, err := LoadAllSummaries(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 0 {
		t.Fatalf("len = %d, want 0", len(summaries))
	}
}

// ---------------------------------------------------------------------------
// findById
// ---------------------------------------------------------------------------

func TestFindByID(t *testing.T) {
	root := newTempRoot(t)
	mustCreate(t, root, "task-7-find-me.md", makeFrontmatter(map[string]interface{}{"id": "task-7", "title": "Find Me"}))
	summary, err := FindByID(root, "task-7")
	if err != nil {
		t.Fatal(err)
	}
	if summary.ID != "task-7" || summary.Title != "Find Me" {
		t.Fatalf("summary = %+v", summary)
	}
}

func TestFindByIDMissing(t *testing.T) {
	root := newTempRoot(t)
	if _, err := FindByID(root, "task-999"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestFindByIDDoesNotMatchLongerID(t *testing.T) {
	root := newTempRoot(t)
	mustCreate(t, root, "task-12-foo.md", makeFrontmatter(map[string]interface{}{"id": "task-12", "title": "Foo"}))
	if _, err := FindByID(root, "task-1"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

// ---------------------------------------------------------------------------
// writeTask
// ---------------------------------------------------------------------------

func TestWriteTaskBumpsUpdated(t *testing.T) {
	root := newTempRoot(t)
	fm := makeFrontmatter(map[string]interface{}{"id": "task-1", "updated": "2020-01-01"})
	created, err := CreateTask(root, fm, "task-1-test-task.md", "")
	if err != nil {
		t.Fatal(err)
	}
	created.Updated = "2020-01-01"
	if err := WriteTask(created); err != nil {
		t.Fatal(err)
	}
	reloaded, err := LoadTask(created.FilePath)
	if err != nil {
		t.Fatal(err)
	}
	want := time.Now().UTC().Format("2006-01-02")
	if reloaded.Updated != want {
		t.Fatalf("updated = %q, want %q", reloaded.Updated, want)
	}
}

func TestWriteTaskPersistsChanges(t *testing.T) {
	root := newTempRoot(t)
	created, err := CreateTask(root, makeFrontmatter(nil), "task-1-test-task.md", "")
	if err != nil {
		t.Fatal(err)
	}
	created.Status = types.StatusInProgress
	created.Priority = types.PriorityHigh
	if err := WriteTask(created); err != nil {
		t.Fatal(err)
	}
	reloaded, err := LoadTask(created.FilePath)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.Status != types.StatusInProgress || reloaded.Priority != types.PriorityHigh {
		t.Fatalf("reloaded = %+v", reloaded)
	}
}

func TestWriteTaskPreservesBodyAndOptionalFields(t *testing.T) {
	root := newTempRoot(t)
	fm := makeFrontmatter(map[string]interface{}{
		"id":     "task-9",
		"parent": "epic-1",
		"tags":   []string{"frontend", "auth"},
	})
	created, err := CreateTask(root, fm, "task-9-nested.md", "")
	if err != nil {
		t.Fatal(err)
	}
	created.Body += "custom note\n"
	if err := WriteTask(created); err != nil {
		t.Fatal(err)
	}
	reloaded, err := LoadTask(created.FilePath)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.Parent != "epic-1" || len(reloaded.Tags) != 2 || reloaded.Tags[0] != "frontend" {
		t.Fatalf("reloaded = %+v", reloaded)
	}
	if !strings.Contains(reloaded.Body, "custom note") {
		t.Fatalf("body lost content: %q", reloaded.Body)
	}
}

// ---------------------------------------------------------------------------
// frontmatter validation (Issue 6)
// ---------------------------------------------------------------------------

func TestRejectsInvalidFrontmatterFiles(t *testing.T) {
	cases := []struct {
		name  string
		front string
	}{
		{"missing id", "title: X\ntype: task\nstatus: todo\npriority: medium\ncreated: 2024-01-01\nupdated: 2024-01-01\n"},
		{"bad type", "id: task-1\ntitle: X\ntype: banana\nstatus: todo\npriority: medium\ncreated: 2024-01-01\nupdated: 2024-01-01\n"},
		{"bad status", "id: task-1\ntitle: X\ntype: task\nstatus: whenever\npriority: medium\ncreated: 2024-01-01\nupdated: 2024-01-01\n"},
		{"bad priority", "id: task-1\ntitle: X\ntype: task\nstatus: todo\npriority: urgent\ncreated: 2024-01-01\nupdated: 2024-01-01\n"},
		{"numeric id", "id: 42\ntitle: X\ntype: task\nstatus: todo\npriority: medium\ncreated: 2024-01-01\nupdated: 2024-01-01\n"},
		{"non-string parent", "id: task-1\ntitle: X\ntype: task\nstatus: todo\npriority: medium\nparent: 42\ncreated: 2024-01-01\nupdated: 2024-01-01\n"},
		{"non-array tags", "id: task-1\ntitle: X\ntype: task\nstatus: todo\npriority: medium\ntags: not-an-array\ncreated: 2024-01-01\nupdated: 2024-01-01\n"},
		{"non-string tag element", "id: task-1\ntitle: X\ntype: task\nstatus: todo\npriority: medium\ntags: [1, 2]\ncreated: 2024-01-01\nupdated: 2024-01-01\n"},
		{"non-array related", "id: task-1\ntitle: X\ntype: task\nstatus: todo\npriority: medium\nrelated: foo\ncreated: 2024-01-01\nupdated: 2024-01-01\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			raw := "---\n" + tc.front + "---\n# Body\n"
			_, err := parseFile(raw, "/tmp/x.md")
			if err == nil {
				t.Fatalf("expected invalid frontmatter to be rejected:\n%s", tc.front)
			}
		})
	}
}

func TestAcceptsValidFrontmatterWithOptionalFields(t *testing.T) {
	raw := "---\n" +
		"id: task-1\ntitle: X\ntype: subtask\nstatus: in-progress\npriority: high\n" +
		"parent: task-2\nrelated: [a, b]\ndocs: [d1]\ntags: [t1]\n" +
		"created: \"2024-01-01\"\nupdated: \"2024-01-01\"\n" +
		"---\n# Body\n"
	summary, err := parseFile(raw, "/tmp/x.md")
	if err != nil {
		t.Fatal(err)
	}
	if summary.Parent != "task-2" || len(summary.Related) != 2 || len(summary.Docs) != 1 || len(summary.Tags) != 1 {
		t.Fatalf("summary = %+v", summary)
	}
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func mustCreate(t *testing.T, root string, filename string, fm types.TaskFrontmatter) types.Task {
	t.Helper()
	task, err := CreateTask(root, fm, filename, "")
	if err != nil {
		t.Fatal(err)
	}
	return task
}
