package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/avressatelier/lokan/internal/id"
	"github.com/avressatelier/lokan/internal/store"
	"github.com/avressatelier/lokan/internal/types"
)

func newTestProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := id.WriteConfig(root, types.LokanConfig{Counter: 0, Version: "1"}); err != nil {
		t.Fatal(err)
	}
	return root
}

func createTestTask(t *testing.T, root string, id, title string, status types.Status, priority types.Priority) {
	t.Helper()
	fm := types.TaskFrontmatter{
		ID:       id,
		Title:    title,
		Type:     types.TypeTask,
		Status:   status,
		Priority: priority,
		Created:  "2024-01-01",
		Updated:  "2024-01-01",
	}
	if _, err := store.CreateTask(root, fm, ""); err != nil {
		t.Fatal(err)
	}
}

func doRequest(t *testing.T, h http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		reader = strings.NewReader(string(raw))
	}
	req := httptest.NewRequest(method, path, reader)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

// ---------------------------------------------------------------------------
// GET /api/tasks
// ---------------------------------------------------------------------------

func TestTasksEmpty(t *testing.T) {
	root := newTestProject(t)
	rec := doRequest(t, New(root).Handler(), "GET", "/api/tasks", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("content-type = %q, want application/json", ct)
	}
	var resp struct {
		Tasks []types.TaskSummary `json:"tasks"`
		Root  string              `json:"root"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Tasks == nil || len(resp.Tasks) != 0 {
		t.Fatalf("tasks = %#v, want empty non-nil slice", resp.Tasks)
	}
	if resp.Root != root {
		t.Fatalf("root = %q, want %q", resp.Root, root)
	}
}

func TestTasksSeeded(t *testing.T) {
	root := newTestProject(t)
	createTestTask(t, root, "1", "Alpha", types.StatusTodo, types.PriorityMedium)
	createTestTask(t, root, "2", "Beta", types.StatusDone, types.PriorityHigh)

	rec := doRequest(t, New(root).Handler(), "GET", "/api/tasks", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var resp struct {
		Tasks []types.TaskSummary `json:"tasks"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Tasks) != 2 {
		t.Fatalf("len(tasks) = %d, want 2", len(resp.Tasks))
	}
	if resp.Tasks[0].ID != "1" || resp.Tasks[1].ID != "2" {
		t.Fatalf("unexpected task order: %+v", resp.Tasks)
	}
}

// ---------------------------------------------------------------------------
// GET /api/task/:id
// ---------------------------------------------------------------------------

func TestTaskFound(t *testing.T) {
	root := newTestProject(t)
	createTestTask(t, root, "1", "Alpha", types.StatusTodo, types.PriorityMedium)

	rec := doRequest(t, New(root).Handler(), "GET", "/api/task/1", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var resp struct {
		Task types.Task `json:"task"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Task.ID != "1" {
		t.Fatalf("task id = %q, want 1", resp.Task.ID)
	}
	if resp.Task.Title != "Alpha" {
		t.Fatalf("task title = %q, want Alpha", resp.Task.Title)
	}
	if resp.Task.Body == "" {
		t.Fatal("task body should be non-empty")
	}
	if !strings.Contains(resp.Task.FilePath, "1") {
		t.Fatalf("task filePath = %q, want it to contain 1", resp.Task.FilePath)
	}
}

func TestTaskNotFound(t *testing.T) {
	root := newTestProject(t)
	rec := doRequest(t, New(root).Handler(), "GET", "/api/task/nope", nil)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	var resp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Error == "" {
		t.Fatal("error message should not be empty")
	}
}

// ---------------------------------------------------------------------------
// POST /api/update
// ---------------------------------------------------------------------------

func TestUpdateStatus(t *testing.T) {
	root := newTestProject(t)
	createTestTask(t, root, "1", "Alpha", types.StatusTodo, types.PriorityMedium)

	rec := doRequest(t, New(root).Handler(), "POST", "/api/update", map[string]string{
		"id": "1", "field": "status", "value": "done",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var resp struct {
		Task types.Task `json:"task"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Task.Status != types.StatusDone {
		t.Fatalf("task status = %q, want done", resp.Task.Status)
	}
	assertStoredField(t, root, "1", func(t types.Task) bool { return t.Status == types.StatusDone })
}

// ---------------------------------------------------------------------------
// POST /api/move
// ---------------------------------------------------------------------------

func taskIDs(t *testing.T, h http.Handler) []string {
	t.Helper()
	rec := doRequest(t, h, "GET", "/api/tasks", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var resp struct {
		Tasks []types.TaskSummary `json:"tasks"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	ids := make([]string, len(resp.Tasks))
	for i, t := range resp.Tasks {
		ids[i] = t.ID
	}
	return ids
}

func TestMoveReordersWithinLane(t *testing.T) {
	root := newTestProject(t)
	h := New(root).Handler()
	createTestTask(t, root, "1", "Alpha", types.StatusTodo, types.PriorityMedium)
	createTestTask(t, root, "2", "Beta", types.StatusTodo, types.PriorityMedium)
	createTestTask(t, root, "3", "Gamma", types.StatusTodo, types.PriorityMedium)

	rec := doRequest(t, h, "POST", "/api/move", map[string]string{"id": "1", "status": "todo", "beforeId": "3"})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (%s)", rec.Code, rec.Body.String())
	}
	if got, want := taskIDs(t, h), []string{"2", "1", "3"}; !slices.Equal(got, want) {
		t.Fatalf("order = %v, want %v", got, want)
	}
}

func TestMoveAcrossLanes(t *testing.T) {
	root := newTestProject(t)
	h := New(root).Handler()
	createTestTask(t, root, "1", "Alpha", types.StatusTodo, types.PriorityMedium)
	createTestTask(t, root, "2", "Beta", types.StatusTodo, types.PriorityMedium)

	rec := doRequest(t, h, "POST", "/api/move", map[string]string{"id": "1", "status": "done", "beforeId": ""})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (%s)", rec.Code, rec.Body.String())
	}
	if got, want := taskIDs(t, h), []string{"2", "1"}; !slices.Equal(got, want) {
		t.Fatalf("order = %v, want %v", got, want)
	}
	var resp struct {
		Task types.Task `json:"task"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Task.Status != types.StatusDone {
		t.Fatalf("status = %s, want done", resp.Task.Status)
	}
}

func TestMoveInvalidStatus(t *testing.T) {
	root := newTestProject(t)
	createTestTask(t, root, "1", "Alpha", types.StatusTodo, types.PriorityMedium)

	rec := doRequest(t, New(root).Handler(), "POST", "/api/move", map[string]string{"id": "1", "status": "bogus"})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	assertError(t, rec, "Invalid status: bogus")
}

func TestMoveUnknownTask(t *testing.T) {
	root := newTestProject(t)
	createTestTask(t, root, "1", "Alpha", types.StatusTodo, types.PriorityMedium)

	rec := doRequest(t, New(root).Handler(), "POST", "/api/move", map[string]string{"id": "99", "status": "todo"})
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	assertError(t, rec, "task not found: 99")
}

func TestMoveAnchorMustBeInTargetLane(t *testing.T) {
	root := newTestProject(t)
	createTestTask(t, root, "1", "Alpha", types.StatusTodo, types.PriorityMedium)

	rec := doRequest(t, New(root).Handler(), "POST", "/api/move", map[string]string{"id": "1", "status": "done", "beforeId": "1"})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	assertError(t, rec, "beforeId must be a task in the target lane")
}

func TestUpdateInvalidStatus(t *testing.T) {
	root := newTestProject(t)
	createTestTask(t, root, "1", "Alpha", types.StatusTodo, types.PriorityMedium)

	rec := doRequest(t, New(root).Handler(), "POST", "/api/update", map[string]string{
		"id": "1", "field": "status", "value": "bogus",
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	assertError(t, rec, "Invalid status: bogus")
}

func TestUpdateInvalidPriority(t *testing.T) {
	root := newTestProject(t)
	createTestTask(t, root, "1", "Alpha", types.StatusTodo, types.PriorityMedium)

	rec := doRequest(t, New(root).Handler(), "POST", "/api/update", map[string]string{
		"id": "1", "field": "priority", "value": "bogus",
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	assertError(t, rec, "Invalid priority: bogus")
}

func TestUpdateUnknownField(t *testing.T) {
	root := newTestProject(t)
	createTestTask(t, root, "1", "Alpha", types.StatusTodo, types.PriorityMedium)

	rec := doRequest(t, New(root).Handler(), "POST", "/api/update", map[string]string{
		"id": "1", "field": "wat", "value": "x",
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	assertError(t, rec, "Unknown field: wat")
}

func TestUpdateTitle(t *testing.T) {
	root := newTestProject(t)
	createTestTask(t, root, "1", "Alpha", types.StatusTodo, types.PriorityMedium)

	rec := doRequest(t, New(root).Handler(), "POST", "/api/update", map[string]string{
		"id": "1", "field": "title", "value": "Renamed Alpha",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var resp struct {
		Task types.Task `json:"task"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Task.Title != "Renamed Alpha" {
		t.Fatalf("task title = %q, want Renamed Alpha", resp.Task.Title)
	}
	assertStoredField(t, root, "1", func(t types.Task) bool { return t.Title == "Renamed Alpha" })
}

func TestUpdateType(t *testing.T) {
	root := newTestProject(t)
	createTestTask(t, root, "1", "Alpha", types.StatusTodo, types.PriorityMedium)

	rec := doRequest(t, New(root).Handler(), "POST", "/api/update", map[string]string{
		"id": "1", "field": "type", "value": "bug",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var resp struct {
		Task types.Task `json:"task"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Task.Type != types.TypeBug {
		t.Fatalf("task type = %q, want bug", resp.Task.Type)
	}
	if resp.Task.ID != "1" {
		t.Fatalf("task id changed to %q, want 1", resp.Task.ID)
	}
	assertStoredField(t, root, "1", func(t types.Task) bool { return t.Type == types.TypeBug })
}

func TestUpdateInvalidType(t *testing.T) {
	root := newTestProject(t)
	createTestTask(t, root, "1", "Alpha", types.StatusTodo, types.PriorityMedium)

	rec := doRequest(t, New(root).Handler(), "POST", "/api/update", map[string]string{
		"id": "1", "field": "type", "value": "bogus",
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	assertError(t, rec, "Invalid type: bogus")
}

func TestUpdateParent(t *testing.T) {
	root := newTestProject(t)
	createTestTask(t, root, "1", "Alpha", types.StatusTodo, types.PriorityMedium)

	rec := doRequest(t, New(root).Handler(), "POST", "/api/update", map[string]string{
		"id": "1", "field": "parent", "value": "2",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	assertStoredField(t, root, "1", func(t types.Task) bool { return t.Parent == "2" })

	rec = doRequest(t, New(root).Handler(), "POST", "/api/update", map[string]string{
		"id": "1", "field": "parent", "value": "",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	assertStoredField(t, root, "1", func(t types.Task) bool { return t.Parent == "" })
}

func TestUpdateTags(t *testing.T) {
	root := newTestProject(t)
	createTestTask(t, root, "1", "Alpha", types.StatusTodo, types.PriorityMedium)

	rec := doRequest(t, New(root).Handler(), "POST", "/api/update", map[string]string{
		"id": "1", "field": "tags", "value": " alpha, beta ,, gamma ",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	want := []string{"alpha", "beta", "gamma"}
	assertStoredField(t, root, "1", func(t types.Task) bool {
		return slices.Equal(t.Tags, want)
	})
}

func TestUpdateBody(t *testing.T) {
	root := newTestProject(t)
	createTestTask(t, root, "1", "Alpha", types.StatusTodo, types.PriorityMedium)

	rec := doRequest(t, New(root).Handler(), "POST", "/api/update", map[string]string{
		"id": "1", "field": "body", "value": "Round trip notes\nwith a second line",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var resp struct {
		Task types.Task `json:"task"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Task.Body != "Round trip notes\nwith a second line\n" {
		t.Fatalf("task body = %q, want the round-tripped notes with a trailing newline", resp.Task.Body)
	}
	assertStoredField(t, root, "1", func(t types.Task) bool {
		return t.Body == "Round trip notes\nwith a second line"
	})
}

// ---------------------------------------------------------------------------
// POST /api/seed
// ---------------------------------------------------------------------------

func TestSeedCreatesDemoTasks(t *testing.T) {
	root := newTestProject(t)
	h := New(root).Handler()

	rec := doRequest(t, h, "POST", "/api/seed", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var resp struct {
		Created int `json:"created"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Created != 28 {
		t.Fatalf("created = %d, want 28", resp.Created)
	}

	rec = doRequest(t, h, "GET", "/api/tasks", nil)
	var tasks struct {
		Tasks []types.TaskSummary `json:"tasks"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(tasks.Tasks) != 28 {
		t.Fatalf("len(tasks) = %d, want 28", len(tasks.Tasks))
	}
}

// ---------------------------------------------------------------------------
// GET / and unknown routes
// ---------------------------------------------------------------------------

func TestIndexServesEmbed(t *testing.T) {
	root := newTestProject(t)
	rec := doRequest(t, New(root).Handler(), "GET", "/", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("content-type = %q, want text/html", ct)
	}
	if !strings.Contains(rec.Body.String(), "lokan") {
		t.Fatalf("body does not contain lokan: %q", rec.Body.String())
	}
}

func TestUnknownRoute(t *testing.T) {
	root := newTestProject(t)
	h := New(root).Handler()
	for _, path := range []string{"/api/nope", "/somewhere"} {
		rec := doRequest(t, h, "GET", path, nil)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("GET %s status = %d, want 404", path, rec.Code)
		}
	}
}

func TestTaskMissingID(t *testing.T) {
	root := newTestProject(t)
	rec := doRequest(t, New(root).Handler(), "GET", "/api/task/", nil)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	assertError(t, rec, "Missing id")
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func assertError(t *testing.T, rec *httptest.ResponseRecorder, want string) {
	t.Helper()
	var resp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Error != want {
		t.Fatalf("error = %q, want %q", resp.Error, want)
	}
}

func assertStoredField(t *testing.T, root, id string, check func(types.Task) bool) {
	t.Helper()
	summary, err := store.FindByID(root, id)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	task, err := store.LoadTask(summary.FilePath)
	if err != nil {
		t.Fatalf("LoadTask: %v", err)
	}
	if !check(task) {
		t.Fatalf("stored task %s does not match expected state: %+v", id, task)
	}
}

// ---------------------------------------------------------------------------
// POST /api/create
// ---------------------------------------------------------------------------

func TestCreateTaskEndpoint(t *testing.T) {
	root := newTestProject(t)
	h := New(root).Handler()

	rec := doRequest(t, h, "POST", "/api/create", map[string]any{
		"title":    "Hello from API",
		"type":     "task",
		"priority": "high",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Task types.Task `json:"task"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Task.Title != "Hello from API" || resp.Task.Status != types.StatusTodo {
		t.Fatalf("unexpected task: %+v", resp.Task)
	}
	// persisted to disk
	assertStoredField(t, root, resp.Task.ID, func(task types.Task) bool {
		return task.Title == "Hello from API" && task.Priority == types.PriorityHigh
	})
}

func TestCreateTaskValidation(t *testing.T) {
	root := newTestProject(t)
	h := New(root).Handler()

	cases := []struct {
		name string
		body map[string]any
		want string
	}{
		{"missing title", map[string]any{"type": "task", "priority": "medium"}, "Missing title"},
		{"bad type", map[string]any{"title": "x", "type": "nope", "priority": "medium"}, "Invalid type: nope"},
		{"bad priority", map[string]any{"title": "x", "type": "task", "priority": "nope"}, "Invalid priority: nope"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := doRequest(t, h, "POST", "/api/create", tc.body)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400", rec.Code)
			}
			assertError(t, rec, tc.want)
		})
	}
}

func TestUpdateTaskNotFound(t *testing.T) {
	root := newTestProject(t)
	rec := doRequest(t, New(root).Handler(), "POST", "/api/update", map[string]any{
		"id": "999", "field": "status", "value": "done",
	})
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestAssetsServed(t *testing.T) {
	root := newTestProject(t)
	h := New(root).Handler()
	for _, path := range []string{"/assets/nonexistent.js", "/assets/nonexistent.css"} {
		rec := doRequest(t, h, "GET", path, nil)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("GET %s status = %d, want 404", path, rec.Code)
		}
	}
}
