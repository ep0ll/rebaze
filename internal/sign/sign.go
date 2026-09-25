package sign

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Envelope struct {
	Algorithm string    `json:"algorithm"`
	KeyID     string    `json:"keyId"`
	Signature string    `json:"signature"`
	Created   time.Time `json:"created"`
	Path      string    `json:"path"`
}

func GenerateKeyPair(privPath, pubPath string) error {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	if err := os.WriteFile(privPath, []byte(base64.StdEncoding.EncodeToString(priv)), 0o600); err != nil {
		return err
	}
	return os.WriteFile(pubPath, []byte(base64.StdEncoding.EncodeToString(pub)), 0o644)
}

func loadKey(path string) ([]byte, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(string(raw))
}

func SignFile(path, privPath, outPath string) (*Envelope, error) {
	payload, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	priv, err := loadKey(privPath)
	if err != nil {
		return nil, fmt.Errorf("read private key: %w", err)
	}
	if len(priv) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid ed25519 private key size")
	}
	sig := ed25519.Sign(ed25519.PrivateKey(priv), payload)
	env := &Envelope{
		Algorithm: "ed25519",
		KeyID:     base64.StdEncoding.EncodeToString(ed25519.PrivateKey(priv).Public().(ed25519.PublicKey)),
		Signature: base64.StdEncoding.EncodeToString(sig),
		Created:   time.Now().UTC(),
		Path:      path,
	}
	data, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(outPath, append(data, '\n'), 0o644); err != nil {
		return nil, err
	}
	return env, nil
}

func VerifyFile(path, pubPath, sigPath string) error {
	payload, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	pub, err := loadKey(pubPath)
	if err != nil {
		return fmt.Errorf("read public key: %w", err)
	}
	if len(pub) != ed25519.PublicKeySize {
		return fmt.Errorf("invalid ed25519 public key size")
	}
	var env Envelope
	data, err := os.ReadFile(sigPath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &env); err != nil {
		return fmt.Errorf("parse signature envelope: %w", err)
	}
	sig, err := base64.StdEncoding.DecodeString(env.Signature)
	if err != nil {
		return err
	}
	if !ed25519.Verify(ed25519.PublicKey(pub), payload, sig) {
		return fmt.Errorf("signature verification failed")
	}
	return nil
}
