// Copyright 2024 Flant JSC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/deckhouse/deckhouse/pkg/log"
	"github.com/pkg/errors"

	"github.com/deckhouse/deckhouse/go_lib/registry-packages-proxy/cache"
)

const HighUsagePercent = 80

type CacheEntry struct {
	lastAccessTime time.Time
	size           uint64
}

type Cache struct {
	storage map[string]*CacheEntry
	sync.RWMutex
	logger        *log.Logger
	root          string
	retentionSize uint64

	metrics *Metrics
}

func NewCache(logger *log.Logger, root string, retentionSize uint64, metrics *Metrics) *Cache {
	storage := make(map[string]*CacheEntry)
	return &Cache{
		storage:       storage,
		logger:        logger,
		root:          root,
		retentionSize: retentionSize,
		metrics:       metrics,
	}
}

func (c *Cache) Get(digest string) (int64, io.ReadCloser, error) {
	path := c.digestToPath(digest)

	// check if cache entry exists
	if !c.storageGetOK(digest) {
		c.logger.Info("entry with digest is not found in the cache", slog.String("digest", digest))
		return 0, nil, cache.ErrEntryNotFound
	}

	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, nil, cache.ErrEntryNotFound
		}

		return 0, nil, err
	}

	stat, err := file.Stat()
	c.logger.Info("found file with size in the cache", slog.String("name", stat.Name()), slog.Int64("size", stat.Size()))

	if err != nil {
		return 0, nil, err
	}

	c.Lock()
	c.storage[digest].lastAccessTime = time.Now()
	c.Unlock()

	return stat.Size(), file, nil
}

func (c *Cache) Set(digest string, size int64, reader io.Reader) error {
	// check if cache entry exists
	if c.storageGetOK(digest) {
		c.logger.Info("entry with digest already exists, skipping", slog.String("digest", digest))
		return nil
	}

	c.logger.Info("write file with digest with size to the cache dir", slog.String("digest", digest), slog.Int64("size", size))

	path := c.digestToPath(digest)

	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}

	c.Lock()
	c.storage[digest] = &CacheEntry{
		lastAccessTime: time.Now(),
		size:           uint64(size),
	}
	c.Unlock()

	c.metrics.CacheSize.Add(float64(size))
	return nil
}

func (c *Cache) Delete(digest string) error {
	// check if cache entry exists
	if !c.storageGetOK(digest) {
		c.logger.Info("entry with digest doesn't exists, skipping", slog.String("digest", digest))
		return nil
	}

	path := c.digestToPath(digest)
	c.logger.Info("remove file with path from the cache dir", slog.String("path", path))

	err := os.Remove(path)
	if err != nil {
		return err
	}

	c.Lock()
	c.metrics.CacheSize.Sub(float64(c.storage[digest].size))
	delete(c.storage, digest)
	c.Unlock()

	return nil
}

func (c *Cache) Reconcile(ctx context.Context) {
	c.logger.Info("starting cache reconcile loop")

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := c.ApplyRetentionPolicy()
			if err != nil {
				c.logger.Error("reconcile loop failed", log.Err(err))
				return
			}
		}
	}
}

func (c *Cache) ApplyRetentionPolicy() error {
	for {
		usagePercent := int(float64(c.calculateCacheSize()) / float64(c.retentionSize) * 100)
		if usagePercent < HighUsagePercent {
			c.logger.Info(fmt.Sprintf("current cache usage less than %d%%, compaction is not needed", HighUsagePercent), slog.Int("usage", usagePercent))
			return nil
		}

		if len(c.storage) == 0 {
			c.logger.Info("storage map is empty")
			return nil
		}

		c.logger.Info(fmt.Sprintf("need to compact cache, current usage more than %d%%", HighUsagePercent), slog.Int("usage", usagePercent))

		// sort descending by last access time
		var oldestDigest string
		lowestTime := time.Now()

		c.Lock()
		for k, v := range c.storage {
			if lowestTime.Compare(v.lastAccessTime) >= 0 {
				oldestDigest = k
			}
		}
		c.Unlock()

		// remove oldest entry
		err := c.Delete(oldestDigest)
		if err != nil {
			return err
		}
	}
}

func (c *Cache) calculateCacheSize() uint64 {
	c.Lock()
	defer c.Unlock()
	var size uint64
	for _, v := range c.storage {
		size += v.size
	}
	return size
}

func (c *Cache) digestToPath(digest string) string {
	// digest format is sha256:1234567....
	// remove sha256: and convert to path
	hash := digest[7:]
	return filepath.Join(c.root, "packages", hash[:2], hash)
}

func (c *Cache) storageGetOK(digest string) bool {
	c.Lock()
	defer c.Unlock()
	_, ok := c.storage[digest]
	return ok
}
