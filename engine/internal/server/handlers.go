package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/avrebarra/lokan/internal/store"
	"github.com/avrebarra/lokan/internal/types"
)

func (s *Server) HandleTasks(w http.ResponseWriter, r *http.Request) {
	// load all summaries and the lane config, defaulting to empty lists
	tasks, err := store.LoadAllSummaries(s.boardPath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if tasks == nil {
		tasks = []types.TaskSummary{}
	}
	cfg, err := store.ReadConfig(s.boardPath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ResponseDataTasks{
		Tasks:    tasks,
		Statuses: cfg.Statuses,
		Root:     s.boardPath,
	})
}

func (s *Server) HandleTask(w http.ResponseWriter, r *http.Request) {
	// extract the id from the path, then load full detail
	id := strings.TrimPrefix(r.URL.Path, "/api/task/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing id")
		return
	}
	summary, err := store.FindByID(s.boardPath, id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	task, err := store.LoadTask(summary.FilePath)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ResponseDataTask{Task: task})
}

func (s *Server) HandleCreate(w http.ResponseWriter, r *http.Request) {
	// decode the request body
	var req RequestDataCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// create via the shared flow — same validation as the CLI
	task, err := store.CreateTaskFromInput(s.boardPath, req.Title, req.Type, req.Priority, req.Parent, nil)
	if err != nil {
		var ve *store.ValidationError
		if errors.As(err, &ve) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ResponseDataCreate{Task: task})
}

func (s *Server) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	// decode the request body
	var req RequestDataUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// load the target task
	summary, err := store.FindByID(s.boardPath, req.ID)
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
		if !types.Contains(types.Priorities, types.Priority(req.Value)) {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("Invalid priority: %s", req.Value))
			return
		}
		task.Priority = types.Priority(req.Value)
	case "title":
		task.Title = req.Value
	case "type":
		if !types.Contains(types.TaskTypes, types.TaskType(req.Value)) {
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
	writeJSON(w, http.StatusOK, ResponseDataUpdate{Task: task})
}

func (s *Server) HandleMove(w http.ResponseWriter, r *http.Request) {
	// decode the request body
	var req RequestDataMove
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// validate the target lane and the anchor task against the board
	if !s.validStatus(req.Status) {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Invalid status: %s", req.Status))
		return
	}
	summaries, err := store.LoadAllSummaries(s.boardPath)
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
	moved, err := store.MoveTask(s.boardPath, req.ID, req.Status, req.BeforeID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ResponseDataMove{Task: moved})
}

func (s *Server) HandleConfigStatuses(w http.ResponseWriter, r *http.Request) {
	// decode the request body
	var req RequestDataConfigStatuses
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
	cfg, err := store.ReadConfig(s.boardPath)
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
	moved, err := store.UpdateLanes(s.boardPath, renames, req.Statuses)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ResponseDataConfigStatuses{
		Statuses: req.Statuses,
		Moved:    moved,
	})
}

func (s *Server) HandleClear(w http.ResponseWriter, r *http.Request) {
	// decode the request body
	var req RequestDataClear
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// run the requested scope and report how many were deleted
	var deleted int
	var err error
	switch req.Scope {
	case "archived":
		deleted, err = store.ClearArchived(s.boardPath)
	case "all":
		deleted, err = store.ClearAll(s.boardPath)
	default:
		writeError(w, http.StatusBadRequest, "Invalid scope: must be \"archived\" or \"all\"")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ResponseDataClear{Deleted: deleted})
}

func (s *Server) HandleSeed(w http.ResponseWriter, r *http.Request) {
	created, err := SeedDemoData(s.boardPath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ResponseDataSeed{Created: created})
}