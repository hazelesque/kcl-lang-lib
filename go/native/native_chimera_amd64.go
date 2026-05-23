//go:build chimera && amd64
// +build chimera,amd64

// Chimera Linux uses musl libc but is ABI-distinct from other musl distros
// (Alpine, Void) — different libc patches, different distro-identity build
// flags, different rustc target triple (x86_64-chimera-linux-musl vs
// x86_64-unknown-linux-musl). This file is the static-link cgo path that
// pulls in the fork's Chimera-specific libkcl.a; build with `-tags chimera`.

// Static-link libkcl, dynamic-link everything else. Chimera ships libc
// (and libatomic) as dynamic-only by default, so a fully-static link
// (`-static`) fails with `unable to find library -lc`. The explicit
// `-Wl,-Bstatic -lkcl -Wl,-Bdynamic` bracket pins libkcl as static
// while leaving libc/libpthread/libdl/libatomic dynamic.

package native

/*
#cgo LDFLAGS: -L${SRCDIR}/../lib/x86_64-chimera-linux-musl -Wl,-Bstatic -lkcl -Wl,-Bdynamic
#include <stdlib.h>
#include "../include/kcl.h"
*/
import "C"
import (
	"runtime"
	"sync"
	"unsafe"

	"kcl-lang.io/lib/go/api"
	"kcl-lang.io/lib/go/plugin"
)

var libInit sync.Once

var (
	client        *NativeServiceClient
	serviceNew    func(uint64) uintptr
	serviceDelete func(uintptr)
	serviceCall   func(uintptr, string, string, uint, *uint) uintptr
	free          func(uintptr, uint)
)

type NativeServiceClient struct {
	svc uintptr
}

func initClient(pluginAgent uint64) {
	libInit.Do(func() {
		serviceNew = func(agent uint64) uintptr {
			return uintptr(C.kcl_service_new(C.uint64_t(agent)))
		}
		serviceDelete = func(svc uintptr) {
			C.kcl_service_delete(C.uintptr_t(svc))
		}
		serviceCall = func(svc uintptr, method, args string, argsLen uint, outSize *uint) uintptr {
			cSvc := C.KclServiceHandle(svc)
			cMethod := C.CString(method)
			defer C.free(unsafe.Pointer(cMethod))
			cArgs := C.CString(args)
			defer C.free(unsafe.Pointer(cArgs))
			cArgsLen := C.uint32_t(argsLen)
			cOutSize := (*C.uint32_t)(unsafe.Pointer(outSize))

			var cResultPtr *C.uint8_t
			cResultPtr = C.kcl_service_call_with_length(cSvc, cMethod, cArgs, cArgsLen, cOutSize)
			return uintptr(unsafe.Pointer(cResultPtr))
		}
		free = func(ptr uintptr, len uint) {
			cPtr := (*C.uint8_t)(unsafe.Pointer(ptr))
			cLen := C.uint32_t(len)
			C.kcl_free(cPtr, cLen)
		}
		client = &NativeServiceClient{
			svc: serviceNew(pluginAgent),
		}
		runtime.SetFinalizer(client, func(x *NativeServiceClient) {
			if x != nil {
				x.Close()
			}
		})
	})
}

func NewNativeServiceClient() api.ServiceClient {
	return NewNativeServiceClientWithPluginAgent(plugin.GetInvokeJsonProxyPtr())
}

func NewNativeServiceClientWithPluginAgent(pluginAgent uint64) *NativeServiceClient {
	initClient(pluginAgent)
	return client
}

func (x *NativeServiceClient) Close() {
	if x.svc != 0 {
		serviceDelete(x.svc)
		x.svc = 0
	}
}
