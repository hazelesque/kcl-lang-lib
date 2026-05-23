package lib

import (
	_ "embed"
)

//go:embed aarch64-pc-windows-msvc/kcl.dll
var CliLib []byte

// NOTE: upstream embedded the amd64 .lib here. The .lib content is
// byte-identical across both windows slots so this isn't a visible bug,
// but the path was wrong; corrected to match the file's GOARCH.
//go:embed aarch64-pc-windows-msvc/kcl.lib
var ExportLib []byte
