package server

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/avrebarra/lokan/internal/store"
	"github.com/avrebarra/lokan/internal/types"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ResponseDataError{Error: message})
}

// validStatus reports whether v is one of the configured lane ids.
func (s *Server) validStatus(v types.Status) bool {
	cfg, err := store.ReadConfig(s.boardPath)
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

// boardLocation reports the board path relative to the nearest git root and
// that root's directory name, for display in the UI. Falls back to the raw
// path and an empty root when no git root is found.
func (s *Server) boardLocation() (path, root string) {
	abs, err := filepath.Abs(s.boardPath)
	if err != nil {
		return s.boardPath, ""
	}
	for dir := abs; ; {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			rel, err := filepath.Rel(dir, abs)
			if err != nil {
				return s.boardPath, ""
			}
			return filepath.ToSlash(rel), filepath.Base(dir)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return s.boardPath, ""
		}
		dir = parent
	}
}
