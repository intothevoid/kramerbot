//go:build prod

package main

import (
	"embed"
	"io/fs"
)

//go:embed all:frontend/dist
var embeddedFiles embed.FS

// staticFiles is the sub-filesystem rooted at frontend/dist.
var staticFiles, _ = fs.Sub(embeddedFiles, "frontend/dist")
