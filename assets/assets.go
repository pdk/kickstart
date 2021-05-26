package assets

import (
	"embed"
)

// FS embeds all assets
//go:embed templates/* css/* js/* img/*
var FS embed.FS
