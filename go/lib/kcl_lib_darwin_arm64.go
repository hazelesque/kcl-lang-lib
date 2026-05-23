package lib

import (
	_ "embed"
)

//go:embed aarch64-apple-darwin/libkcl.dylib
var CliLib []byte
