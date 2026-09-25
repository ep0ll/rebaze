package apply

import (
	"testing"
)

func TestMergeAndUnsetEnv(t *testing.T) {
	got := mergeEnv([]string{"PATH=/bin", "FOO=1"}, map[string]string{"FOO": "2", "BAR": "3"})
	want := map[string]string{"PATH": "/bin", "FOO": "2", "BAR": "3"}
	gotMap := map[string]string{}
	for _, kv := range got {
		k, v, _ := splitOnce(kv)
		gotMap[k] = v
	}
	for k, v := range want {
		if gotMap[k] != v {
			t.Fatalf("env %s = %q, want %q", k, gotMap[k], v)
		}
	}

	cleared := unsetEnv(got, []string{"FOO"})
	for _, kv := range cleared {
		k, _, _ := splitOnce(kv)
		if k == "FOO" {
			t.Fatal("FOO should have been removed")
		}
	}
}

func splitOnce(kv string) (string, string, bool) {
	for i := 0; i < len(kv); i++ {
		if kv[i] == '=' {
			return kv[:i], kv[i+1:], true
		}
	}
	return kv, "", false
}
