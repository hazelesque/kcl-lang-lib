//go:build linux && amd64 && !chimera
// +build linux,amd64,!chimera

package lib

import (
	_ "embed"
)

//go:embed x86_64-unknown-linux-gnu/libkcl.so
var CliLib []byte
