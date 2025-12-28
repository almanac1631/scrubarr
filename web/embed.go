package internal

import "embed"

//go:embed all:templates
var Templates embed.FS

//go:embed all:assets
var Assets embed.FS
