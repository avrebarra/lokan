package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/avrebarra/lokan/internal/store"
	"github.com/avrebarra/lokan/internal/types"
)

// runCLI executes the CLI with args from within dir, returning exit code,
// stdout, and stderr.
func runCLI(t *testing.T, dir string, args ...string) (int, string, string) {
	t.Helper()
	t.Chdir(dir)
	var stdout, stderr bytes.Buffer
	code := run(args, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

func initProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if code, _, stderr := runCLI(t, root, "init", bp(root)); code != 0 {
		t.Fatalf("init failed: %s", stderr)
	}
	return root
}

// bp returns the board path for a test project dir.
func bp(dir string) string {
	return filepath.Join(dir, "docs", "board.md")
}

// runBoard runs the CLI from dir, injecting the project's board path as the
// first positional argument.
func runBoard(t *testing.T, dir string, args ...string) (int, string, string) {
	t.Helper()
	full := append([]string{args[0], bp(dir)}, args[1:]...)
	return runCLI(t, dir, full...)
}

func mustCreate(t *testing.T, dir string, args ...string) string {
	t.Helper()
	args = append([]string{"create", bp(dir)}, args...)
	if code, _, stderr := runCLI(t, dir, args...); code != 0 {
		t.Fatalf("create %v failed: %s", args, stderr)
	}
	return dir
}

func taskIDs(t *testing.T, dir string) []string {
	t.Helper()
	all, err := store.LoadAllSummaries(bp(dir))
	if err != nil {
		t.Fatal(err)
	}
	ids := make([]string, len(all))
	for i, s := range all {
		ids[i] = s.ID
	}
	return ids
}

// ---------------------------------------------------------------------------
// init
// ---------------------------------------------------------------------------

func TestInitCreatesConfig(t *testing.T) {
	root := t.TempDir()
	code, stdout, stderr := runBoard(t, root, "init")
	if code != 0 {
		t.Fatalf("init failed: %s", stderr)
	}
	if !strings.Contains(stdout, "Initialized lokan project.") {
		t.Fatalf("stdout = %q", stdout)
	}
	cfg, err := store.ReadConfig(bp(root))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Counter != 0 || cfg.Version != "1" {
		t.Fatalf("cfg = %+v", cfg)
	}
	if _, err := os.Stat(bp(root)); err != nil {
		t.Fatalf("board file missing: %v", err)
	}
}

func TestInitIdempotent(t *testing.T) {
	root := t.TempDir()
	if code, _, stderr := runBoard(t, root, "init"); code != 0 {
		t.Fatalf("first init failed: %s", stderr)
	}
	code, stdout, _ := runBoard(t, root, "init")
	if code != 0 {
		t.Fatalf("second init should exit 0")
	}
	if !strings.Contains(stdout, "Already a lokan board") {
		t.Fatalf("stdout = %q", stdout)
	}
}

func TestInitWithCustomBoard(t *testing.T) {
	root := t.TempDir()
	board := filepath.Join(root, "docs", "roadmap.md")
	code, stdout, stderr := runCLI(t, root, "init", board)
	if code != 0 {
		t.Fatalf("init failed: %s", stderr)
	}
	if !strings.Contains(stdout, "Initialized lokan project.") {
		t.Fatalf("stdout = %q", stdout)
	}
	if _, err := os.Stat(board); err != nil {
		t.Fatalf("custom board missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "docs", "board.md")); !os.IsNotExist(err) {
		t.Fatalf("default board should not exist, err = %v", err)
	}
	// the config lives inside the custom board
	cfg, err := store.ReadConfig(board)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Counter != 0 || cfg.Version != "1" {
		t.Fatalf("cfg = %+v", cfg)
	}
}

// ---------------------------------------------------------------------------
// explicit board rule
// ---------------------------------------------------------------------------

func TestMissingBoardPath(t *testing.T) {
	root := t.TempDir()
	code, _, stderr := runCLI(t, root, "list")
	if code == 0 {
		t.Fatalf("list should fail without a board path")
	}
	if !strings.Contains(stderr, "accepts 1 arg(s), received 0") {
		t.Fatalf("stderr = %q", stderr)
	}
	if !strings.HasPrefix(stderr, "error: ") {
		t.Fatalf("stderr should be prefixed with error: %q", stderr)
	}
}

func TestNotABoardError(t *testing.T) {
	root := t.TempDir()
	for _, args := range [][]string{
		{"list"},
		{"create", "hello"},
		{"get", "1"},
		{"edit", "1", "--status", "done"},
		{"subtasks", "1"},
	} {
		code, _, stderr := runBoard(t, root, args...)
		if code == 0 {
			t.Fatalf("%v should fail without a board", args)
		}
		if !strings.Contains(stderr, "not a lokan board") {
			t.Fatalf("%v stderr = %q", args, stderr)
		}
	}
}

// ---------------------------------------------------------------------------
// create
// ---------------------------------------------------------------------------

func TestCreateValid(t *testing.T) {
	root := initProject(t)
	code, stdout, stderr := runBoard(t, root, "create", "hello")
	if code != 0 {
		t.Fatalf("create failed: %s", stderr)
	}
	if !strings.Contains(stdout, "Created 1") {
		t.Fatalf("stdout = %q", stdout)
	}
	summary, err := store.FindByID(bp(root), "1")
	if err != nil {
		t.Fatal(err)
	}
	if summary.Type != types.TypeTask || summary.Status != types.StatusBacklog || summary.Priority != types.PriorityMedium {
		t.Fatalf("summary = %+v", summary)
	}
	if summary.Title != "hello" {
		t.Fatalf("title = %q", summary.Title)
	}
}

func TestCreateWithTypePriorityAndTags(t *testing.T) {
	root := initProject(t)
	code, _, stderr := runBoard(t, root, "create", "--type", "bug", "--priority", "high", "--tag", "a", "--tag", "b", "broken thing")
	if code != 0 {
		t.Fatalf("create failed: %s", stderr)
	}
	summary, err := store.FindByID(bp(root), "1")
	if err != nil {
		t.Fatal(err)
	}
	if summary.Type != types.TypeBug || summary.Priority != types.PriorityHigh {
		t.Fatalf("summary = %+v", summary)
	}
	if len(summary.Tags) != 2 || summary.Tags[0] != "a" || summary.Tags[1] != "b" {
		t.Fatalf("tags = %v", summary.Tags)
	}
}

func TestCreateInvalidType(t *testing.T) {
	root := initProject(t)
	code, _, stderr := runBoard(t, root, "create", "--type", "banana", "hello")
	if code == 0 {
		t.Fatalf("invalid type should fail")
	}
	if !strings.Contains(stderr, `Invalid type "banana". Must be one of: epic, task, subtask, bug`) {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestCreateInvalidPriority(t *testing.T) {
	root := initProject(t)
	code, _, stderr := runBoard(t, root, "create", "--priority", "urgent", "hello")
	if code == 0 {
		t.Fatalf("invalid priority should fail")
	}
	if !strings.Contains(stderr, `Invalid priority "urgent". Must be one of: critical, high, medium, low`) {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestCreateParentNotFound(t *testing.T) {
	root := initProject(t)
	code, _, stderr := runBoard(t, root, "create", "--parent", "99", "hello")
	if code == 0 {
		t.Fatalf("missing parent should fail")
	}
	if !strings.Contains(stderr, "Parent task not found: 99") {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestCreateParentTypeNotAllowed(t *testing.T) {
	root := initProject(t)
	// 1 can hold task
	if code, _, stderr := runBoard(t, root, "create", "--type", "epic", "e"); code != 0 {
		t.Fatalf("create epic failed: %s", stderr)
	}
	if code, _, stderr := runBoard(t, root, "create", "--parent", "1", "t"); code != 0 {
		t.Fatalf("create task under epic failed: %s", stderr)
	}
	// task cannot be created under a task
	if code, _, stderr := runBoard(t, root, "create", "--parent", "2", "sub"); code == 0 {
		t.Fatalf("task under task should fail")
	} else if !strings.Contains(stderr, "Allowed parents: epic") {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestCreateEpicNoParent(t *testing.T) {
	root := initProject(t)
	mustCreate(t, root, "parent task")
	if code, _, stderr := runBoard(t, root, "create", "--type", "epic", "--parent", "1", "e"); code == 0 {
		t.Fatalf("epic with parent should fail")
	} else if !strings.Contains(stderr, "Allowed parents: none") {
		t.Fatalf("stderr = %q", stderr)
	}
}

// ---------------------------------------------------------------------------
// get
// ---------------------------------------------------------------------------

func TestGet(t *testing.T) {
	root := initProject(t)
	mustCreate(t, root, "hello")
	code, stdout, stderr := runBoard(t, root, "get", "1")
	if code != 0 {
		t.Fatalf("get failed: %s", stderr)
	}
	if !strings.Contains(stdout, "1 · task") {
		t.Fatalf("stdout = %q", stdout)
	}
	if !strings.Contains(stdout, "Title:    hello") {
		t.Fatalf("stdout = %q", stdout)
	}
	if !strings.Contains(stdout, "## Notes") {
		t.Fatalf("stdout = %q", stdout)
	}
}

func TestGetNotFound(t *testing.T) {
	root := initProject(t)
	code, _, stderr := runBoard(t, root, "get", "99")
	if code == 0 {
		t.Fatalf("get missing should fail")
	}
	if !strings.Contains(stderr, "task not found: 99") {
		t.Fatalf("stderr = %q", stderr)
	}
}

// ---------------------------------------------------------------------------
// edit
// ---------------------------------------------------------------------------

func TestEditStatus(t *testing.T) {
	root := initProject(t)
	mustCreate(t, root, "hello")
	code, stdout, stderr := runBoard(t, root, "edit", "1", "--status", "done")
	if code != 0 {
		t.Fatalf("edit failed: %s", stderr)
	}
	if !strings.Contains(stdout, "Updated 1") {
		t.Fatalf("stdout = %q", stdout)
	}
	summary, err := store.FindByID(bp(root), "1")
	if err != nil {
		t.Fatal(err)
	}
	if summary.Status != types.StatusDone {
		t.Fatalf("status = %q", summary.Status)
	}
}

func TestEditInvalidStatus(t *testing.T) {
	root := initProject(t)
	mustCreate(t, root, "hello")
	code, _, stderr := runBoard(t, root, "edit", "1", "--status", "whenever")
	if code == 0 {
		t.Fatalf("invalid status should fail")
	}
	if !strings.Contains(stderr, `Invalid status "whenever". Must be one of: backlog, todo, in-progress, done, cancelled`) {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestEditInvalidPriority(t *testing.T) {
	root := initProject(t)
	mustCreate(t, root, "hello")
	code, _, stderr := runBoard(t, root, "edit", "1", "--priority", "urgent")
	if code == 0 {
		t.Fatalf("invalid priority should fail")
	}
	if !strings.Contains(stderr, `Invalid priority "urgent".`) {
		t.Fatalf("stderr = %q", stderr)
	}
}

// Title edits update the task block in place — per-task file renames are gone
// with single-board storage.
func TestEditTitleUpdatesInPlace(t *testing.T) {
	root := initProject(t)
	mustCreate(t, root, "hello")
	code, _, stderr := runBoard(t, root, "edit", "1", "--title", "hello world")
	if code != 0 {
		t.Fatalf("edit title failed: %s", stderr)
	}
	summary, err := store.FindByID(bp(root), "1")
	if err != nil {
		t.Fatal(err)
	}
	if summary.Title != "hello world" {
		t.Fatalf("title = %q", summary.Title)
	}
	if summary.FilePath != store.VirtualPath(bp(root), "1") {
		t.Fatalf("filePath = %q, want %q", summary.FilePath, store.VirtualPath(bp(root), "1"))
	}
}

// Editing the title twice must keep the id stable.
func TestEditTitleTwiceKeepsID(t *testing.T) {
	root := initProject(t)
	mustCreate(t, root, "hello")
	if code, _, stderr := runBoard(t, root, "edit", "1", "--title", "second"); code != 0 {
		t.Fatalf("first title edit failed: %s", stderr)
	}
	if code, _, stderr := runBoard(t, root, "edit", "1", "--title", "third"); code != 0 {
		t.Fatalf("second title edit failed: %s", stderr)
	}
	if code, _, stderr := runBoard(t, root, "get", "1"); code != 0 {
		t.Fatalf("get after title edits failed: %s", stderr)
	}
	if len(taskIDs(t, root)) != 1 {
		t.Fatalf("ids = %v", taskIDs(t, root))
	}
}

func TestEditParentClear(t *testing.T) {
	root := initProject(t)
	mustCreate(t, root, "--type", "epic", "e")
	mustCreate(t, root, "--parent", "1", "child")
	code, _, stderr := runBoard(t, root, "edit", "2", "--parent", "")
	if code != 0 {
		t.Fatalf("parent clear failed: %s", stderr)
	}
	summary, err := store.FindByID(bp(root), "2")
	if err != nil {
		t.Fatal(err)
	}
	if summary.Parent != "" {
		t.Fatalf("parent = %q, want empty", summary.Parent)
	}
}

func TestEditParentSet(t *testing.T) {
	root := initProject(t)
	mustCreate(t, root, "--type", "epic", "e")
	mustCreate(t, root, "child")
	code, _, stderr := runBoard(t, root, "edit", "2", "--parent", "1")
	if code != 0 {
		t.Fatalf("parent set failed: %s", stderr)
	}
	summary, err := store.FindByID(bp(root), "2")
	if err != nil {
		t.Fatal(err)
	}
	if summary.Parent != "1" {
		t.Fatalf("parent = %q", summary.Parent)
	}
}

func TestEditNoChanges(t *testing.T) {
	root := initProject(t)
	mustCreate(t, root, "hello")
	code, stdout, stderr := runBoard(t, root, "edit", "1", "--status", "backlog")
	if code != 0 {
		t.Fatalf("edit failed: %s", stderr)
	}
	if !strings.Contains(stdout, "No changes for 1.") {
		t.Fatalf("stdout = %q", stdout)
	}
}

// ---------------------------------------------------------------------------
// list
// ---------------------------------------------------------------------------

func TestListAll(t *testing.T) {
	root := initProject(t)
	mustCreate(t, root, "alpha")
	mustCreate(t, root, "beta")
	code, stdout, stderr := runBoard(t, root, "list")
	if code != 0 {
		t.Fatalf("list failed: %s", stderr)
	}
	if !strings.Contains(stdout, "1") || !strings.Contains(stdout, "2") {
		t.Fatalf("stdout = %q", stdout)
	}
	if !strings.Contains(stdout, "TITLE") {
		t.Fatalf("stdout missing header: %q", stdout)
	}
}

func TestListFilters(t *testing.T) {
	root := initProject(t)
	mustCreate(t, root, "a")
	mustCreate(t, root, "b")
	if code, _, stderr := runBoard(t, root, "edit", "1", "--status", "done"); code != 0 {
		t.Fatalf("edit failed: %s", stderr)
	}
	if code, _, stderr := runBoard(t, root, "edit", "2", "--priority", "high"); code != 0 {
		t.Fatalf("edit failed: %s", stderr)
	}

	code, stdout, stderr := runBoard(t, root, "list", "--status", "done")
	if code != 0 {
		t.Fatalf("list failed: %s", stderr)
	}
	if !strings.Contains(stdout, "1") || strings.Contains(stdout, "2") {
		t.Fatalf("stdout = %q", stdout)
	}

	code, stdout, _ = runBoard(t, root, "list", "--priority", "high")
	if code != 0 {
		t.Fatalf("list failed")
	}
	if !strings.Contains(stdout, "2") || strings.Contains(stdout, "1") {
		t.Fatalf("stdout = %q", stdout)
	}

	code, stdout, _ = runBoard(t, root, "list", "--type", "epic")
	if code != 0 {
		t.Fatalf("list failed")
	}
	if strings.Contains(stdout, "1") {
		t.Fatalf("stdout = %q", stdout)
	}
}

func TestListEmpty(t *testing.T) {
	root := initProject(t)
	code, stdout, stderr := runBoard(t, root, "list")
	if code != 0 {
		t.Fatalf("list failed: %s", stderr)
	}
	if !strings.Contains(stdout, "No tasks found.") {
		t.Fatalf("stdout = %q", stdout)
	}
}

func TestListMarkdown(t *testing.T) {
	root := initProject(t)
	mustCreate(t, root, "alpha")
	mustCreate(t, root, "beta")
	if code, _, stderr := runBoard(t, root, "edit", "2", "--status", "done"); code != 0 {
		t.Fatalf("edit failed: %s", stderr)
	}
	code, stdout, stderr := runBoard(t, root, "list", "--md")
	if code != 0 {
		t.Fatalf("list --md failed: %s", stderr)
	}
	if !strings.Contains(stdout, "# Board — 1 active, 1 archived") {
		t.Fatalf("missing header: %q", stdout)
	}
	if !strings.Contains(stdout, "## backlog") || !strings.Contains(stdout, "## done") {
		t.Fatalf("missing status groups: %q", stdout)
	}
	if !strings.Contains(stdout, "- 1 [medium] alpha") || !strings.Contains(stdout, "- 2 [medium] beta") {
		t.Fatalf("missing task lines: %q", stdout)
	}
}

// ---------------------------------------------------------------------------
// subtasks
// ---------------------------------------------------------------------------

func TestSubtasks(t *testing.T) {
	root := initProject(t)
	mustCreate(t, root, "--type", "epic", "e")
	mustCreate(t, root, "--parent", "1", "one")
	mustCreate(t, root, "--parent", "1", "two")
	code, stdout, stderr := runBoard(t, root, "subtasks", "1")
	if code != 0 {
		t.Fatalf("subtasks failed: %s", stderr)
	}
	if !strings.Contains(stdout, "Subtasks of 1") {
		t.Fatalf("stdout = %q", stdout)
	}
	if !strings.Contains(stdout, "2") || !strings.Contains(stdout, "3") {
		t.Fatalf("stdout = %q", stdout)
	}
	if strings.Contains(stdout, "1") == false {
		t.Fatalf("stdout = %q", stdout)
	}
}

func TestSubtasksNone(t *testing.T) {
	root := initProject(t)
	mustCreate(t, root, "only")
	code, stdout, stderr := runBoard(t, root, "subtasks", "1")
	if code != 0 {
		t.Fatalf("subtasks failed: %s", stderr)
	}
	if !strings.Contains(stdout, "No subtasks for 1.") {
		t.Fatalf("stdout = %q", stdout)
	}
}

func TestSubtasksNotFound(t *testing.T) {
	root := initProject(t)
	code, _, stderr := runBoard(t, root, "subtasks", "99")
	if code == 0 {
		t.Fatalf("missing task should fail")
	}
	if !strings.Contains(stderr, "task not found: 99") {
		t.Fatalf("stderr = %q", stderr)
	}
}
