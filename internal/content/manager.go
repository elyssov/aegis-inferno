package content

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	"aegis-inferno/internal/crypto"
)

// Manager handles encrypted game content.
// Content files are embedded in the binary as .enc files.
// Plain .json files are also supported (dev mode).
type Manager struct {
	fs    fs.FS
	key   []byte
	cache map[string][]byte
	mu    sync.RWMutex
}

// NewManager creates a content manager from an embedded FS.
func NewManager(contentFS fs.FS, key []byte) (*Manager, error) {
	return &Manager{
		fs:    contentFS,
		key:   key,
		cache: make(map[string][]byte),
	}, nil
}

// Get retrieves decrypted content by name (e.g., "ru/chat_her").
// Looks for .enc first (encrypted), falls back to .json (dev mode).
func (m *Manager) Get(name string) ([]byte, error) {
	m.mu.RLock()
	if data, ok := m.cache[name]; ok {
		m.mu.RUnlock()
		return data, nil
	}
	m.mu.RUnlock()

	// Try encrypted first
	encPath := filepath.Join("content", name+".enc")
	data, err := fs.ReadFile(m.fs, encPath)
	if err == nil {
		plain, err := crypto.Decrypt(data, m.key)
		if err != nil {
			return nil, fmt.Errorf("decrypt %s: %w", name, err)
		}
		m.mu.Lock()
		m.cache[name] = plain
		m.mu.Unlock()
		return plain, nil
	}

	// Fallback: plain JSON (dev mode)
	jsonPath := filepath.Join("content", name+".json")
	data, err = fs.ReadFile(m.fs, jsonPath)
	if err != nil {
		return nil, fmt.Errorf("content %s: not found", name)
	}
	m.mu.Lock()
	m.cache[name] = data
	m.mu.Unlock()
	return data, nil
}

// List returns all available content names.
func (m *Manager) List() ([]string, error) {
	var names []string
	err := fs.WalkDir(m.fs, "content", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		name := strings.TrimPrefix(path, "content/")
		name = strings.TrimSuffix(name, ".enc")
		name = strings.TrimSuffix(name, ".json")
		names = append(names, name)
		return nil
	})
	return names, err
}
