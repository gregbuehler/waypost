package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/miekg/dns"
)

type Waypost struct {
	Cache     *Repository
	Blacklist *Repository
	Whitelist *Repository
	Upstreams *Pool

	BlockingEnabled bool
}

func getEmptyResponse() *dns.Msg {
	m := new(dns.Msg)

	return m
}

func InitWaypost() Waypost {
	c := InitRepository(5)
	b := InitRepository(0)
	w := InitRepository(0)
	u := InitPool([]string{})

	config, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	warnIfErr(err)
	for s := range config.Servers {
		n := config.Servers[s]
		u.Set(fmt.Sprintf("%s:53", n))
	}

	return Waypost{
		Cache:     &c,
		Blacklist: &b,
		Whitelist: &w,
		Upstreams: &u,
	}
}

func (w *Waypost) Resolve(r *dns.Msg) (*dns.Msg, error) {
	if len(r.Question) <= 0 {
		return getEmptyResponse(), nil
	}
	var m *dns.Msg

	q := r.Question[0]
	c, cached := w.Cache.Get(q.Name)
	if cached {
		log.Printf("Serving %s from cache", q.Name)
		c, ok := c.(*dns.Msg)
		m = c
		if !ok {
			return getEmptyResponse(), errors.New("Failed to reassemble response from cache")
		}
	} else {

		if w.BlockingEnabled {
			_, blacklisted := w.Blacklist.Get(q.Name)
			_, whitelisted := w.Whitelist.Get(q.Name)

			if blacklisted && !whitelisted {
				m := getEmptyResponse()
				return m, nil
			}
		}

		upstream := w.Upstreams.Select()
		c, err := dns.Exchange(r, upstream)
		m = c
		warnIfErr(err)

		w.Cache.Set(q.Name, m)
	}

	return m, nil
}
