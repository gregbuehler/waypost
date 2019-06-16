package waypost

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// Server defines the runtime state of the waypost service
type Server struct {
	mutex sync.RWMutex

	started bool

	Logger    *logrus.Logger
	Config    *Config
	Cache     *Repository
	Hosts     *Repository
	Blacklist *Repository
	Whitelist *Repository
	Resolvers *Pool
	Search    *Repository

	dnsServer  *dns.Server
	mgmtServer *MgmtServer
}

func getEmptyResponse() *dns.Msg {
	m := new(dns.Msg)
	return m
}

func NewServer(config *Config) (*Server, error) {
	s := new(Server)
	s.Config = config
	s.bootstrap()

	return s, nil
}

// InitServer() initializes a waypost service definition
func (s *Server) bootstrap() {
	log.SetOutput(os.Stdout)
	config := s.Config
	logLevel, err := log.ParseLevel(s.Config.Server.Logging.Level)
	if err != nil {
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.SetLevel(logLevel)
	}

	log.Tracef("Initializing Cache")
	c := InitRepository(config.Cache.TTL)
	s.Cache = &c
	log.Debugf("Initialized Cache")

	log.Tracef("Initializing Hosts")
	h := InitRepository(0)
	for k, v := range config.Hosts.List {
		log.Tracef("Adding host map %s %s", k, v)
		h.Set(k, v)
	}
	s.Hosts = &h
	log.Debugf("Initialized Hosts")
	go func() {
		log.Tracef("Loading Hosts from upstreams...")
		for _, n := range config.Hosts.Upstreams {
			log.Debugf("Hosts loading %s", n)
			c, err := h.Fetch(n)
			if err != nil {
				log.Warnf("Hosts failed to load %s. %s", n, err)
			}
			log.Infof("Hosts loaded %d entries from %s", c, n)
		}
		log.Debugf("Loaded Hosts from upstreams")
	}()

	log.Tracef("Initializing Blacklist")
	b := InitRepository(0)
	for _, n := range config.Filtering.Blacklist.List {
		b.Set(n, nil)
	}
	s.Blacklist = &b
	log.Debugf("Initialized Blacklist")
	log.Tracef("Loading Blacklists from upstreams...")
	go func() {
		for _, n := range config.Filtering.Blacklist.Upstreams {
			log.Debugf("Blacklist loading %s", n)
			c, err := b.Fetch(n)
			if err != nil {
				log.Warnf("Blacklist failed to load %s. %s", n, err)
			}
			log.Infof("Blacklist loaded %d entries from %s", c, n)
		}
		log.Debug("Loaded Blacklists from upstreams")
	}()

	log.Tracef("Initializing Whitelist")
	w := InitRepository(0)
	for _, n := range config.Filtering.Whitelist.List {
		w.Set(n, nil)
	}
	s.Whitelist = &w
	log.Debugf("Initialized Whitelist")
	log.Tracef("Loading Whitelists from upstreams...")
	go func() {
		for _, n := range config.Filtering.Whitelist.Upstreams {
			log.Debugf("Whitelist loading %s", n)
			c, err := w.Fetch(n)
			if err != nil {
				log.Warnf("Whitelist failed to load %s. %s", n, err)
			}
			log.Infof("Whitelist loaded %d entries from %s", c, n)
		}
		log.Debug("Loaded Whitelists from upstreams")
	}()

	log.Tracef("Initializing Resolvers")
	r := InitPool([]string{})
	for _, n := range config.Resolve.Nameservers {
		r.Set(fmt.Sprintf("%s:53", n))
	}
	s.Resolvers = &r
	log.Debugf("Initialized Resolvers")

	log.Tracef("Initializing Search Domains")
	d := InitRepository(0)
	for _, n := range config.Resolve.Search {
		d.Set(n, nil)
	}
	s.Search = &d
	log.Debugf("Initialized Search Domains")

	host := config.Server.Host
	dnsPort := config.Server.DNSPort
	mgmtPort := config.Server.MgmtPort

	log.Tracef("Initializing DNS Server")
	s.dnsServer = &dns.Server{
		Addr: fmt.Sprintf("%s:%d", host, dnsPort),
		Net:  "udp",
	}
	log.Debugf("Initialized DNS Server")

	log.Tracef("Initializing Management Server")
	s.mgmtServer = &MgmtServer{
		Addr:    fmt.Sprintf("%s:%d", host, mgmtPort),
		Waypost: s,
	}
	log.Debugf("Initialized Management Server")
}

func (s *Server) printInfo() {
	fmt.Printf("waypost pid=%d host=%s dnsPort=%d mgmtPort=%d",
		os.Getpid(),
		s.Config.Server.Host,
		s.Config.Server.DNSPort,
		s.Config.Server.MgmtPort,
	)
}

// Shutdown gracefully shuts down the Waypost server
func (s *Server) Shutdown() {
	s.dnsServer.Shutdown()
}

// ListenAndServe starts the waypost server
func (s *Server) ListenAndServe() error {
	dns.HandleFunc(".", s.handleRequest)
	go s.dnsServer.ListenAndServe()
	go s.mgmtServer.ListenAndServe()

	sigc := make(chan os.Signal)
	signal.Notify(sigc,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		//syscall.SIGINFO, // TODO: SIGINFO doesn't exist in *nix, just BSD?
		syscall.SIGTERM,
	)

	exit := make(chan int)
	go func() {
		for {
			sigv := <-sigc
			switch sigv {
			case syscall.SIGINT:
				log.Infof("SIGINT received, stopping...")
				exit <- 0
			case syscall.SIGHUP:
				// TODO: Implement SIGHUP
				log.Infof("SIGHUP received, reloading...")
				log.Warnf("Not yet implemented!")
			case syscall.SIGUSR1:
				// TODO: Implement SIGUSR1
				log.Infof("SIGUSR1 received, reloading...")
				log.Warnf("Not yet implemented!")
			case syscall.SIGUSR2:
				// TODO: Implement SIGUSR2
				log.Infof("SIGUSR1 received, reloading...")
				log.Warnf("Not yet implemented!")
			// case syscall.SIGINFO:
			// 	s.printInfo()
			case syscall.SIGTERM:
				log.Infof("SIGTERM received, stopping...")
				s.Shutdown()
				exit <- 0
			default:
				log.Errorf("Signal received, unknown type")
			}
		}
	}()

	code := <-exit
	os.Exit(code)

	return nil
}

func (s *Server) handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	switch r.Opcode {
	case dns.OpcodeQuery:
		m, err := s.resolve(r)
		if err != nil {
			log.Warnf("Problem resolving from upstream: %v", err)
		}

		m.SetReply(r)
		w.WriteMsg(m)
	}
}

func examineName(n string) (string, string, bool) {
	if len(n) == 0 {
		return n, n, false
	}
	if n[len(n)-1] == '.' {
		s := strings.TrimRight(n, ".")
		return n, s, true
	}

	return n, n, false
}

func (s *Server) resolveFromCache(name string) (*dns.Msg, bool, error) {
	c, cached := s.Cache.Get(name)
	if cached {
		log.Tracef("cache hit for %s", name)
		m, ok := c.(*dns.Msg)
		if !ok {
			return nil, false, fmt.Errorf("Failed to reassemble response from cache")
		}
		return m, true, nil
	}

	return nil, false, nil
}

func (s *Server) resolveFromHosts(name string) (*dns.Msg, bool, error) {
	l, exists := s.Hosts.Get(name)
	log.Tracef("action=resolveFromHosts name=%s exists=%t lookup=%s", name, exists, l)
	if exists {
		a := net.ParseIP(l.(string))
		if a != nil {
			t := &dns.A{
				Hdr: dns.RR_Header{
					Name:   fmt.Sprintf("%s.", name),
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    0,
				},
				A: a,
			}
			m := new(dns.Msg)
			m.Compress = false
			m.Answer = append(m.Answer, t)

			return m, true, nil
		}
	}
	return nil, false, nil
}

func (s *Server) resolve(r *dns.Msg) (*dns.Msg, error) {
	if len(r.Question) == 0 {
		return nil, fmt.Errorf("Query with no Question receieved")
	}

	_, sname, authoritative := examineName(r.Question[0].Name)
	qtype := r.Question[0].Qtype
	qslug := fmt.Sprintf("%s-%d", sname, qtype)

	if s.Config.Filtering.Enabled && s.Config.Filtering.Blacklist.Enabled {
		log.Tracef("filtering enabled; checking lists for %s", sname)
		_, _, blacklisted := s.Blacklist.Search(sname)
		if s.Config.Filtering.Blacklist.Enabled == false {
			blacklisted = false
		}

		whitelisted := false
		if s.Config.Filtering.Whitelist.Enabled {
			_, _, whitelisted = s.Whitelist.Search(sname)
		}

		blocked := blacklisted && !whitelisted
		log.Tracef("action=filter sname=%s blacklist=%t whitelist=%t blocked=%t", sname, blacklisted, whitelisted, blocked)
		if blacklisted && !whitelisted {
			log.Infof("action=filter sname=%s blacklist=%t whitelist=%t blocked=%t", sname, blacklisted, whitelisted, blocked)
			m := getEmptyResponse()
			return m, nil
		}
	}

	if s.Config.Cache.Enabled {
		m, hit, err := s.resolveFromCache(qslug)
		if hit && err == nil {
			return m, nil
		}
		if err != nil {
			log.Errorf("error checking cache: %s", err.Error())
		}
		log.Tracef("cache miss for '%s'...", qslug)
	}

	if s.Config.Hosts.Enabled {
		log.Tracef("hosts enabled; checking lists for '%s'", sname)
		m, hit, err := s.resolveFromHosts(sname)
		if hit && err == nil {
			log.Tracef("action=resolveFromHostsResult name=%s hit=%t m=%#v", sname, hit, m)
			return m, nil
		}
		if err != nil {
			log.Errorf("error checking hosts: %s", err.Error())
		}
		log.Tracef("host record missing for '%s'...", sname)
	}

	searchDomains := s.Search.List()
	if authoritative || searchDomains == nil || len(searchDomains) == 0 {
		searchDomains = []string{""}
	}
	log.Tracef("action=lookup sname=%s qtype=%d authoritative=%t, domains=%v", sname, qtype, authoritative, searchDomains)
	for _, domain := range searchDomains {
		sname := sname
		if !authoritative {
			sname = fmt.Sprintf("%s.%s", sname, domain)
		}
		for _, upstream := range s.Resolvers.List() {
			log.Tracef("resolver %s, lookup up '%s'...", upstream, sname)
			m, err := dns.Exchange(r, upstream)
			if err != nil {
				log.Tracef("resolver %s, lookup '%s' error, trying next", upstream, sname)
				log.Print(err)
				continue
			}

			if m == nil || m.Answer == nil {
				log.Tracef("resolver %s, lookup '%s' failed, trying next", upstream, sname)
				continue
			}

			log.Debugf("caching '%s': %s", qslug, m.Answer[0].String())
			s.Cache.Set(qslug, m)
			return m, nil
		}
	}

	return getEmptyResponse(), nil
}
