package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
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

func newTempBoard(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	board := filepath.Join(dir, "docs", "board.md")
	writeInitialBoard(t, board)
	return board
}

// writeInitialBoard scaffolds a fresh board at the given path.
func writeInitialBoard(t *testing.T, board string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(board), 0o755); err != nil {
		t.Fatal(err)
	}
	raw, err := InitialBoard(types.LokanConfig{Counter: 0, Version: "1"})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(board, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustCreate(t *testing.T, board string, fm types.TaskFrontmatter) types.Task {
	t.Helper()
	task, err := CreateTask(board, fm, "")
	if err != nil {
		t.Fatal(err)
	}
	return task
}

// ---------------------------------------------------------------------------
// createTask
// ---------------------------------------------------------------------------

func TestCreateTaskWritesBoard(t *testing.T) {
	board := newTempBoard(t)
	if _, err := CreateTask(board, makeFrontmatter(nil), ""); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(board); err != nil {
		t.Fatalf("board file missing: %v", err)
	}
	summaries, err := LoadAllSummaries(board)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 1 {
		t.Fatalf("len = %d, want 1", len(summaries))
	}
}

func TestCreateTaskReturnsFields(t *testing.T) {
	board := newTempBoard(t)
	fm := makeFrontmatter(map[string]interface{}{"id": "5", "title": "Big Epic", "type": types.TypeEpic})
	task, err := CreateTask(board, fm, "")
	if err != nil {
		t.Fatal(err)
	}
	if task.ID != "5" || task.Title != "Big Epic" {
		t.Fatalf("task = %+v", task)
	}
	if task.FilePath != VirtualPath(board, "5") {
		t.Fatalf("filePath = %q", task.FilePath)
	}
}

// ---------------------------------------------------------------------------
// loadTask round-trip
// ---------------------------------------------------------------------------

func TestLoadTaskRoundTrip(t *testing.T) {
	board := newTempBoard(t)
	fm := makeFrontmatter(map[string]interface{}{"id": "42", "title": "Round-trip Task", "priority": types.PriorityHigh})
	created, err := CreateTask(board, fm, "")
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
	board := newTempBoard(t)
	if _, err := LoadTask(VirtualPath(board, "999")); !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

// ---------------------------------------------------------------------------
// loadAllSummaries
// ---------------------------------------------------------------------------

func TestLoadAllSummariesReturnsAll(t *testing.T) {
	board := newTempBoard(t)
	mustCreate(t, board, makeFrontmatter(map[string]interface{}{"id": "1", "title": "First"}))
	mustCreate(t, board, makeFrontmatter(map[string]interface{}{"id": "2", "title": "Second"}))

	summaries, err := LoadAllSummaries(board)
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
	board := newTempBoard(t)
	raw := "# Kanlo Board\n\n## Active\n\n" +
		"<!-- lokan:1 -->\n" +
		"---\nid: \"1\"\ntitle: Valid\ntype: task\nstatus: todo\npriority: medium\ncreated: \"2024-01-01\"\nupdated: \"2024-01-01\"\n---\n# Valid\n\n" +
		"<!-- lokan:bad -->\n" +
		"this is not frontmatter\n"
	if err := os.MkdirAll(filepath.Dir(board), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(board, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}

	summaries, err := LoadAllSummaries(board)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 1 || summaries[0].ID != "1" {
		t.Fatalf("summaries = %+v, want only 1", summaries)
	}
}

func TestLoadAllSummariesEmptyBoard(t *testing.T) {
	board := newTempBoard(t)
	if err := os.MkdirAll(filepath.Dir(board), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(board, []byte("# Kanlo Board\n\n## Active\n\n## Archive\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	summaries, err := LoadAllSummaries(board)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 0 {
		t.Fatalf("len = %d, want 0", len(summaries))
	}
}

func TestLoadAllSummariesMissingBoard(t *testing.T) {
	board := filepath.Join(t.TempDir(), "board.md")
	if _, err := LoadAllSummaries(board); err == nil {
		t.Fatalf("missing board should error")
	}
}

// ---------------------------------------------------------------------------
// findById
// ---------------------------------------------------------------------------

func TestFindByID(t *testing.T) {
	board := newTempBoard(t)
	mustCreate(t, board, makeFrontmatter(map[string]interface{}{"id": "7", "title": "Find Me"}))
	summary, err := FindByID(board, "7")
	if err != nil {
		t.Fatal(err)
	}
	if summary.ID != "7" || summary.Title != "Find Me" {
		t.Fatalf("summary = %+v", summary)
	}
}

func TestFindByIDMissing(t *testing.T) {
	board := newTempBoard(t)
	if _, err := FindByID(board, "999"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestFindByIDDoesNotMatchLongerID(t *testing.T) {
	board := newTempBoard(t)
	mustCreate(t, board, makeFrontmatter(map[string]interface{}{"id": "12", "title": "Foo"}))
	if _, err := FindByID(board, "1"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

// ---------------------------------------------------------------------------
// writeTask
// ---------------------------------------------------------------------------

func TestWriteTaskBumpsUpdated(t *testing.T) {
	board := newTempBoard(t)
	fm := makeFrontmatter(map[string]interface{}{"id": "1", "updated": "2020-01-01"})
	created, err := CreateTask(board, fm, "")
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
	board := newTempBoard(t)
	created, err := CreateTask(board, makeFrontmatter(nil), "")
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
	board := newTempBoard(t)
	fm := makeFrontmatter(map[string]interface{}{
		"id":     "9",
		"parent": "1",
		"tags":   []string{"frontend", "auth"},
	})
	created, err := CreateTask(board, fm, "")
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
	board := newTempBoard(t)
	created, err := CreateTask(board, makeFrontmatter(nil), "")
	if err != nil {
		t.Fatal(err)
	}
	created.Status = types.StatusDone
	if err := WriteTask(created); err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile(board)
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

	summaries, err := LoadAllSummaries(board)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 1 || summaries[0].ID != "1" {
		t.Fatalf("summaries = %+v, want 1", summaries)
	}
}

func TestWriteTaskUnarchivesOnReopen(t *testing.T) {
	board := newTempBoard(t)
	created, err := CreateTask(board, makeFrontmatter(nil), "")
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

	raw, err := os.ReadFile(board)
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
	board := newTempBoard(t)
	const n = 30
	var wg sync.WaitGroup
	errs := make(chan error, n)
	for i := 1; i <= n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fm := makeFrontmatter(map[string]interface{}{"id": fmt.Sprintf("%d", i), "title": fmt.Sprintf("T%d", i)})
			if _, err := CreateTask(board, fm, ""); err != nil {
				errs <- err
			}
		}(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatalf("CreateTask: %v", err)
	}

	summaries, err := LoadAllSummaries(board)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != n {
		t.Fatalf("got %d tasks, want %d (lost updates)", len(summaries), n)
	}
}

// ---------------------------------------------------------------------------
// moveTask
// ---------------------------------------------------------------------------

func boardIDs(t *testing.T, board string) []string {
	t.Helper()
	summaries, err := LoadAllSummaries(board)
	if err != nil {
		t.Fatal(err)
	}
	ids := make([]string, len(summaries))
	for i, s := range summaries {
		ids[i] = s.ID
	}
	return ids
}

func TestMoveTaskReordersWithinLane(t *testing.T) {
	board := newTempBoard(t)
	for _, id := range []string{"1", "2", "3"} {
		mustCreate(t, board, makeFrontmatter(map[string]interface{}{"id": id}))
	}
	// move 1 before 3 → 2,1,3
	if _, err := MoveTask(board, "1", types.StatusTodo, "3"); err != nil {
		t.Fatal(err)
	}
	if got, want := boardIDs(t, board), []string{"2", "1", "3"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("order = %v, want %v", got, want)
	}
}

func TestMoveTaskToEndOfLane(t *testing.T) {
	board := newTempBoard(t)
	for _, id := range []string{"1", "2", "3"} {
		mustCreate(t, board, makeFrontmatter(map[string]interface{}{"id": id}))
	}
	// empty anchor = append at the lane's end → 2,3,1
	if _, err := MoveTask(board, "1", types.StatusTodo, ""); err != nil {
		t.Fatal(err)
	}
	if got, want := boardIDs(t, board), []string{"2", "3", "1"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("order = %v, want %v", got, want)
	}
}

func TestMoveTaskAcrossLanesArchives(t *testing.T) {
	board := newTempBoard(t)
	for _, id := range []string{"1", "2"} {
		mustCreate(t, board, makeFrontmatter(map[string]interface{}{"id": id}))
	}
	// move 1 into done → archived at the end, status applied
	moved, err := MoveTask(board, "1", types.StatusDone, "")
	if err != nil {
		t.Fatal(err)
	}
	if moved.Status != types.StatusDone {
		t.Fatalf("status = %s, want done", moved.Status)
	}
	if got, want := boardIDs(t, board), []string{"2", "1"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("order = %v, want %v", got, want)
	}
}

func TestMoveTaskToEmptyActiveLane(t *testing.T) {
	board := newTempBoard(t)
	mustCreate(t, board, makeFrontmatter(map[string]interface{}{"id": "1"}))
	mustCreate(t, board, makeFrontmatter(map[string]interface{}{"id": "2", "status": types.StatusDone}))
	// backlog is empty: lands before the archive section, not after it
	if _, err := MoveTask(board, "1", types.StatusBacklog, ""); err != nil {
		t.Fatal(err)
	}
	if got, want := boardIDs(t, board), []string{"1", "2"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("order = %v, want %v", got, want)
	}
}

func TestMoveTaskNotFound(t *testing.T) {
	board := newTempBoard(t)
	mustCreate(t, board, makeFrontmatter(nil))
	if _, err := MoveTask(board, "99", types.StatusTodo, ""); !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
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
			_, err := parseFile(raw, "/tmp/x.md", types.DefaultStatusDefs())
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
	summary, err := parseFile(raw, "/tmp/x.md", types.DefaultStatusDefs())
	if err != nil {
		t.Fatal(err)
	}
	if summary.Parent != "2" || len(summary.Related) != 2 || len(summary.Docs) != 1 || len(summary.Tags) != 1 {
		t.Fatalf("summary = %+v", summary)
	}
}

// ---------------------------------------------------------------------------
// configurable lanes
// ---------------------------------------------------------------------------

// writeLanes persists a custom lane set into the board's config block.
func writeLanes(t *testing.T, board string, lanes []types.StatusDef) {
	t.Helper()
	cfg, _ := ReadConfig(board)
	cfg.Statuses = lanes
	if err := WriteConfig(board, cfg); err != nil {
		t.Fatal(err)
	}
}

func TestParseAcceptsConfiguredStatus(t *testing.T) {
	// a custom lane parses only when the parser is given its lane set
	custom := []types.StatusDef{{ID: "doing"}, {ID: "done", Archived: true}}
	raw := "---\nid: \"1\"\ntitle: X\ntype: task\nstatus: doing\npriority: medium\ncreated: \"2024-01-01\"\nupdated: \"2024-01-01\"\n---\n# Body\n"
	if _, err := parseFile(raw, "/tmp/x.md", custom); err != nil {
		t.Fatalf("custom status should parse with matching lanes: %v", err)
	}
	if _, err := parseFile(raw, "/tmp/x.md", types.DefaultStatusDefs()); err == nil {
		t.Fatalf("custom status should be rejected by default lanes")
	}
}

func TestBoardRoundTripsConfiguredStatuses(t *testing.T) {
	// a board holding a custom status loads fully only when the project
	// config knows the lane — never silently dropped
	board := newTempBoard(t)
	writeLanes(t, board, []types.StatusDef{{ID: "doing"}, {ID: "done", Archived: true}})
	fm := makeFrontmatter(nil)
	fm.Status = "doing"
	mustCreate(t, board, fm)

	summaries, err := LoadAllSummaries(board)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 1 || summaries[0].Status != "doing" {
		t.Fatalf("summaries = %+v", summaries)
	}
}

func TestMoveLaneRewritesTaskStatuses(t *testing.T) {
	board := newTempBoard(t)
	// the target lane must be configured or the rewritten board won't parse
	writeLanes(t, board, []types.StatusDef{{ID: "todo"}, {ID: "doing"}, {ID: "done", Archived: true}})
	mustCreate(t, board, makeFrontmatter(nil))
	mustCreate(t, board, makeFrontmatter(map[string]interface{}{"id": "2", "status": types.StatusDone}))

	moved, err := MoveLane(board, types.StatusTodo, "doing")
	if err != nil {
		t.Fatal(err)
	}
	if moved != 1 {
		t.Fatalf("moved = %d, want 1", moved)
	}
	task, err := FindByID(board, "1")
	if err != nil {
		t.Fatal(err)
	}
	if task.Status != "doing" {
		t.Fatalf("status = %q, want doing", task.Status)
	}
	// untouched lanes keep their status
	task2, err := FindByID(board, "2")
	if err != nil {
		t.Fatal(err)
	}
	if task2.Status != types.StatusDone {
		t.Fatalf("status = %q, want done", task2.Status)
	}
}

func TestClearArchivedDeletesOnlyArchivedLanes(t *testing.T) {
	board := newTempBoard(t)
	mustCreate(t, board, makeFrontmatter(nil))
	mustCreate(t, board, makeFrontmatter(map[string]interface{}{"id": "2", "status": types.StatusDone}))
	mustCreate(t, board, makeFrontmatter(map[string]interface{}{"id": "3", "status": types.StatusCancelled}))

	deleted, err := ClearArchived(board)
	if err != nil {
		t.Fatal(err)
	}
	if deleted != 2 {
		t.Fatalf("deleted = %d, want 2", deleted)
	}
	remaining, err := LoadAllSummaries(board)
	if err != nil {
		t.Fatal(err)
	}
	if len(remaining) != 1 || remaining[0].ID != "1" {
		t.Fatalf("remaining = %+v", remaining)
	}
}

func TestClearAllDeletesEveryTask(t *testing.T) {
	board := newTempBoard(t)
	mustCreate(t, board, makeFrontmatter(nil))
	mustCreate(t, board, makeFrontmatter(map[string]interface{}{"id": "2", "status": types.StatusDone}))

	deleted, err := ClearAll(board)
	if err != nil {
		t.Fatal(err)
	}
	if deleted != 2 {
		t.Fatalf("deleted = %d, want 2", deleted)
	}
	remaining, err := LoadAllSummaries(board)
	if err != nil {
		t.Fatal(err)
	}
	if len(remaining) != 0 {
		t.Fatalf("remaining = %+v, want empty", remaining)
	}
}

func TestClearArchivedUsesConfiguredArchivedFlag(t *testing.T) {
	// a custom lane marked archived counts for clear-archived
	board := newTempBoard(t)
	writeLanes(t, board, []types.StatusDef{{ID: "doing"}, {ID: "shipped", Archived: true}})
	mustCreate(t, board, makeFrontmatter(map[string]interface{}{"status": types.Status("shipped")}))

	deleted, err := ClearArchived(board)
	if err != nil {
		t.Fatal(err)
	}
	if deleted != 1 {
		t.Fatalf("deleted = %d, want 1", deleted)
	}
}

// ---------------------------------------------------------------------------
// self-contained config & explicit board paths

func TestOpsTargetExplicitBoardPath(t *testing.T) {
	dir := t.TempDir()
	// two boards in the same project — targeting is explicit per call
	boardA := filepath.Join(dir, "docs", "board.md")
	boardB := filepath.Join(dir, "notes", "roadmap.md")
	writeInitialBoard(t, boardA)
	writeInitialBoard(t, boardB)

	task := mustCreate(t, boardB, makeFrontmatter(nil))
	summaries, err := LoadAllSummaries(boardB)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 1 || summaries[0].ID != task.ID {
		t.Fatalf("summaries = %+v", summaries)
	}

	// the other board stays untouched
	aSummaries, err := LoadAllSummaries(boardA)
	if err != nil {
		t.Fatal(err)
	}
	if len(aSummaries) != 0 {
		t.Fatalf("board A should stay empty, got %+v", aSummaries)
	}
}

func TestConfigLivesInBoard(t *testing.T) {
	board := newTempBoard(t)
	if err := WriteConfig(board, types.LokanConfig{
		Counter:  41,
		Version:  "2",
		Statuses: []types.StatusDef{{ID: "doing"}, {ID: "done", Archived: true}},
	}); err != nil {
		t.Fatal(err)
	}
	cfg, err := ReadConfig(board)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Counter != 41 || cfg.Version != "2" || len(cfg.Statuses) != 2 {
		t.Fatalf("cfg = %+v", cfg)
	}

	// the config block survives task writes
	mustCreate(t, board, makeFrontmatter(nil))
	cfg, err = ReadConfig(board)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Counter != 41 || len(cfg.Statuses) != 2 {
		t.Fatalf("config lost after task write: %+v", cfg)
	}
}

func TestNextCounterIncrementsInBoard(t *testing.T) {
	board := newTempBoard(t)
	for want := 1; want <= 3; want++ {
		if got, err := NextCounter(board); err != nil || got != want {
			t.Fatalf("NextCounter = %d, %v, want %d", got, err, want)
		}
	}
	cfg, err := ReadConfig(board)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Counter != 3 {
		t.Fatalf("counter = %d, want 3", cfg.Counter)
	}
}

func TestNextCounterConcurrentUnique(t *testing.T) {
	board := newTempBoard(t)
	const n = 50

	var wg sync.WaitGroup
	counts := make(chan int, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c, err := NextCounter(board)
			if err != nil {
				t.Errorf("NextCounter: %v", err)
				return
			}
			counts <- c
		}()
	}
	wg.Wait()
	close(counts)

	seen := make(map[int]bool, n)
	for c := range counts {
		if seen[c] {
			t.Fatalf("duplicate counter handed out: %d", c)
		}
		if c < 1 || c > n {
			t.Fatalf("counter out of range: %d", c)
		}
		seen[c] = true
	}
	if len(seen) != n {
		t.Fatalf("got %d unique counters, want %d", len(seen), n)
	}
}
