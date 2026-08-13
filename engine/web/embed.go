// Package web embeds the built static UI served by the lokan HTTP server.
package web

import "embed"

// FS holds the built Vite app (web/dist), copied here by the build.
// A placeholder dist/index.html is committed so `go build` works before
// the frontend has been built.
//
//go:embed all:dist
var FS embed.FS
