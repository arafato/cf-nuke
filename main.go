package main

import (
	"fmt"
	"os"

	_ "github.com/arafato/cf-nuke/resources"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/spf13/cobra"
)

// Global flags that can be used across commands
var (
	configFile string
	key        string
	accountId  string
	noDryRun   bool
)

var rootCmd = &cobra.Command{
	Use:   "cf-nuke",
	Short: "cf-nuke removes every resource from your cloudflare account",
	Long:  `A tool which removes every resource from an cloudflare account.  Use it with caution, since it cannot distinguish between production and non-production.`,
}

var nukeCmd = &cobra.Command{
	Use:   "nuke",
	Short: "Execute nuke operation",
	Long: `Nuke command performs destructive operations based on the provided configuration.
Use with caution and review the dry-run output before executing.`,
	Run: func(cmd *cobra.Command, args []string) {
		// This is where you'll add your custom logic
		executeNuke()
	},
}

func init() {
	// Add the nuke command to root
	rootCmd.AddCommand(nukeCmd)

	// Define flags for the nuke command
	nukeCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to configuration file (required)")
	nukeCmd.Flags().StringVarP(&key, "key", "k", "", "Key for operation (required)")
	nukeCmd.Flags().StringVarP(&accountId, "account-id", "a", "", "Cloudflare account id (required)")
	nukeCmd.Flags().BoolVar(&noDryRun, "no-dry-run", false, "Execute without dry run")

	// Make config and key required
	// nukeCmd.MarkFlagRequired("config")
	nukeCmd.MarkFlagRequired("account-id")
	nukeCmd.MarkFlagRequired("key")
}

func executeNuke() {
	fmt.Printf("Executing nuke command with:\n")
	fmt.Printf("  Config file: %s\n", configFile)
	fmt.Printf("  Account ID: %s\n", accountId)
	fmt.Printf("  No dry run: %t\n", noDryRun)

	// PLACEHOLDER: Load and parse configuration file
	/*
		config, err := loadConfig(configFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
		_ = config // Use the config in your logic
	*/
	if noDryRun {
		fmt.Println("Executing actual nuke operation...")

	} else {
		fmt.Println("Performing dry run...")

		resources := infrastructure.ProcessCollection(&types.Credentials{
			AccountID: accountId,
			APIKey:    key,
		})

		fmt.Printf("Scan complete: Found %d resources in account %s:\n", len(resources), accountId)
		for _, resource := range resources {
			fmt.Printf("%s - \033[32m%s\033[0m - %s\n", resource.ProductName, resource.ResourceName, resource.ResourceID)
		}
	}
}

func loadConfig(configPath string) {

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
