package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	blacklistCmd = &cobra.Command{
		Use:   "blacklist",
		Short: "waypost blacklist",
	}
	addBlacklistCmd = &cobra.Command{
		Use:   "add [domain]",
		Short: "waypost blacklist add [domain]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			msg := fmt.Sprintf("blacklist add %s", strings.Join(args, " "))
			r, err := mgmtClient.IssueCommand(msg)
			if err != nil {
				fmt.Printf("error: %s", err.Error())
				os.Exit(1)
			}
			fmt.Print(r)
			os.Exit(0)
		},
	}
	delBlacklistCmd = &cobra.Command{
		Use:   "del [domain]",
		Short: "waypost blacklist del [domain]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			msg := fmt.Sprintf("blacklist del %s", strings.Join(args, " "))
			r, err := mgmtClient.IssueCommand(msg)
			if err != nil {
				fmt.Printf("error: %s", err.Error())
				os.Exit(1)
			}
			fmt.Print(r)
			os.Exit(0)
		},
	}
	listBlacklistCmd = &cobra.Command{
		Use:   "list",
		Short: "waypost blacklist list",
		Run: func(cmd *cobra.Command, args []string) {
			r, err := mgmtClient.IssueCommand("blacklist list")
			if err != nil {
				fmt.Printf("error: %s", err.Error())
				os.Exit(1)
			}
			fmt.Print(r)
			os.Exit(0)
		},
	}
	enableBlacklistCmd = &cobra.Command{
		Use:   "enable",
		Short: "waypost blacklist enable",
		Run: func(cmd *cobra.Command, args []string) {
			r, err := mgmtClient.IssueCommand("blacklist enable")
			if err != nil {
				fmt.Printf("error: %s", err.Error())
				os.Exit(1)
			}
			fmt.Print(r)
			os.Exit(0)
		},
	}
	disableBlacklistCmd = &cobra.Command{
		Use:   "disable",
		Short: "waypost blacklist disable",
		Run: func(cmd *cobra.Command, args []string) {
			r, err := mgmtClient.IssueCommand("blacklist disable")
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
	blacklistCmd.AddCommand(addBlacklistCmd)
	blacklistCmd.AddCommand(delBlacklistCmd)
	blacklistCmd.AddCommand(listBlacklistCmd)
	blacklistCmd.AddCommand(enableBlacklistCmd)
	blacklistCmd.AddCommand(disableBlacklistCmd)
}
