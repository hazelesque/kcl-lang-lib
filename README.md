# KCL Language Bindings (private fork)

Bindings to the [KCL language core](../kcl-lang) for Rust (primary) and Go.

This is a private fork of [`kcl-lang/lib`](https://github.com/kcl-lang/lib). The other host-language bindings (Python, Node.js, Java, Kotlin, C++, .NET, Swift, Lua, Zig, WASM) have been removed; the C ABI surface is retained because Go's cgo path depends on it.

## Bindings

### Rust

```rust
use kcl_lang::*;
use anyhow::Result;

fn main() -> Result<()> {
    let api = API::default();
    let args = &ExecProgramArgs {
        k_filename_list: vec!["main.k".to_string()],
        k_code_list: vec!["a = 1".to_string()],
        ..Default::default()
    };
    let exec_result = api.exec_program(args)?;
    println!("{}", exec_result.yaml_result);
    Ok(())
}
```

The Rust binding re-exports `kcl_api` and friends from the local `../kcl-lang` checkout via path deps in `Cargo.toml`. There is no version-pinning to GitHub.

A cleaner embedding surface returning structured `ValueRef` instead of YAML/JSON strings will land via the `kcl-lang/crates/embed/` crate; see the restructuring plan.

### Go

```go
package main

import (
    "fmt"

    "kcl-lang.io/lib/go/api"
    "kcl-lang.io/lib/go/native"
)

func main() {
    client := native.NewNativeServiceClient()
    result, err := client.ExecProgram(&api.ExecProgramArgs{
        KFilenameList: []string{"main.k"},
        KCodeList:     []string{"a = 1"},
    })
    if err != nil {
        panic(err)
    }
    fmt.Println(result.YamlResult)
}
```

Go uses `purego` + cgo to load `libkcl.{so,dylib,dll}` (built from `../kcl-lang/crates/lib`). Pre-built per-platform binaries are no longer bundled — the local fork's cdylib is consumed directly.

### C ABI

`c/include/kcl_ffi.h` defines the single C entry point (`call_native`) used by Go's cgo path. It is not a user-facing binding in this fork; it is the FFI moat that Go consumes.

## Building

1. Build the fork's cdylib:
   ```sh
   cd ../kcl-lang
   cargo build -p kcl-lib --release
   ```
   Produces `../kcl-lang/target/release/libkcl.{so,dylib,dll}`.

2. Build the Rust binding:
   ```sh
   cargo build --workspace
   ```

3. Build/test Go:
   ```sh
   cd go && go build ./...
   ```
   Go's cgo wrapper loads the cdylib from a known path; see `go/native/` for the lookup logic.

## License

Apache-2.0 (preserved from upstream).
