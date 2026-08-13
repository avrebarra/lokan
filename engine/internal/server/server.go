// Package server implements the lokan HTTP API, matching docs/api.md.
package server

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/avressatelier/lokan/internal/store"
	"github.com/avressatelier/lokan/internal/types"
	"github.com/avressatelier/lokan/web"
)

// Server serves the lokan API for a single project root.
type Server struct {
	board string
}

// New returns a Server serving the given board file.
func New(board string) *Server {
	return &Server{board: board}
}

// Handler returns the HTTP handler implementing the frozen API contract.
func (s *Server) Handler() http.Handler {
	// register the static app, assets, and api routes
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.serveIndex)
	mux.HandleFunc("GET /assets/", s.serveAssets)
	mux.HandleFunc("GET /api/tasks", s.handleTasks)
	mux.HandleFunc("GET /api/task/", s.handleTask)
	mux.HandleFunc("POST /api/create", s.handleCreate)
	mux.HandleFunc("POST /api/update", s.handleUpdate)
	mux.HandleFunc("POST /api/move", s.handleMove)
	mux.HandleFunc("POST /api/config/statuses", s.handleConfigStatuses)
	mux.HandleFunc("POST /api/clear", s.handleClear)
	mux.HandleFunc("POST /api/seed", s.handleSeed)
	return mux
}

func (s *Server) serveIndex(w http.ResponseWriter, r *http.Request) {
	// only the root path serves the app
	if r.URL.Path != "/" {
		writeError(w, http.StatusNotFound, "Not Found")
		return
	}

	// serve the embedded index.html
	index, err := web.FS.ReadFile("dist/index.html")
	if err != nil {
		writeError(w, http.StatusNotFound, "Not Found")
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(index)
}

func (s *Server) serveAssets(w http.ResponseWriter, r *http.Request) {
	// serve bundled assets from the embedded dist
	sub, err := fs.Sub(web.FS, "dist")
	if err != nil {
		writeError(w, http.StatusNotFound, "Not Found")
		return
	}
	http.FileServer(http.FS(sub)).ServeHTTP(w, r)
}

func (s *Server) handleTasks(w http.ResponseWriter, r *http.Request) {
	// load all summaries and the lane config, defaulting to empty lists
	tasks, err := store.LoadAllSummaries(s.board)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if tasks == nil {
		tasks = []types.TaskSummary{}
	}
	cfg, err := store.ReadConfig(s.board)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"tasks":    tasks,
		"statuses": cfg.Statuses,
		"root":     s.board,
	})
}

func (s *Server) handleTask(w http.ResponseWriter, r *http.Request) {
	// extract the id from the path, then load full detail
	id := strings.TrimPrefix(r.URL.Path, "/api/task/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing id")
		return
	}
	summary, err := store.FindByID(s.board, id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	task, err := store.LoadTask(summary.FilePath)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"task": task})
}

type createRequest struct {
	Title    string         `json:"title"`
	Type     types.TaskType `json:"type"`
	Priority types.Priority `json:"priority"`
	Parent   string         `json:"parent"`
}

func (s *Server) handleCreate(w http.ResponseWriter, r *http.Request) {
	// decode the request body
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// validate title and enums before creating
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "Missing title")
		return
	}
	if !contains(types.TaskTypes, req.Type) {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Invalid type: %s", req.Type))
		return
	}
	if !contains(types.Priorities, req.Priority) {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Invalid priority: %s", req.Priority))
		return
	}

	// create and respond with the new task
	task, err := CreateTask(s.board, req.Title, req.Type, req.Priority, req.Parent)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"task": task})
}

type updateRequest struct {
	ID    string `json:"id"`
	Field string `json:"field"`
	Value string `json:"value"`
}

func (s *Server) handleUpdate(w http.ResponseWriter, r *http.Request) {
	// decode the request body
	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// load the target task
	summary, err := store.FindByID(s.board, req.ID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	task, err := store.LoadTask(summary.FilePath)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	// apply the requested field with enum validation
	switch req.Field {
	case "status":
		if !s.validStatus(types.Status(req.Value)) {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("Invalid status: %s", req.Value))
			return
		}
		task.Status = types.Status(req.Value)
	case "priority":
		if !contains(types.Priorities, types.Priority(req.Value)) {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("Invalid priority: %s", req.Value))
			return
		}
		task.Priority = types.Priority(req.Value)
	case "title":
		task.Title = req.Value
	case "type":
		if !contains(types.TaskTypes, types.TaskType(req.Value)) {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("Invalid type: %s", req.Value))
			return
		}
		task.Type = types.TaskType(req.Value)
	case "parent":
		task.Parent = req.Value
	case "tags":
		tags := []string{}
		for _, tag := range strings.Split(req.Value, ",") {
			if tag = strings.TrimSpace(tag); tag != "" {
				tags = append(tags, tag)
			}
		}
		task.Tags = tags
	case "body":
		// bodies must end with a newline so the next board section doesn't
		// glue onto the last line of the serialized block
		task.Body = strings.TrimRight(req.Value, "\n") + "\n"
	default:
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Unknown field: %s", req.Field))
		return
	}

	// persist and respond
	if err := store.WriteTask(task); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"task": task})
}

type moveRequest struct {
	ID       string       `json:"id"`
	Status   types.Status `json:"status"`
	BeforeID string       `json:"beforeId"`
}

func (s *Server) handleMove(w http.ResponseWriter, r *http.Request) {
	// decode the request body
	var req moveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// validate the target lane and the anchor task against the board
	if !s.validStatus(req.Status) {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Invalid status: %s", req.Status))
		return
	}
	summaries, err := store.LoadAllSummaries(s.board)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	foundID := false
	foundBefore := false
	for _, t := range summaries {
		if t.ID == req.ID {
			foundID = true
		}
		if t.ID == req.BeforeID && req.BeforeID != "" && t.Status == req.Status {
			foundBefore = true
		}
	}
	if !foundID {
		writeError(w, http.StatusNotFound, fmt.Sprintf("task not found: %s", req.ID))
		return
	}
	if req.BeforeID != "" && !foundBefore {
		writeError(w, http.StatusBadRequest, "beforeId must be a task in the target lane")
		return
	}

	// apply the move and respond
	moved, err := store.MoveTask(s.board, req.ID, req.Status, req.BeforeID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"task": moved})
}

// validStatus reports whether v is one of the configured lane ids.
func (s *Server) validStatus(v types.Status) bool {
	cfg, err := store.ReadConfig(s.board)
	if err != nil {
		return false
	}
	for _, st := range cfg.Statuses {
		if st.ID == v {
			return true
		}
	}
	return false
}

type configStatusesRequest struct {
	Statuses []types.StatusDef `json:"statuses"`
}

func (s *Server) handleConfigStatuses(w http.ResponseWriter, r *http.Request) {
	// decode the request body
	var req configStatusesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// validate: at least one lane, non-empty unique ids
	if len(req.Statuses) == 0 {
		writeError(w, http.StatusBadRequest, "At least one status is required")
		return
	}
	seen := map[types.Status]bool{}
	for _, st := range req.Statuses {
		if st.ID == "" {
			writeError(w, http.StatusBadRequest, "Status id cannot be empty")
			return
		}
		if seen[st.ID] {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("Duplicate status: %s", st.ID))
			return
		}
		seen[st.ID] = true
	}

	// diff against the current lanes: positional renames first, then moves
	// for removed lanes (tasks land in the leftmost remaining lane)
	cfg, err := store.ReadConfig(s.board)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	oldSet := map[types.Status]bool{}
	for _, st := range cfg.Statuses {
		oldSet[st.ID] = true
	}
	renamedFrom := map[types.Status]bool{}
	var renames [][2]types.Status
	for i := 0; i < len(cfg.Statuses) && i < len(req.Statuses); i++ {
		o, n := cfg.Statuses[i].ID, req.Statuses[i].ID
		if o == n || oldSet[n] || seen[o] {
			continue
		}
		renames = append(renames, [2]types.Status{o, n})
		renamedFrom[o] = true
	}
	leftmost := req.Statuses[0].ID
	for _, st := range cfg.Statuses {
		if !seen[st.ID] && !renamedFrom[st.ID] {
			renames = append(renames, [2]types.Status{st.ID, leftmost})
		}
	}

	// apply renames and persist the new lane set in one atomic rewrite
	moved, err := store.UpdateLanes(s.board, renames, req.Statuses)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"statuses": req.Statuses,
		"moved":    moved,
	})
}

type clearRequest struct {
	Scope string `json:"scope"`
}

func (s *Server) handleClear(w http.ResponseWriter, r *http.Request) {
	// decode the request body
	var req clearRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// run the requested scope and report how many were deleted
	var deleted int
	var err error
	switch req.Scope {
	case "archived":
		deleted, err = store.ClearArchived(s.board)
	case "all":
		deleted, err = store.ClearAll(s.board)
	default:
		writeError(w, http.StatusBadRequest, "Invalid scope: must be \"archived\" or \"all\"")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"deleted": deleted})
}

func (s *Server) handleSeed(w http.ResponseWriter, r *http.Request) {
	created, err := SeedDemoData(s.board)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"created": created})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func contains[T ~string](list []T, v T) bool {
	for _, item := range list {
		if item == v {
			return true
		}
	}
	return false
}
