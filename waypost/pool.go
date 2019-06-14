package waypost

import (
	"math/rand"
	"strings"
	"sync"
)

// Member defines a pool member
// TODO: not currently utilized
type Member struct {
	Value  interface{}
	Weight int
}

// Pool defines a pool of ~Members~ strings
type Pool struct {
	mutex   sync.RWMutex
	members []string
}

// InitPool creates a Pool with the specified seed.
func InitPool(seed []string) Pool {
	return Pool{
		members: seed,
	}
}

// Set adds a member to the Pool.
// A member is unique to the Pool.
func (p *Pool) Set(v string) {
	v = strings.TrimSpace(v)
	if len(v) == 0 {
		return
	}
	p.mutex.Lock()

	for _, m := range p.members {
		if m == v {
			p.mutex.Unlock()
			return
		}
	}

	p.members = append(p.members, v)
	p.mutex.Unlock()
}

// Unset removes a member to the Pool.
// A member is unique to the Pool.
func (p *Pool) Unset(v string) {
	v = strings.TrimSpace(v)

	p.mutex.Lock()
	for idx, m := range p.members {
		if m == v {
			p.members[idx] = p.members[len(p.members)-1]
			p.members = p.members[:len(p.members)-1]
			break
		}
	}
	p.mutex.Unlock()
}

// Select selects a member at random from the Pool.
func (p *Pool) Select() string {
	p.mutex.RLock()
	s := p.members[rand.Intn(len(p.members))]
	p.mutex.RUnlock()
	return s
}

// List returns the pool members.
func (p *Pool) List() []string {
	p.mutex.RLock()
	s := p.members
	p.mutex.RUnlock()

	// TODO: make this sorted by weight
	result := []string{}
	seen := make(map[string]string)
	for _, v := range s {
		if _, exists := seen[v]; !exists {
			seen[v] = v
			result = append(result, v)
		}
	}

	return s
}
