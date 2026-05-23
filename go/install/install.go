package install

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gofrs/flock"
)

// KCL_VERSION is a sentinel used to invalidate cached cdylibs under
// $XDG_CACHE_HOME/kcl/kcl. The "-fork-N" suffix differentiates this
// fork's embedded binaries from the upstream-released ones — bump
// the integer suffix whenever the embedded contents change in a way
// users' cached copies need to be replaced.
//
// -fork-2: slot rename to full rustc target triples; linux-amd64 slot
//          restored to upstream glibc binary (was Chimera-musl in
//          -fork-1). Chimera consumers must build with `-tags chimera`
//          and link against `x86_64-chimera-linux-musl/libkcl.a` via
//          native_chimera_amd64.go.
//
// -fork-3: Phase 2.5 Context cleanup. The kcl_value_plan_to_{json,yaml}
//          C-API functions no longer mirror their output into the
//          context — returned ValueRef::str now own their backing
//          strings independently. The static archive in
//          x86_64-chimera-linux-musl/libkcl.a is rebuilt against
//          kcl-lang commit a08e12b7. The embedded .so slots (glibc,
//          Alpine-musl) still ship upstream's stock binary and are
//          unaffected by this bump.
const KCL_VERSION = "v0.12.3-fork-3"

func getVersion() string {
	return fmt.Sprintf("%s-%s-%s", KCL_VERSION, runtime.GOOS, runtime.GOARCH)
}

func checkVersion(kclVersionDir string) (bool, error) {
	kclVersionPath := filepath.Join(kclVersionDir, "kcl.version")
	_, err := os.Stat(kclVersionPath)

	if os.IsNotExist(err) {
		err := os.MkdirAll(kclVersionDir, 0777)
		if err != nil {
			return false, err
		}
		versionFile, err := os.Create(kclVersionPath)
		defer func() {
			versionFile.Close()
		}()
		return false, err
	}
	version, err := os.ReadFile(kclVersionPath)

	if err != nil {
		return false, err
	}

	return getVersion() == string(version), nil

}

func InstallKcl(installRoot string) error {
	installRoot, err := filepath.Abs(installRoot)
	if err != nil {
		return err
	}
	err = os.MkdirAll(installRoot, 0777)
	if err != nil {
		return err
	}
	// Create a lock file for installing.
	lockFilePath := filepath.Join(installRoot, "install.lock")
	fileLock := flock.New(lockFilePath)

	// Try to obtain a lock with a timeout.
	err = fileLock.Lock()
	if err != nil {
		return err
	}
	// Ensure the lock is released when done.
	defer fileLock.Unlock()

	// Check the lib is installed.
	versionMatched, err := checkVersion(installRoot)
	if err != nil {
		return err
	}

	// Install kcl libs.
	err = installLib(installRoot, "kcl", versionMatched)
	if err != nil {
		return err
	}

	if !versionMatched {
		kclVersionPath := filepath.Join(installRoot, "kcl.version")
		err = os.WriteFile(kclVersionPath, []byte(getVersion()), os.FileMode(os.O_WRONLY|os.O_TRUNC))
		if err != nil {
			return err
		}
	}

	os.Setenv("PATH", os.Getenv("PATH")+string(os.PathListSeparator)+installRoot)

	return nil
}
