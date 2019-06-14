package waypost

import (
	"bytes"
	"fmt"
	"net"
	"strings"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

type MgmtServer struct {
	Addr    string
	Waypost *Server
}

func (m *MgmtServer) ListenAndServe() error {
	log.Debugf("Starting managment interface...")
	l, err := net.Listen("tcp", m.Addr)
	if err != nil {
		return err
	}
	log.Infof("Started managment interface...")

	defer l.Close()
	for {
		err := func() error {
			conn, err := l.Accept()
			if err != nil {
				log.Fatal(err)
			}

			go m.handleMgmtRequest(conn)

			return nil
		}()
		if err != nil {
			return err
		}
	}
}

func (m *MgmtServer) handleMgmtRequest(conn net.Conn) {
	b := make([]byte, 1024)
	_, err := conn.Read(b)
	if err != nil {
		log.Error(err)
		conn.Close()
		return
	}
	t := strings.TrimSpace(string(bytes.Trim(b, "\x00")))
	r, err := m.parseAndExecuteCmd(strings.Fields(t))
	if err != nil {
		log.Error(err)
		conn.Close()
		return
	}

	conn.Write([]byte(r))
	conn.Close()
}

func parseCommand(s []string) (string, string, []string) {
	var namespace string
	var action string
	var values []string

	if len(s) == 1 {
		namespace = s[0]
	}
	if len(s) >= 2 {
		namespace = s[0]
		action = s[1]
		for _, v := range s[2:] {
			v = strings.TrimSpace(v)
			if v != "" {
				values = append(values, v)
			}
		}
	}

	return namespace, action, values
}

func (m *MgmtServer) parseAndExecuteCmd(args []string) (string, error) {
	namespace, action, values := parseCommand(args)

	// TODO: refactor mgmt command map. nested switches are gross.
	switch namespace {
	case "cache":
		switch action {
		case "flush":
			for _, k := range m.Waypost.Cache.List() {
				m.Waypost.Cache.Unset(k)
			}
			return "", nil
		case "list":
			return strings.Join(m.Waypost.Cache.List(), "\n"), nil
		case "enable":
			m.Waypost.Config.Cache.Enabled = true
			return "", nil
		case "disable":
			m.Waypost.Config.Cache.Enabled = false
			return "", nil
		default:
			return "", fmt.Errorf("unknown cache action '%s'", action)
		}
	case "nameservers":
		switch action {
		case "add":
			for _, v := range values {
				log.Infof("Adding resolver %s", v)
				m.Waypost.Resolvers.Set(v)
			}
			return "", nil
		case "del":
			for _, v := range values {
				log.Infof("Removing resolver %s", v)
				m.Waypost.Resolvers.Unset(v)
			}
			return "", nil
		case "list":
			return strings.Join(m.Waypost.Resolvers.List(), "\n"), nil
		default:
			return "", fmt.Errorf("unknown nameservers action '%s'", action)
		}
	case "search":
		switch action {
		case "add":
			for _, v := range values {
				log.Infof("Adding search domain %s", v)
				m.Waypost.Search.Set(v, nil)
			}
			return "", nil
		case "del":
			for _, v := range values {
				log.Infof("Removing search domain %s", v)
				m.Waypost.Search.Unset(v)
			}
			return "", nil
		case "list":
			return strings.Join(m.Waypost.Search.List(), "\n"), nil
		default:
			return "", fmt.Errorf("unknown search domain action '%s'", action)
		}
	case "blacklist":
		switch action {
		case "add":
			for _, v := range values {
				log.Infof("Adding blacklist %s", v)
				m.Waypost.Blacklist.Set(v, nil)
				m.Waypost.Cache.Unset(v)
			}
			return "", nil
		case "del":
			for _, v := range values {
				log.Infof("Removing blacklist %s", v)
				m.Waypost.Blacklist.Unset(v)
				m.Waypost.Cache.Unset(v)
			}
			return "", nil
		case "list":
			return strings.Join(m.Waypost.Blacklist.List(), "\n"), nil
		case "enable":
			m.Waypost.Config.Filtering.Blacklist.Enabled = true
			return "", nil
		case "disable":
			m.Waypost.Config.Filtering.Blacklist.Enabled = false
			return "", nil
		default:
			return "", fmt.Errorf("unknown blacklist action '%s'", action)
		}
	case "whitelist":
		switch action {
		case "add":
			for _, v := range values {
				log.Infof("Adding whitelist %s", v)
				m.Waypost.Whitelist.Set(v, nil)
				m.Waypost.Cache.Unset(v)
			}
			return "", nil
		case "del":
			for _, v := range values {
				log.Infof("Removing whitelist %s", v)
				m.Waypost.Whitelist.Unset(v)
				m.Waypost.Cache.Unset(v)
			}
			return "", nil
		case "list":
			return strings.Join(m.Waypost.Whitelist.List(), "\n"), nil
		case "enable":
			m.Waypost.Config.Filtering.Whitelist.Enabled = true
			return "", nil
		case "disable":
			m.Waypost.Config.Filtering.Whitelist.Enabled = false
			return "", nil
		default:
			return "", fmt.Errorf("unknown whitelist action '%s'", action)
		}
	case "hosts":
		switch action {
		case "add":
			if len(values) != 2 {
				return "", fmt.Errorf("hosts add arguments error")
			}

			log.Infof("Adding hosts mapping %s to %s", values[0], values[1])
			m.Waypost.Hosts.Set(values[0], values[1])
			m.Waypost.Cache.Unset(values[0])

			return "", nil
		case "del":
			if len(values) != 2 {
				return "", fmt.Errorf("hosts add arguments error")
			}

			log.Infof("Removing hosts mapping %s to %s", values[0], values[1])
			m.Waypost.Hosts.Set(values[0], values[1])
			m.Waypost.Cache.Unset(values[0])

			return "", nil
		case "list":
			l := make([]string, 0)
			for _, k := range m.Waypost.Hosts.List() {
				v, _ := m.Waypost.Hosts.Get(k)
				l = append(l, fmt.Sprintf("%s\t%s", k, v))
			}
			return strings.Join(l, "\n"), nil
		case "enable":
			m.Waypost.Config.Hosts.Enabled = true
			return "", nil
		case "disable":
			m.Waypost.Config.Hosts.Enabled = false
			return "", nil
		default:
			return "", fmt.Errorf("unknown hosts action '%s'", action)
		}
	case "config":
		switch action {
		case "reload":
			return "", fmt.Errorf("config action 'reload' is not yet implemented")
		case "dump":
			hostMap := make(map[string]string, 0)
			for _, k := range m.Waypost.Hosts.List() {
				v, exists := m.Waypost.Hosts.Get(k)
				if exists {
					hostMap[k] = v.(string)
				}
			}
			m.Waypost.Config.Hosts.List = hostMap

			buffer := new(bytes.Buffer)
			encoder := toml.NewEncoder(buffer)
			err := encoder.Encode(m.Waypost.Config)
			if err != nil {
				return "", fmt.Errorf("failed to dump config: %s", err.Error())
			}
			return buffer.String(), nil
		default:
			return "", fmt.Errorf("unknown hosts action '%s'", action)
		}
	default:
		return "", fmt.Errorf("unknown namespace '%s'", namespace)
	}
	return "", fmt.Errorf("failed to handle managment command")
}
