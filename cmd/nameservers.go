package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	nameserversCmd = &cobra.Command{
		Use:   "nameservers",
		Short: "waypost nameservers",
	}

	addNameserverCmd = &cobra.Command{
		Use:   "add [ip]",
		Short: "waypost nameservers add [ip]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			msg := fmt.Sprintf("nameservers add %s", strings.Join(args, " "))
			r, err := mgmtClient.IssueCommand(msg)
			if err != nil {
				fmt.Printf("error: %s", err.Error())
				os.Exit(1)
			}
			fmt.Print(r)
			os.Exit(0)
		},
	}
	delNameserverCmd = &cobra.Command{
		Use:   "del [ip]",
		Short: "waypost nameservers del [ip]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			msg := fmt.Sprintf("nameservers del %s", strings.Join(args, " "))
			r, err := mgmtClient.IssueCommand(msg)
			if err != nil {
				fmt.Printf("error: %s", err.Error())
				os.Exit(1)
			}
			fmt.Print(r)
			os.Exit(0)
		},
	}
	listNameserversCmd = &cobra.Command{
		Use:   "list",
		Short: "waypost nameservers list",
		Run: func(cmd *cobra.Command, args []string) {
			r, err := mgmtClient.IssueCommand("nameservers list")
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
	nameserversCmd.AddCommand(addNameserverCmd)
	nameserversCmd.AddCommand(delNameserverCmd)
	nameserversCmd.AddCommand(listNameserversCmd)
}
