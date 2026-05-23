//go:build chimera
// +build chimera

// On Chimera Linux the native cgo path (`native_chimera_amd64.go`)
// statically links against `x86_64-chimera-linux-musl/libkcl.a`, so the
// embed-then-install runtime extraction path is unused. This file exists
// only to satisfy `var CliLib []byte` and `var ExportLib []byte` symbols
// referenced by `install/install_lib_*.go`; the slices are deliberately
// empty because nothing in the chimera build calls `install.InstallKcl`.

package lib

var CliLib []byte
var ExportLib []byte
