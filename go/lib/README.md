# Embedded KCL runtime libraries

Per-platform prebuilt KCL runtime libraries consumed by Go's cgo wrappers.
This fork uses full **rustc target triples** for slot names instead of
upstream's Go-style `linux-amd64` / `linux-musl-amd64` shorthand, because
the shorthand collides distinct ABIs (Chimera musl is not Alpine musl, and
upstream's `linux-musl-*` slots contain Alpine-shaped builds).

## Slot scheme

| Slot                                 | Library             | Linker     | Go consumer                                |
| ------------------------------------ | ------------------- | ---------- | ------------------------------------------ |
| `x86_64-unknown-linux-gnu/`          | `libkcl.so`         | glibc      | `native_nonmusl.go` (default, dlopen)      |
| `aarch64-unknown-linux-gnu/`         | `libkcl.so`         | glibc      | `native_nonmusl.go`                        |
| `x86_64-unknown-linux-musl/`         | `libkcl.a`          | musl (Alpine-shaped) | `native_musl_amd64.go` (build with `-tags musl`) |
| `aarch64-unknown-linux-musl/`        | `libkcl.a`          | musl (Alpine-shaped) | `native_musl_arm64.go` (build with `-tags musl`) |
| `x86_64-chimera-linux-musl/`         | `libkcl.a`          | Chimera musl | `native_chimera_amd64.go` (build with `-tags chimera`) |
| `x86_64-apple-darwin/`               | `libkcl.dylib`      | system     | `native_nonmusl.go`                        |
| `aarch64-apple-darwin/`              | `libkcl.dylib`      | system     | `native_nonmusl.go`                        |
| `x86_64-pc-windows-msvc/`            | `kcl.dll` + `kcl.lib` | MSVC     | `native_nonmusl.go`                        |
| `aarch64-pc-windows-msvc/`           | `kcl.dll` + `kcl.lib` | MSVC     | `native_nonmusl.go`                        |

## Build tag conventions

| Build invocation                      | Active cgo file                | Library source                                                    |
| ------------------------------------- | ------------------------------ | ----------------------------------------------------------------- |
| `go build ./...` (default)            | `native_nonmusl.go`            | Embed-then-dlopen; reads `lib.CliLib` and extracts to `$XDG_CACHE_HOME/kcl/kcl/`. Per-platform shared library from the embedded blob. |
| `go build -tags musl ./...`           | `native_musl_{amd64,arm64}.go` | Static-link via cgo against `x86_64-unknown-linux-musl/libkcl.a`. Alpine-shaped musl. |
| `go build -tags chimera ./...`        | `native_chimera_amd64.go`      | Static-link via cgo against `x86_64-chimera-linux-musl/libkcl.a`. Chimera's distinct musl ABI. |

The `chimera` and `musl` tags are mutually exclusive; each gates the others
out. The default `!musl && !chimera` path uses the embed-then-dlopen
machinery (`install/install.go` + `loader.go`).

## Chimera-specific notes

Chimera Linux uses musl but its rustc target triple is
`x86_64-chimera-linux-musl`, distinct from `x86_64-unknown-linux-musl`
(the generic / Alpine-shaped target). The libcs differ in subtle ways
(libc patches, distro-identity build flags, occasionally symbol coverage)
that mean a binary built for one is not safely consumable by the other.

The `x86_64-chimera-linux-musl/libkcl.a` slot is populated from a local
build of the fork (`cargo build -p kcl-lib --release` on a Chimera host
produces both `libkcl.so` and `libkcl.a` under `target/release/`; the `.a`
goes in this slot).

## Cache invalidation

`go/install/install.go::KCL_VERSION` is the cache-busting sentinel for
the embed-then-dlopen path. Bump the `-fork-N` integer whenever embedded
content changes in a way users' cached copies need to be replaced. Slot
renames count.

## Known upstream-content weirdness

- `windows-amd64/` and `windows-arm64/` (now `x86_64-pc-windows-msvc/` and
  `aarch64-pc-windows-msvc/`) contained byte-identical `kcl.dll` and
  `kcl.lib` in the upstream release. Likely an upstream CI misconfiguration
  (the arm64 build is actually an amd64 binary). Not corrected here — out
  of scope and the fork doesn't ship Windows consumers anyway.
- `kcl_lib_windows_arm64.go` upstream embedded `windows-amd64/kcl.lib`
  (note the `amd64`) for the arm64 build. Corrected during the rename
  because the embed path had to change anyway; documented in the file's
  comment.
