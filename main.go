package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/arafato/cf-nuke/resources"
	"github.com/arafato/cf-nuke/utils"

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

	resources := infrastructure.ProcessCollection(&types.Credentials{
		AccountID: accountId,
		APIKey:    key,
	})

	fmt.Printf("Scan complete: Found %d removable resources in account %s:\n", resources.NumOf(types.Ready), accountId)
	utils.PrettyPrintStatus(resources)

	if !noDryRun {
		fmt.Println("Dry run complete.")

	} else {
		fmt.Println("Executing actual nuke operation... do you really want to continue (yes/no)?")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Nuke operation aborted.")
			return
		}
		fmt.Println("Nuke operation confirmed.")
		utils.PrettyPrintStatus(resources)

		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		wg.Add(1)
		go utils.PrintStatusWithContext(&wg, ctx, resources)

		if err := infrastructure.RemoveCollection(ctx, resources); err != nil {
			log.Printf("Error removing resources: %v", err)
		}

		cancel()
		// Waiting for everything to finish, in this case the status printer
		wg.Wait()

		fmt.Println("Process finished.")
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
