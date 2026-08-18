// Package server implements the lokan HTTP API, matching docs/api.md.
package server

import (
	"io/fs"
	"net/http"

	"github.com/avrebarra/lokan/web"
)

// Server serves the lokan API for a single project root.
type Server struct {
	boardPath string
}

// New returns a Server serving the given board file.
func New(boardPath string) *Server {
	return &Server{boardPath: boardPath}
}

// Handler returns the HTTP handler implementing the frozen API contract.
func (s *Server) Handler() http.Handler {
	// register the static app, assets, and api routes
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.serveIndex)
	mux.HandleFunc("GET /assets/", s.serveAssets)
	mux.HandleFunc("GET /api/tasks", s.HandleTasks)
	mux.HandleFunc("GET /api/task/", s.HandleTask)
	mux.HandleFunc("POST /api/create", s.HandleCreate)
	mux.HandleFunc("POST /api/update", s.HandleUpdate)
	mux.HandleFunc("POST /api/move", s.HandleMove)
	mux.HandleFunc("POST /api/config/statuses", s.HandleConfigStatuses)
	mux.HandleFunc("POST /api/clear", s.HandleClear)
	mux.HandleFunc("POST /api/delete", s.HandleDelete)
	mux.HandleFunc("POST /api/seed", s.HandleSeed)
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
