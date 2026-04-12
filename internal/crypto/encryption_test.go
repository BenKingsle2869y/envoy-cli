package crypto_test

import (
	"testing"

	"github.com/envoy-cli/envoy/internal/crypto"
)

func TestDeriveKey(t *testing.T) {
	key := crypto.DeriveKey("my-secret-passphrase")
	if len(key) != 32 {
		t.Fatalf("expected key length 32, got %d", len(key))
	}

	// Same passphrase should produce the same key
	key2 := crypto.DeriveKey("my-secret-passphrase")
	for i := range key {
		if key[i] != key2[i] {
			t.Fatal("expected deterministic key derivation")
		}
	}
}

func TestEncryptDecrypt(t *testing.T) {
	passphrase := "test-passphrase-123"
	key := crypto.DeriveKey(passphrase)
	plaintext := []byte("DB_PASSWORD=supersecret\nAPI_KEY=abc123")

	encrypted, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if len(encrypted) == 0 {
		t.Fatal("expected non-empty encrypted string")
	}

	decrypted, err := crypto.Decrypt(key, encrypted)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Fatalf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	key := crypto.DeriveKey("correct-passphrase")
	wrongKey := crypto.DeriveKey("wrong-passphrase")
	plaintext := []byte("SECRET=value")

	encrypted, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = crypto.Decrypt(wrongKey, encrypted)
	if err == nil {
		t.Fatal("expected decryption to fail with wrong key")
	}
}

func TestEncryptProducesUniqueOutputs(t *testing.T) {
	key := crypto.DeriveKey("passphrase")
	plaintext := []byte("SAME_INPUT=value")

	enc1, _ := crypto.Encrypt(key, plaintext)
	enc2, _ := crypto.Encrypt(key, plaintext)

	if enc1 == enc2 {
		t.Fatal("expected different ciphertexts due to random nonce")
	}
}
