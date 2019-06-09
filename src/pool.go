package main

import (
	"math/rand"
	"sync"
)

type Member struct {
	Value  interface{}
	Weight int
}

type Pool struct {
	mutex   sync.RWMutex
	members []string
}

func InitPool(seed []string) Pool {
	return Pool{
		members: seed,
	}
}

func (p *Pool) Set(v string) {
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

func (p *Pool) Unset(v string) {
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

func (p *Pool) Select() string {
	p.mutex.RLock()
	s := p.members[rand.Intn(len(p.members))]
	p.mutex.RUnlock()
	return s
}

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
