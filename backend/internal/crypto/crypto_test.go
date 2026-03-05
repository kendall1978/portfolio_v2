package crypto_test

import (
	"testing"

	"github.com/kendall1978/project_salary/backend/internal/crypto"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	testKey := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	if err := crypto.SetKey(testKey); err != nil {
		t.Fatalf("SetKey failed: %v", err)
	}

	original := "JBSWY3DPEHPK3PXP"
	encrypted, err := crypto.Encrypt(original)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if encrypted == original {
		t.Error("Encrypted value should differ from original")
	}

	decrypted, err := crypto.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if decrypted != original {
		t.Errorf("Expected %q, got %q", original, decrypted)
	}
}

func TestSetKey_InvalidLength(t *testing.T) {
	err := crypto.SetKey("tooshort")
	if err == nil {
		t.Error("Expected error for short key")
	}
}

func TestEncrypt_DifferentCiphertexts(t *testing.T) {
	testKey := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	crypto.SetKey(testKey)

	a, _ := crypto.Encrypt("same-secret")
	b, _ := crypto.Encrypt("same-secret")
	if a == b {
		t.Error("Two encryptions of the same value should produce different ciphertext (random nonce)")
	}
}
