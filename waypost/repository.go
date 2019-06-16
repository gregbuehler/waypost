package waypost

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
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
	for _, i := range l {
		k := strings.Split(i, ";")[0]
		if k == t {
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

func (r *Repository) loadFromString(data string) (int, error) {
	commentPattern := regexp.MustCompile(`(#.*)`)
	lines := strings.Split(data, "\n")
	count := 0

	for _, line := range lines {
		if commentPattern.MatchString(line) {
			line = commentPattern.ReplaceAllString(line, "")
		}

		fields := strings.Fields(line)
		if len(fields) > 0 {
			if len(fields) == 1 {
				r.Set(fields[0], nil)
				count = count + 1
			}
			if len(fields) >= 2 {
				r.Set(fields[1], fields[0])
				count = count + 1
			}
		}
	}

	return count, nil
}

func (r *Repository) Populate(static map[string]string, remotes []string) (int, int, int, error) {
	for k, v := range static {
		r.Set(k, v)
	}

	re := 0
	for _, n := range remotes {
		c, err := r.Fetch(n)
		if err != nil {
			log.Warnf("Failed to load entries from %s: %s", n, err.Error())
		}
		re = re + c
		log.Debugf("Loaded %d entries from %s", c, n)
	}

	return len(static), re, len(r.List()), nil
}

// Fetch populates a repository from a url
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
	l, err := r.loadFromString(string(body))
	return l, nil
}
