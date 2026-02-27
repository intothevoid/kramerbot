//go:build !prod

package main

import "io/fs"

// staticFiles is nil in dev mode — no embedded frontend.
var staticFiles fs.FS
