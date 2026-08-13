package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/avressatelier/lokan/internal/types"
)

func makeFrontmatter(overrides map[string]interface{}) types.TaskFrontmatter {
	fm := types.TaskFrontmatter{
		ID:       "1",
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
	return t.TempDir()
}

func mustCreate(t *testing.T, root string, fm types.TaskFrontmatter) types.Task {
	t.Helper()
	task, err := CreateTask(root, fm, "")
	if err != nil {
		t.Fatal(err)
	}
	return task
}

// ---------------------------------------------------------------------------
// createTask
// ---------------------------------------------------------------------------

func TestCreateTaskWritesBoard(t *testing.T) {
	root := newTempRoot(t)
	if _, err := CreateTask(root, makeFrontmatter(nil), ""); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(BoardPath(root)); err != nil {
		t.Fatalf("board file missing: %v", err)
	}
	if _, err := os.Stat(TasksDir(root)); !os.IsNotExist(err) {
		t.Fatalf("tasks dir should not exist: %v", err)
	}
	summaries, err := LoadAllSummaries(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 1 {
		t.Fatalf("len = %d, want 1", len(summaries))
	}
}

func TestCreateTaskReturnsFields(t *testing.T) {
	root := newTempRoot(t)
	fm := makeFrontmatter(map[string]interface{}{"id": "5", "title": "Big Epic", "type": types.TypeEpic})
	task, err := CreateTask(root, fm, "")
	if err != nil {
		t.Fatal(err)
	}
	if task.ID != "5" || task.Title != "Big Epic" {
		t.Fatalf("task = %+v", task)
	}
	if task.FilePath != VirtualPath(root, "5") {
		t.Fatalf("filePath = %q", task.FilePath)
	}
}

// ---------------------------------------------------------------------------
// loadTask round-trip
// ---------------------------------------------------------------------------

func TestLoadTaskRoundTrip(t *testing.T) {
	root := newTempRoot(t)
	fm := makeFrontmatter(map[string]interface{}{"id": "42", "title": "Round-trip Task", "priority": types.PriorityHigh})
	created, err := CreateTask(root, fm, "")
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

func TestLoadTaskMissing(t *testing.T) {
	root := newTempRoot(t)
	if _, err := LoadTask(VirtualPath(root, "999")); !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

// ---------------------------------------------------------------------------
// loadAllSummaries
// ---------------------------------------------------------------------------

func TestLoadAllSummariesReturnsAll(t *testing.T) {
	root := newTempRoot(t)
	mustCreate(t, root, makeFrontmatter(map[string]interface{}{"id": "1", "title": "First"}))
	mustCreate(t, root, makeFrontmatter(map[string]interface{}{"id": "2", "title": "Second"}))

	summaries, err := LoadAllSummaries(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 2 {
		t.Fatalf("len = %d, want 2", len(summaries))
	}
	ids := []string{summaries[0].ID, summaries[1].ID}
	if (ids[0] != "1" || ids[1] != "2") && (ids[0] != "2" || ids[1] != "1") {
		t.Fatalf("ids = %v", ids)
	}
}

func TestLoadAllSummariesSkipsInvalid(t *testing.T) {
	root := newTempRoot(t)
	raw := "# Kanlo Board\n\n## Active\n\n" +
		"<!-- lokan:1 -->\n" +
		"---\nid: \"1\"\ntitle: Valid\ntype: task\nstatus: todo\npriority: medium\ncreated: \"2024-01-01\"\nupdated: \"2024-01-01\"\n---\n# Valid\n\n" +
		"<!-- lokan:bad -->\n" +
		"this is not frontmatter\n"
	if err := os.MkdirAll(filepath.Dir(BoardPath(root)), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(BoardPath(root), []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}

	summaries, err := LoadAllSummaries(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 1 || summaries[0].ID != "1" {
		t.Fatalf("summaries = %+v, want only 1", summaries)
	}
}

func TestLoadAllSummariesEmptyBoard(t *testing.T) {
	root := newTempRoot(t)
	if err := os.MkdirAll(filepath.Dir(BoardPath(root)), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(BoardPath(root), []byte("# Kanlo Board\n\n## Active\n\n## Archive\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	summaries, err := LoadAllSummaries(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 0 {
		t.Fatalf("len = %d, want 0", len(summaries))
	}
}

func TestLoadAllSummariesMissingBoard(t *testing.T) {
	root := newTempRoot(t)
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
	mustCreate(t, root, makeFrontmatter(map[string]interface{}{"id": "7", "title": "Find Me"}))
	summary, err := FindByID(root, "7")
	if err != nil {
		t.Fatal(err)
	}
	if summary.ID != "7" || summary.Title != "Find Me" {
		t.Fatalf("summary = %+v", summary)
	}
}

func TestFindByIDMissing(t *testing.T) {
	root := newTempRoot(t)
	if _, err := FindByID(root, "999"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestFindByIDDoesNotMatchLongerID(t *testing.T) {
	root := newTempRoot(t)
	mustCreate(t, root, makeFrontmatter(map[string]interface{}{"id": "12", "title": "Foo"}))
	if _, err := FindByID(root, "1"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

// ---------------------------------------------------------------------------
// writeTask
// ---------------------------------------------------------------------------

func TestWriteTaskBumpsUpdated(t *testing.T) {
	root := newTempRoot(t)
	fm := makeFrontmatter(map[string]interface{}{"id": "1", "updated": "2020-01-01"})
	created, err := CreateTask(root, fm, "")
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
	created, err := CreateTask(root, makeFrontmatter(nil), "")
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
		"id":     "9",
		"parent": "1",
		"tags":   []string{"frontend", "auth"},
	})
	created, err := CreateTask(root, fm, "")
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
	if reloaded.Parent != "1" || len(reloaded.Tags) != 2 || reloaded.Tags[0] != "frontend" {
		t.Fatalf("reloaded = %+v", reloaded)
	}
	if !strings.Contains(reloaded.Body, "custom note") {
		t.Fatalf("body lost content: %q", reloaded.Body)
	}
}

// ---------------------------------------------------------------------------
// archive grouping
// ---------------------------------------------------------------------------

func TestArchiveSectionGroupsDoneTasks(t *testing.T) {
	root := newTempRoot(t)
	created, err := CreateTask(root, makeFrontmatter(nil), "")
	if err != nil {
		t.Fatal(err)
	}
	created.Status = types.StatusDone
	if err := WriteTask(created); err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile(BoardPath(root))
	if err != nil {
		t.Fatal(err)
	}
	content := string(raw)
	if !strings.Contains(content, "## Archive\n\n<!-- lokan:1 -->") {
		t.Fatalf("done task not under Archive section:\n%s", content)
	}
	if strings.Contains(content, "## Active\n\n<!-- lokan:1 -->") {
		t.Fatalf("done task still under Active:\n%s", content)
	}

	summaries, err := LoadAllSummaries(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 1 || summaries[0].ID != "1" {
		t.Fatalf("summaries = %+v, want 1", summaries)
	}
}

func TestWriteTaskUnarchivesOnReopen(t *testing.T) {
	root := newTempRoot(t)
	created, err := CreateTask(root, makeFrontmatter(nil), "")
	if err != nil {
		t.Fatal(err)
	}
	created.Status = types.StatusDone
	if err := WriteTask(created); err != nil {
		t.Fatal(err)
	}
	created.Status = types.StatusTodo
	if err := WriteTask(created); err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile(BoardPath(root))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "## Active\n\n<!-- lokan:1 -->") {
		t.Fatalf("reopened task not back under Active:\n%s", raw)
	}
}

// ---------------------------------------------------------------------------
// concurrency
// ---------------------------------------------------------------------------

func TestConcurrentCreateNoLostUpdates(t *testing.T) {
	root := newTempRoot(t)
	const n = 30
	var wg sync.WaitGroup
	errs := make(chan error, n)
	for i := 1; i <= n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fm := makeFrontmatter(map[string]interface{}{"id": fmt.Sprintf("%d", i), "title": fmt.Sprintf("T%d", i)})
			if _, err := CreateTask(root, fm, ""); err != nil {
				errs <- err
			}
		}(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatalf("CreateTask: %v", err)
	}

	summaries, err := LoadAllSummaries(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != n {
		t.Fatalf("got %d tasks, want %d (lost updates)", len(summaries), n)
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
		{"bad type", "id: \"1\"\ntitle: X\ntype: banana\nstatus: todo\npriority: medium\ncreated: 2024-01-01\nupdated: 2024-01-01\n"},
		{"bad status", "id: \"1\"\ntitle: X\ntype: task\nstatus: whenever\npriority: medium\ncreated: 2024-01-01\nupdated: 2024-01-01\n"},
		{"bad priority", "id: \"1\"\ntitle: X\ntype: task\nstatus: todo\npriority: urgent\ncreated: 2024-01-01\nupdated: 2024-01-01\n"},
		{"numeric id", "id: 42\ntitle: X\ntype: task\nstatus: todo\npriority: medium\ncreated: 2024-01-01\nupdated: 2024-01-01\n"},
		{"non-string parent", "id: \"1\"\ntitle: X\ntype: task\nstatus: todo\npriority: medium\nparent: 42\ncreated: 2024-01-01\nupdated: 2024-01-01\n"},
		{"non-array tags", "id: \"1\"\ntitle: X\ntype: task\nstatus: todo\npriority: medium\ntags: not-an-array\ncreated: 2024-01-01\nupdated: 2024-01-01\n"},
		{"non-string tag element", "id: \"1\"\ntitle: X\ntype: task\nstatus: todo\npriority: medium\ntags: [1, 2]\ncreated: 2024-01-01\nupdated: 2024-01-01\n"},
		{"non-array related", "id: \"1\"\ntitle: X\ntype: task\nstatus: todo\npriority: medium\nrelated: foo\ncreated: 2024-01-01\nupdated: 2024-01-01\n"},
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
		"id: \"1\"\ntitle: X\ntype: subtask\nstatus: in-progress\npriority: high\n" +
		"parent: \"2\"\nrelated: [a, b]\ndocs: [d1]\ntags: [t1]\n" +
		"created: \"2024-01-01\"\nupdated: \"2024-01-01\"\n" +
		"---\n# Body\n"
	summary, err := parseFile(raw, "/tmp/x.md")
	if err != nil {
		t.Fatal(err)
	}
	if summary.Parent != "2" || len(summary.Related) != 2 || len(summary.Docs) != 1 || len(summary.Tags) != 1 {
		t.Fatalf("summary = %+v", summary)
	}
}
