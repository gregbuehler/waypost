package cmd

import (
	"log"

	"github.com/gregbuehler/waypost/waypost"
	"github.com/spf13/cobra"
)

var (
	serveCmd = &cobra.Command{
		Use:     "serve",
		Short:   "waypost serve",
		Aliases: []string{"s"},
		Run:     serve,
	}
)

func init() {

}

func serve(cmd *cobra.Command, args []string) {
	config, err := waypost.NewConfig(configLocation)
	if err != nil {
		log.Fatal(err)
	}

	if host != defaultHost {
		config.Server.Host = host
	}

	if dnsPort != defaultPort {
		config.Server.DNSPort = dnsPort
	}

	if mgmtPort != defaultMgmtPort {
		config.Server.MgmtPort = mgmtPort
	}

	if logLevel != defaultLogLevel {
		config.Server.Logging.Level = logLevel
	}

	server, err := waypost.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.ListenAndServe())
}
