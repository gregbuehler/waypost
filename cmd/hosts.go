package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	hostsCmd = &cobra.Command{
		Use:   "hosts",
		Short: "waypost hosts",
	}

	addHostCmd = &cobra.Command{
		Use:   "add [name] [ip]",
		Short: "waypost hosts add [name] [ip]",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			msg := fmt.Sprintf("hosts add %s", strings.Join(args, " "))
			r, err := mgmtClient.IssueCommand(msg)
			if err != nil {
				fmt.Printf("error: %s", err.Error())
				os.Exit(1)
			}
			fmt.Print(r)
			os.Exit(0)
		},
	}
	delHostCmd = &cobra.Command{
		Use:   "del [name] [ip]",
		Short: "waypost hosts del [name]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			msg := fmt.Sprintf("hosts del %s", strings.Join(args, " "))
			r, err := mgmtClient.IssueCommand(msg)
			if err != nil {
				fmt.Printf("error: %s", err.Error())
				os.Exit(1)
			}
			fmt.Print(r)
			os.Exit(0)
		},
	}
	listHostsCmd = &cobra.Command{
		Use:   "list",
		Short: "waypost hosts list",
		Run: func(cmd *cobra.Command, args []string) {
			r, err := mgmtClient.IssueCommand("hosts list")
			if err != nil {
				fmt.Printf("error: %s", err.Error())
				os.Exit(1)
			}
			fmt.Print(r)
			os.Exit(0)
		},
	}
)

func init() {
	hostsCmd.AddCommand(addHostCmd)
	hostsCmd.AddCommand(delHostCmd)
	hostsCmd.AddCommand(listHostsCmd)
}
