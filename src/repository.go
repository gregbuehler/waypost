package main

import (
	"sync"
	"time"
)

type Entry struct {
	Value   interface{}
	Created int64
}

type Repository struct {
	entries map[string]Entry
	mutex   sync.RWMutex
	ttl     int64
}

func InitRepository(ttl int64) Repository {
	return Repository{
		entries: make(map[string]Entry),
		ttl:     ttl,
	}
}

func (r *Repository) Get(k string) (interface{}, bool) {
	r.mutex.RLock()

	v, exists := r.entries[k]
	if !exists {
		r.mutex.RUnlock()
		return nil, false
	}

	if r.ttl > 0 {
		now := time.Now().Unix()
		if now-r.ttl > v.Created {
			r.mutex.RUnlock()
			return nil, false
		}
	}

	r.mutex.RUnlock()
	return v.Value, true
}

func (r *Repository) Set(k string, v interface{}) {
	r.mutex.Lock()

	r.entries[k] = Entry{
		Value:   v,
		Created: time.Now().Unix(),
	}

	r.mutex.Unlock()
}
