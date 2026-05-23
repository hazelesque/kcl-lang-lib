package lib

import (
	_ "embed"
)

//go:embed aarch64-unknown-linux-gnu/libkcl.so
var CliLib []byte
