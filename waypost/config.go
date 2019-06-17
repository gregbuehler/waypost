package waypost

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type filterListConfig struct {
	Enabled   bool     `toml:"enabled"`
	List      []string `toml:"list"`
	Upstreams []string `toml:"upstreams"`
}

type serverConfig struct {
	Host     string              `toml:"host"`
	DNSPort  int                 `toml:"dnsPort"`
	MgmtPort int                 `toml:"mgmtPort"`
	Logging  serverLoggingConfig `toml:"logging"`
}

type serverLoggingConfig struct {
	Level    string `toml:"level"`
	Location string `toml:"location"`
}

type resolveConfig struct {
	Nameservers []string `toml:"nameservers"`
	Search      []string `toml:"search"`
}

type cacheConfig struct {
	Enabled bool  `toml:"enabled"`
	TTL     int64 `toml:"ttl"`
}

type hostsConfig struct {
	Enabled   bool              `toml:"enabled"`
	List      map[string]string `toml:"list"`
	Upstreams []string          `toml:"upstreams"`
}

type filteringConfig struct {
	Enabled   bool             `toml:"enabled"`
	Blacklist filterListConfig `toml:"blacklist"`
	Whitelist filterListConfig `toml:"whitelist"`
}

// Config describes the Waypost configuration
type Config struct {
	Server    serverConfig    `toml:"server"`
	Resolve   resolveConfig   `toml:"resolve"`
	Cache     cacheConfig     `toml:"cache"`
	Hosts     hostsConfig     `toml:"hosts"`
	Filtering filteringConfig `toml:"filtering"`
}

// NewConfig parses arguments to create a Config
func NewConfig(configPath string) (*Config, error) {
	return LoadConfigFile(configPath)
}

func LoadConfigFile(path string) (*Config, error) {
	var c Config
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	_, err = toml.Decode(string(data), &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
