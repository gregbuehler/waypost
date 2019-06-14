package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:     "version",
		Short:   "waypost version",
		Aliases: []string{"v"},
		Run:     showVersion,
	}
)

const version = "0.0.1"

func showVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("waypost %s", version)
}
