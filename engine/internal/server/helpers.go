package server

import (
	"encoding/json"
	"net/http"

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