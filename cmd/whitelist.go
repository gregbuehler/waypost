package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	whitelistCmd = &cobra.Command{
		Use:   "whitelist",
		Short: "waypost whitelist",
	}
	addWhitelistCmd = &cobra.Command{
		Use:   "add [domain]",
		Short: "waypost whitelist add [domain]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			msg := fmt.Sprintf("whitelist add %s", strings.Join(args, " "))
			r, err := mgmtClient.IssueCommand(msg)
			if err != nil {
				fmt.Printf("error: %s", err.Error())
				os.Exit(1)
			}
			fmt.Print(r)
			os.Exit(0)
		},
	}
	delWhitelistCmd = &cobra.Command{
		Use:   "del [domain]",
		Short: "waypost whitelist del [domain]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			msg := fmt.Sprintf("whitelist del %s", strings.Join(args, " "))
			r, err := mgmtClient.IssueCommand(msg)
			if err != nil {
				fmt.Printf("error: %s", err.Error())
				os.Exit(1)
			}
			fmt.Print(r)
			os.Exit(0)
		},
	}
	listWhitelistCmd = &cobra.Command{
		Use:   "list",
		Short: "waypost whitelist list",
		Run: func(cmd *cobra.Command, args []string) {
			r, err := mgmtClient.IssueCommand("whitelist list")
			if err != nil {
				fmt.Printf("error: %s", err.Error())
				os.Exit(1)
			}
			fmt.Print(r)
			os.Exit(0)
		},
	}
	enableWhitelistCmd = &cobra.Command{
		Use:   "enable",
		Short: "waypost whitelist enable",
		Run: func(cmd *cobra.Command, args []string) {
			r, err := mgmtClient.IssueCommand("whitelist enable")
			if err != nil {
				fmt.Printf("error: %s", err.Error())
				os.Exit(1)
			}
			fmt.Print(r)
			os.Exit(0)
		},
	}
	disableWhitelistCmd = &cobra.Command{
		Use:   "disable",
		Short: "waypost whitelist disable",
		Run: func(cmd *cobra.Command, args []string) {
			r, err := mgmtClient.IssueCommand("whitelist disable")
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
	whitelistCmd.AddCommand(addWhitelistCmd)
	whitelistCmd.AddCommand(delWhitelistCmd)
	whitelistCmd.AddCommand(listWhitelistCmd)
	whitelistCmd.AddCommand(enableWhitelistCmd)
	whitelistCmd.AddCommand(disableWhitelistCmd)
}
