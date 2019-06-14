package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cacheCmd = &cobra.Command{
		Use:   "cache",
		Short: "waypost cache [flush|list|enable|disable]",
	}
	flushCacheCmd = &cobra.Command{
		Use:   "flush",
		Short: "waypost cache flush",
		Run: func(cmd *cobra.Command, args []string) {
			r, err := mgmtClient.IssueCommand("cache flush")
			if err != nil {
				fmt.Printf("error: %s", err.Error())
				os.Exit(1)
			}
			fmt.Print(r)
			os.Exit(0)
		},
	}
	listCacheCmd = &cobra.Command{
		Use:   "list",
		Short: "waypost cache list",
		Run: func(cmd *cobra.Command, args []string) {
			r, err := mgmtClient.IssueCommand("cache list")
			if err != nil {
				fmt.Printf("error: %s", err.Error())
				os.Exit(1)
			}
			fmt.Print(r)
			os.Exit(0)
		},
	}
	enableCacheCmd = &cobra.Command{
		Use:   "enable",
		Short: "waypost cache enable",
		Run: func(cmd *cobra.Command, args []string) {
			r, err := mgmtClient.IssueCommand("cache enable")
			if err != nil {
				fmt.Printf("error: %s", err.Error())
				os.Exit(1)
			}
			fmt.Print(r)
			os.Exit(0)
		},
	}
	disableCacheCmd = &cobra.Command{
		Use:   "disable",
		Short: "waypost cache disable",
		Run: func(cmd *cobra.Command, args []string) {
			r, err := mgmtClient.IssueCommand("cache disable")
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
	cacheCmd.AddCommand(flushCacheCmd)
	cacheCmd.AddCommand(listCacheCmd)
	cacheCmd.AddCommand(enableCacheCmd)
	cacheCmd.AddCommand(disableCacheCmd)
}
