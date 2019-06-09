package main

import (
	"flag"
	"log"
)

type Config struct {
	CacheExpiration   *int
	Debug             *bool
	DynamicReloading  *bool
	BlacklistLocation *string
	WhitelistLocation *string
}

func InitConfig() (Config, error) {
	fileName := flag.String("config", "waypost.json", "configuration file")
	cacheExpiration := flag.Int("expiration", 5, "cache expiration (in seconds)")
	debugLogging := flag.Bool("debug", false, "enable debug logging")
	dynamicReloading := flag.Bool("dynamic", false, "enable dynamic configuration reloading")
	blacklistLocation := flag.String("blacklists", "./blacklists.d", "directory containing blacklists")
	whitelistLocation := flag.String("whitelists", "./whitelists.d", "directory containing whitelists")
	flag.Parse()

	// TODO: load from fileName
	log.Printf("TODO: load from %s", fileName)

	return Config{
		CacheExpiration:   cacheExpiration,
		Debug:             debugLogging,
		DynamicReloading:  dynamicReloading,
		BlacklistLocation: blacklistLocation,
		WhitelistLocation: whitelistLocation,
	}, nil
}
