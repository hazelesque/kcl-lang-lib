package lib

import (
	_ "embed"
)

//go:embed x86_64-apple-darwin/libkcl.dylib
var CliLib []byte
