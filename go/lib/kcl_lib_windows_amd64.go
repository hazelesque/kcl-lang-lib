package lib

import (
	_ "embed"
)

//go:embed x86_64-pc-windows-msvc/kcl.dll
var CliLib []byte

//go:embed x86_64-pc-windows-msvc/kcl.lib
var ExportLib []byte
