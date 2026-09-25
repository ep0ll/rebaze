package sign

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSignAndVerify(t *testing.T) {
	dir := t.TempDir()
	priv := filepath.Join(dir, "ed25519.key")
	pub := filepath.Join(dir, "ed25519.pub")
	payload := filepath.Join(dir, "plan.json")
	sig := filepath.Join(dir, "plan.json.sig")

	if err := os.WriteFile(payload, []byte(`{"image":"alpine:latest"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := GenerateKeyPair(priv, pub); err != nil {
		t.Fatal(err)
	}
	if _, err := SignFile(payload, priv, sig); err != nil {
		t.Fatal(err)
	}
	if err := VerifyFile(payload, pub, sig); err != nil {
		t.Fatal(err)
	}
}
