package cache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Data   interface{}
	Expiry int64
}

func (ce CacheEntry) IsExpired() bool {
	return time.Now().Unix() > ce.Expiry
}

type InMemoryStorage struct {
	entries map[string]CacheEntry
	lock    *sync.RWMutex
}

func (ims InMemoryStorage) Get(key string) interface{} {
	ims.lock.RLock()
	defer ims.lock.RUnlock()

	if entry, ok := ims.entries[key]; !ok {
		return nil
	} else if entry.IsExpired() {
		delete(ims.entries, key)
		return nil
	} else {
		return entry.Data
	}
}

func (ims InMemoryStorage) Set(key string, data interface{}, cacheTime time.Duration) {
	ims.lock.Lock()
	defer ims.lock.Unlock()

	ims.entries[key] = CacheEntry{
		Data:   data,
		Expiry: time.Now().Add(cacheTime).Unix(),
	}
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		entries: make(map[string]CacheEntry),
		lock:    &sync.RWMutex{},
	}
}
