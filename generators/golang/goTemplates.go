package golang

import "embed"

//go:embed templates/*.tmpl templates/**/*.tmpl
var goTemplates embed.FS
