package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Global flags that can be used across commands
var (
	configFile string
	key        string
	accountId  string
	noDryRun   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cf-nuke",
	Short: "cf-nuke removes every resource from your cloudflare account",
	Long:  `A tool which removes every resource from an cloudflare account.  Use it with caution, since it cannot distinguish between production and non-production.`,
}

// nukeCmd represents the nuke command
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
	// PLACEHOLDER: Execute your main logic
	if noDryRun {
		fmt.Println("Executing actual nuke operation...")
		// Your actual destructive operation here

	} else {
		fmt.Println("Performing dry run...")
		// Your dry run logic here
	}

	fmt.Println("Nuke operation completed successfully!")
}

// PLACEHOLDER: Implement your config loading logic
func loadConfig(configPath string) {

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

/*
 var allItems []*Item
    var scanWg sync.WaitGroup
    itemCollectChan := make(chan *Item, 1000)

    // Start all scanners
    scanners := []Scanner{dnsScanner, workerScanner, pageRuleScanner}
    for _, scanner := range scanners {
        scanWg.Add(1)
        go func(s Scanner) {
            defer scanWg.Done()
            s.Scan(itemCollectChan) // Just collect, don't delete yet
        }(scanner)
    }

    // Collect items while scanners run
    go func() {
        scanWg.Wait()           // Wait for ALL scanners to finish
        close(itemCollectChan)  // Then close collection channel
    }()

    // Gather all items
    for item := range itemCollectChan {
        allItems = append(allItems, item)
    }

    fmt.Printf("✅ Scanning complete. Found %d items\n", len(allItems))

    // Phase 2: NOW start deletion with complete list
    startDeletionPhase(allItems)
*/
