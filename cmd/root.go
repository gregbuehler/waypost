package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
)

var (
	mgmtClient     *MgmtClient
	host           string
	dnsPort        int
	mgmtPort       int
	configLocation string
	logLevel       string

	helpFlag bool
)

const (
	defaultConfigLocation = "/etc/waypost/waypost.toml"
	defaultHost           = "127.0.0.1"
	defaultPort           = 53
	defaultMgmtPort       = 5554
	defaultLogLevel       = "info"
)

type MgmtClient struct {
	Addr string
}

func (m *MgmtClient) IssueCommand(cmd string) (string, error) {
	conn, err := net.Dial("tcp", m.Addr)
	if err != nil {
		return "", err
	}

	_, err = fmt.Fprint(conn, fmt.Sprintf("%s\n", cmd))
	if err != nil {
		return "", err
	}

	var b []byte
	b, err = ioutil.ReadAll(conn)
	return string(b), err
}

var RootCmd = &cobra.Command{
	Use:   "waypost",
	Short: "waypost",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		mgmtAddr := fmt.Sprintf("%s:%d", host, mgmtPort)
		mgmtClient = &MgmtClient{Addr: mgmtAddr}
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func init() {
	// reassign help to free up `-h` for `--host`
	RootCmd.PersistentFlags().BoolVarP(&helpFlag, "help", "", false, "Help for Waypost")

	// set version on root cmd to cheat cobra version into play
	RootCmd.Version = "0.0.6"

	RootCmd.PersistentFlags().StringVarP(&configLocation, "config", "c", defaultConfigLocation, "configuration file")
	RootCmd.PersistentFlags().StringVarP(&host, "host", "h", defaultHost, "host to bind")
	RootCmd.PersistentFlags().IntVarP(&dnsPort, "port", "p", defaultPort, "dns port")
	RootCmd.PersistentFlags().IntVarP(&mgmtPort, "mport", "m", defaultMgmtPort, "management port")
	RootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", defaultLogLevel, "log level [ERROR|WARN|INFO|DEBUG|TRACE]")

	//RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(serveCmd)
	RootCmd.AddCommand(configCmd)
	RootCmd.AddCommand(cacheCmd)
	RootCmd.AddCommand(nameserversCmd)
	RootCmd.AddCommand(searchCmd)
	RootCmd.AddCommand(hostsCmd)
	RootCmd.AddCommand(blacklistCmd)
	RootCmd.AddCommand(whitelistCmd)
}
