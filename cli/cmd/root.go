package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
	commit  = "dev"
	date    = "2026-10-27"

	rootCmd = &cobra.Command{
		Use:   "launchpad",
		Short: "Launchpad CLI - Developer interface for the Internal Developer Platform",
		Long: `Launchpad is an Internal Developer Platform (IDP) CLI that empowers software engineers 
to scaffold, validate, render, and monitor production-ready microservices with secure-by-default guardrails.`,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(renderCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(versionCmd)
}

