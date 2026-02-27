package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type entry struct {
	data      any
	expiresAt time.Time
}

// Cache provides in-memory TTL caching with optional disk persistence.
type Cache struct {
	mu      sync.RWMutex
	items   map[string]entry
	diskDir string
}

// New creates a cache. If diskDir is non-empty, disk persistence is enabled.
func New(diskDir string) *Cache {
	if diskDir != "" {
		os.MkdirAll(diskDir, 0o755)
	}
	return &Cache{
		items:   make(map[string]entry),
		diskDir: diskDir,
	}
}

// Get retrieves a value from the cache. Returns nil if expired or missing.
func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.items[key]
	if !ok || time.Now().After(e.expiresAt) {
		return nil, false
	}
	return e.data, true
}

// Set stores a value with a TTL.
func (c *Cache) Set(key string, data any, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = entry{data: data, expiresAt: time.Now().Add(ttl)}
}

// Clear removes a specific key.
func (c *Cache) Clear(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// ClearAll removes all cached entries.
func (c *Cache) ClearAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]entry)
}

// Age returns how long ago a key was last set (approximated from TTL).
func (c *Cache) Age(key string, originalTTL time.Duration) time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.items[key]
	if !ok {
		return 0
	}
	remaining := time.Until(e.expiresAt)
	if remaining < 0 {
		return originalTTL
	}
	return originalTTL - remaining
}

// diskPath returns the file path for a disk-cached key.
func (c *Cache) diskPath(key string) string {
	return filepath.Join(c.diskDir, key+".json")
}

// diskEntry is the on-disk format.
type diskEntry struct {
	Data      json.RawMessage `json:"data"`
	ExpiresAt time.Time       `json:"expires_at"`
}

// SaveToDisk persists a value to disk with a TTL.
func (c *Cache) SaveToDisk(key string, data any, ttl time.Duration) error {
	if c.diskDir == "" {
		return nil
	}
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	de := diskEntry{Data: raw, ExpiresAt: time.Now().Add(ttl)}
	b, err := json.Marshal(de)
	if err != nil {
		return err
	}
	return os.WriteFile(c.diskPath(key), b, 0o644)
}

// LoadFromDisk loads a value from disk into target. Returns false if expired or missing.
func (c *Cache) LoadFromDisk(key string, target any) bool {
	if c.diskDir == "" {
		return false
	}
	b, err := os.ReadFile(c.diskPath(key))
	if err != nil {
		return false
	}
	var de diskEntry
	if err := json.Unmarshal(b, &de); err != nil {
		return false
	}
	if time.Now().After(de.ExpiresAt) {
		os.Remove(c.diskPath(key))
		return false
	}
	return json.Unmarshal(de.Data, target) == nil
}
