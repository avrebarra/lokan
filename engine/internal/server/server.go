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
	root string
}

// New returns a Server rooted at the given project directory.
func New(root string) *Server {
	return &Server{root: root}
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
	// load all summaries and respond, defaulting to an empty list
	tasks, err := store.LoadAllSummaries(s.root)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if tasks == nil {
		tasks = []types.TaskSummary{}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"tasks": tasks,
		"root":  s.root,
	})
}

func (s *Server) handleTask(w http.ResponseWriter, r *http.Request) {
	// extract the id from the path, then load full detail
	id := strings.TrimPrefix(r.URL.Path, "/api/task/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing id")
		return
	}
	summary, err := store.FindByID(s.root, id)
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
	task, err := CreateTask(s.root, req.Title, req.Type, req.Priority, req.Parent)
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
	summary, err := store.FindByID(s.root, req.ID)
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
		if !contains(types.Statuses, types.Status(req.Value)) {
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

func (s *Server) handleSeed(w http.ResponseWriter, r *http.Request) {
	created, err := SeedDemoData(s.root)
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
