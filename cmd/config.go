package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "waypost config",
	}

	reloadConfigCmd = &cobra.Command{
		Use:   "reload [name] [ip]",
		Short: "waypost config reload [name] [ip]",
		Run: func(cmd *cobra.Command, args []string) {
			msg := fmt.Sprintf("config reload")
			r, err := mgmtClient.IssueCommand(msg)
			if err != nil {
				fmt.Printf("error: %s", err.Error())
				os.Exit(1)
			}
			fmt.Print(r)
			os.Exit(0)
		},
	}
	dumpConfigCmd = &cobra.Command{
		Use:   "dump [name] [ip]",
		Short: "waypost config dump [name]",
		Run: func(cmd *cobra.Command, args []string) {
			msg := fmt.Sprintf("config dump")
			r, err := mgmtClient.IssueCommand(msg)
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
	configCmd.AddCommand(reloadConfigCmd)
	configCmd.AddCommand(dumpConfigCmd)
}
