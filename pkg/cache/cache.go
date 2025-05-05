package cache

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// DefaultCachePath defines a default path for storing file cache entries.
const DefaultCachePath = "./cache"

// ErrNoCache defines an error that may occur if the requested cache entry was
// not found.
var ErrNoCache = errors.New("cache entry not found")

// Option defines a type for providing configuration options for Cache instance.
type Option func(*Cache)

// WithCacheDir allows setting a custom path for storing cache files. By default,
// DefaultCachePath will be used.
func WithCacheDir(dir string) Option {
	return func(c *Cache) {
		c.cachePath = dir
	}
}

// WithCompress enables or disables compression of cache files. By default it's
// enabled.
func WithCompress(enable bool) Option {
	return func(c *Cache) {
		c.compressionEnabled = enable
	}
}

// Cache provides a file cache API. It saves and restores an arbitrary data types
// which can marshaled/unmarshalled as JSON into file cache. The cache entries
// can be restored by the corresponding key. The cache files will be stored at
// configured path or default DefaultCachePath will be used.
type Cache struct {
	cachePath          string
	compressionEnabled bool
}

func NewCache(opts ...Option) *Cache {
	c := &Cache{
		cachePath:          DefaultCachePath,
		compressionEnabled: true,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Get retrieves a cache entry by its key and maps data to the provided target.
// The target instance must be a pointer type and support JSON unmarshalling.
// If the corresponding cache entry was not found the ErrNoCache error will be
// returned.
func (c *Cache) Get(key string, target any) error {
	cacheFile, err := os.Open(filepath.Join(c.cachePath, c.keyHash(key)))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNoCache
		}

		return err
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(cacheFile)

	rawCache, err := c.decompress(cacheFile)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(rawCache, target); err != nil {
		return err
	}

	return nil
}

// Set creates a new cache entry for specified key or overrides the existing one.
// The val instance must support JSON marshalling.
func (c *Cache) Set(key string, val any) error {
	cacheFile, err := c.initEntryCache(key)
	if err != nil {
		return err
	}

	defer func(f io.WriteCloser) {
		_ = f.Close()
	}(cacheFile)

	rawData, err := json.Marshal(val)
	if err != nil {
		return err
	}

	if err = c.compress(cacheFile, rawData); err != nil {
		return err
	}

	return nil
}

func (c *Cache) compress(w io.WriteCloser, data []byte) error {
	if c.compressionEnabled {
		w = gzip.NewWriter(w)
	}

	defer func(wc io.WriteCloser) {
		_ = wc.Close()
	}(w)

	if _, err := w.Write(data); err != nil {
		return err
	}

	return nil
}

func (c *Cache) decompress(r io.ReadCloser) ([]byte, error) {
	var err error

	if c.compressionEnabled {
		if r, err = gzip.NewReader(r); err != nil {
			return nil, err
		}
	}

	defer func(rc io.ReadCloser) {
		_ = rc.Close()
	}(r)

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Cache) initEntryCache(key string) (io.WriteCloser, error) {
	cachePath := filepath.Join(c.cachePath, c.keyHash(key))

	cacheFile, err := os.Create(cachePath)
	if err != nil {
		return nil, err
	}

	return cacheFile, nil
}

func (c *Cache) keyHash(key string) string {
	if len(key) == 0 {
		return ""
	}

	h := sha256.New()
	h.Write([]byte(key))

	return hex.EncodeToString(h.Sum(nil))
}
