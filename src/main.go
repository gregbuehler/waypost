package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/miekg/dns"
)

var (
	resolverPort   int
	managementPort int

	whitelist []string
	blacklist []string

	whitelistEnabled bool
	blacklistEnabled bool
)

var (
	waypost Waypost
)

func handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	switch r.Opcode {
	case dns.OpcodeQuery:
		m, err := waypost.Resolve(r)
		warnIfErr(err)

		m.SetReply(r)
		w.WriteMsg(m)
		return
	}
}

func main() {
	waypost = InitWaypost()
	waypost.BlockingEnabled = true
	// TODO: do proper list loading
	waypost.Blacklist.Set("doubleclick.net.", nil)

	host := ""
	port := 5553
	bind := fmt.Sprintf("%s:%d", host, port)

	server := &dns.Server{
		Addr: bind,
		Net:  "udp",
	}

	go server.ListenAndServe()
	dns.HandleFunc(".", handleRequest)

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

forever:
	for {
		select {
		case <-sig:
			// TODO: SIGUSR1=reload, SIGHUP=restart
			log.Print("Signal received, stopping...")
			break forever
		}
	}
}
