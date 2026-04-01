// encrypt — encrypts all .json content files to .enc
// Usage: go run cmd/encrypt/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"aegis-inferno/internal/crypto"
)

func main() {
	key := crypto.DeriveKey()
	contentDir := "content"

	err := filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		plain, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		enc, err := crypto.Encrypt(plain, key)
		if err != nil {
			return fmt.Errorf("encrypt %s: %w", path, err)
		}

		outPath := strings.TrimSuffix(path, ".json") + ".enc"
		if err := os.WriteFile(outPath, enc, 0644); err != nil {
			return fmt.Errorf("write %s: %w", outPath, err)
		}

		fmt.Printf("Encrypted: %s -> %s (%d -> %d bytes)\n",
			path, outPath, len(plain), len(enc))
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Done. All content encrypted.")
}
