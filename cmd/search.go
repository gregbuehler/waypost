package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	searchCmd = &cobra.Command{
		Use:   "search",
		Short: "waypost search",
	}

	addSearchCmd = &cobra.Command{
		Use:   "add [ip]",
		Short: "waypost search add [domain]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			msg := fmt.Sprintf("search add %s", strings.Join(args, " "))
			r, err := mgmtClient.IssueCommand(msg)
			if err != nil {
				fmt.Printf("error: %s", err.Error())
				os.Exit(1)
			}
			fmt.Print(r)
			os.Exit(0)
		},
	}
	delSearchCmd = &cobra.Command{
		Use:   "del [ip]",
		Short: "waypost search del [domain]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			msg := fmt.Sprintf("search del %s", strings.Join(args, " "))
			r, err := mgmtClient.IssueCommand(msg)
			if err != nil {
				fmt.Printf("error: %s", err.Error())
				os.Exit(1)
			}
			fmt.Print(r)
			os.Exit(0)
		},
	}
	listSearchsCmd = &cobra.Command{
		Use:   "list",
		Short: "waypost search list",
		Run: func(cmd *cobra.Command, args []string) {
			r, err := mgmtClient.IssueCommand("search list")
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
	searchCmd.AddCommand(addSearchCmd)
	searchCmd.AddCommand(delSearchCmd)
	searchCmd.AddCommand(listSearchsCmd)
}
