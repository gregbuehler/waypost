package waypost

import (
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Entry describes a value and timestamp within a Repository
type Entry struct {
	Value   interface{}
	Created int64
}

// Repository describes a cache-like datastructure
type Repository struct {
	entries map[string]Entry
	mutex   sync.RWMutex
	ttl     int64
}

// InitRepository intializes a repository with entries and TTL
func InitRepository(ttl int64) Repository {
	return Repository{
		entries: make(map[string]Entry),
		ttl:     ttl,
	}
}

func (r *Repository) Search(t string) (string, interface{}, bool) {
	l := r.List()
	for _, k := range l {
		if strings.Contains(k, t) {
			v, _ := r.Get(k)
			return k, v, true
		}
	}

	return "", nil, false
}

// Get returns the value of a Repository entry
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

// Set specifies the value of an entry value
func (r *Repository) Set(k string, v interface{}) {
	k = strings.TrimSpace(k)
	if len(k) == 0 {
		return
	}

	r.mutex.Lock()
	r.entries[k] = Entry{
		Value:   v,
		Created: time.Now().Unix(),
	}

	r.mutex.Unlock()
}

// List keys
func (r *Repository) List() []string {
	l := make([]string, 0, len(r.entries))
	r.mutex.RLock()
	for k := range r.entries {
		l = append(l, k)
	}
	r.mutex.RUnlock()

	return l
}

// Unset removes an item from the Repository
func (r *Repository) Unset(k string) {
	r.mutex.Lock()
	if _, exists := r.entries[k]; exists {
		delete(r.entries, k)
	}
	r.mutex.Unlock()
}

func (r *Repository) Fetch(url string) (int, error) {
	// TODO: handle getting/setting a local cache of remote sources
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()

	data := strings.Split(string(body), "\n")
	for _, v := range data {
		fields := strings.Fields(v)
		if len(fields) == 2 {
			r.Set(fields[0], fields[1])
		}
	}
	return len(data), nil
}
