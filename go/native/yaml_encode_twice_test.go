// Phase 2.5 regression test, viewed through the Go binding.
//
// Pre-cleanup, the runtime's kcl_value_plan_to_{json,yaml} C-API functions
// mirrored their output into Context::{json,yaml}_result and returned a
// ValueRef::str pointing into context-resident storage. A second call to
// either function would overwrite those fields; any consumer that held the
// first result's `*char` past the second call read corrupted data.
//
// Post-cleanup (kcl-lang Phase 2.5), the returned ValueRef owns its String
// independently and the two calls cannot interfere. This Go test mirrors
// the Rust-side `test_yaml_encode_twice_independent_results` in
// kcl-lang/crates/runner/src/tests.rs; it exercises the runtime through
// the actual Go cgo path so we confirm the fix is observable from the
// host-language consumer Go consumers care about.

package native

import (
	"strings"
	"testing"

	"kcl-lang.io/lib/go/api"
)

func TestYamlEncodeTwiceIndependentResults(t *testing.T) {
	client := NewNativeServiceClient()
	args := &api.ExecProgramArgs{
		KFilenameList: []string{"test.k"},
		KCodeList: []string{
			`import yaml
first = yaml.encode({a = 1})
second = yaml.encode({b = 2})
`,
		},
	}

	result, err := client.ExecProgram(args)
	if err != nil {
		t.Fatalf("ExecProgram failed: %v", err)
	}
	if result.ErrMessage != "" {
		t.Fatalf("unexpected ErrMessage: %s", result.ErrMessage)
	}

	if !strings.Contains(result.YamlResult, "a: 1") {
		t.Errorf("first did not encode {a=1}; yaml_result = %q", result.YamlResult)
	}
	if !strings.Contains(result.YamlResult, "b: 2") {
		t.Errorf("second did not encode {b=2}; yaml_result = %q", result.YamlResult)
	}
	if strings.Count(result.YamlResult, "a: 1") != 1 {
		t.Errorf("expected exactly one 'a: 1', got %d: %q",
			strings.Count(result.YamlResult, "a: 1"), result.YamlResult)
	}
	if strings.Count(result.YamlResult, "b: 2") != 1 {
		t.Errorf("expected exactly one 'b: 2', got %d: %q",
			strings.Count(result.YamlResult, "b: 2"), result.YamlResult)
	}
}
