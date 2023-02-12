package main

import (
	"embed"

	"github.com/dinkur/dinkur-desktop/cmd"
	"github.com/dinkur/dinkur-desktop/pkg/app"
)

//go:embed all:frontend/build
var assets embed.FS

//go:embed icons/dinkur-small-64.png
var iconBytes []byte

func main() {
	app.IconBytes = iconBytes
	app.Assets = assets

	cmd.Execute()
}
