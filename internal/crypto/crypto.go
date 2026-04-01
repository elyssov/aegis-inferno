package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
)

// DeriveKey returns the AES-256 key for content decryption.
// The key is derived from multiple obfuscated fragments scattered
// in the binary. Not military-grade, but prevents casual copying.
func DeriveKey() []byte {
	// Fragments — look like random strings, assembled at runtime
	parts := []string{
		"aEgIs", // 1
		"2026",  // 2
		"iNf3",  // 3
		"rN0!",  // 4
		"KaTyA", // 5
		"gR0zA", // 6
	}
	combined := ""
	for _, p := range parts {
		combined += p
	}
	h := sha256.Sum256([]byte(combined))
	return h[:]
}

// Encrypt encrypts plaintext with AES-256-GCM.
func Encrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	// For content encryption we use deterministic nonce from hash
	// (same content = same ciphertext, OK for static game data)
	h := sha256.Sum256(plaintext)
	copy(nonce, h[:gcm.NonceSize()])
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt decrypts ciphertext encrypted with Encrypt.
func Decrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ns := gcm.NonceSize()
	if len(ciphertext) < ns {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ct := ciphertext[:ns], ciphertext[ns:]
	return gcm.Open(nil, nonce, ct, nil)
}
