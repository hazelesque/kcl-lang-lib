package lib

// Side-effect imports of the per-slot sub-packages so `go mod` keeps
// the slot directories in the module graph. The actual library bytes
// are consumed via //go:embed in the per-platform kcl_lib_*.go files;
// these imports just keep the dirs visible to module tooling.
import (
	_ "kcl-lang.io/lib/go/lib/aarch64-apple-darwin"
	_ "kcl-lang.io/lib/go/lib/aarch64-pc-windows-msvc"
	_ "kcl-lang.io/lib/go/lib/aarch64-unknown-linux-gnu"
	_ "kcl-lang.io/lib/go/lib/x86_64-apple-darwin"
	_ "kcl-lang.io/lib/go/lib/x86_64-pc-windows-msvc"
	_ "kcl-lang.io/lib/go/lib/x86_64-unknown-linux-gnu"
)
