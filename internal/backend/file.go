package backend

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// FileBackend resolves secrets from a dotenv-style file.
type FileBackend struct {
	path   string
	cache  map[string]string
	loaded bool
}

// NewFileBackend creates a FileBackend that reads from the given path.
func NewFileBackend(path string) *FileBackend {
	return &FileBackend{path: path}
}

// load parses the dotenv file into the in-memory cache (once).
func (f *FileBackend) load() error {
	if f.loaded {
		return nil
	}
	file, err := os.Open(f.path)
	if err != nil {
		return fmt.Errorf("file backend: open %q: %w", f.path, err)
	}
	defer file.Close()

	f.cache = make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		f.cache[key] = val
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("file backend: scan %q: %w", f.path, err)
	}
	f.loaded = true
	return nil
}

// Get returns the value for key from the dotenv file.
func (f *FileBackend) Get(key string) (string, error) {
	if err := f.load(); err != nil {
		return "", err
	}
	val, ok := f.cache[key]
	if !ok {
		return "", ErrSecretNotFound{Key: key}
	}
	return val, nil
}

// List returns all keys present in the dotenv file.
func (f *FileBackend) List() ([]string, error) {
	if err := f.load(); err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(f.cache))
	for k := range f.cache {
		keys = append(keys, k)
	}
	return keys, nil
}

// String returns a human-readable description of the backend.
func (f *FileBackend) String() string {
	return fmt.Sprintf("file(%s)", f.path)
}
