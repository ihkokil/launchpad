package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print Launchpad CLI version and environment information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Launchpad IDP CLI version %s (commit: %s, built: %s, %s/%s)\n",
			version, commit, date, runtime.GOOS, runtime.GOARCH)
	},
}
