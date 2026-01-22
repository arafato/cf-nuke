package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/arafato/cf-nuke/config"
	"github.com/arafato/cf-nuke/infrastructure"
	_ "github.com/arafato/cf-nuke/resources"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
	"github.com/arafato/cf-nuke/version"
	"github.com/spf13/cobra"
)

// Global flags that can be used across commands
var (
	configFile   string
	key          string
	accountId    string
	user         string
	mode         string
	noDryRun     bool
	shortVersion bool
)

var rootCmd = &cobra.Command{
	Use:   "cf-nuke",
	Short: "cf-nuke removes every resource from your cloudflare account",
	Long:  `A tool which removes every resource from an cloudflare account.  Use it with caution, since it cannot distinguish between production and non-production.`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of cf-nuke",
	Run: func(cmd *cobra.Command, args []string) {
		if shortVersion {
			version.PrintShort(os.Stdout)
		} else {
			version.Print(os.Stdout)
		}
	},
}

var nukeCmd = &cobra.Command{
	Use:   "nuke",
	Short: "Execute nuke operation",
	Long: `Nuke command performs destructive operations based on the provided configuration.
Use with caution and review the dry-run output before executing.`,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		mode, _ := cmd.Flags().GetString("mode")
		if mode != "token" && mode != "account" {
			return fmt.Errorf("invalid mode '%s', must be 'token' or 'account'", mode)
		}

		if mode == "account" {
			user, _ := cmd.Flags().GetString("user")
			if user == "" {
				return fmt.Errorf("--user flag is required when --mode is 'account'")
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		executeNuke()
	},
}

func init() {
	rootCmd.AddCommand(nukeCmd)
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolVar(&shortVersion, "short", false, "Print short version string")

	nukeCmd.Flags().StringVarP(&mode, "mode", "m", "", "The mode of operation ('token' or 'account')")
	nukeCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to configuration file (required)")
	nukeCmd.Flags().StringVarP(&key, "key", "k", "", "Key for operation (required)")
	nukeCmd.Flags().StringVarP(&accountId, "account-id", "a", "", "Cloudflare account id (required)")
	nukeCmd.Flags().BoolVar(&noDryRun, "no-dry-run", false, "Execute without dry run")
	nukeCmd.Flags().StringVarP(&user, "user", "u", "", "The user identifier (required only for 'account' mode)")

	nukeCmd.MarkFlagRequired("config")
	nukeCmd.MarkFlagRequired("account-id")
	nukeCmd.MarkFlagRequired("key")
}

func executeNuke() {
	config, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
		os.Exit(1)
	}

	creds := &types.Credentials{
		AccountID: accountId,
		APIKey:    key,
		User:      user,
		Mode:      types.Mode(mode), // this is safe due to pre-check in PreRunE
	}

	creds.S3AccessKeyID, creds.S3AccessSecret, err = utils.CreateTemporaryR2Token(creds)
	time.Sleep(3 * time.Second)
	if err != nil {
		log.Fatalf("Error creating temporary S3/R2 token: %v", err)
		os.Exit(1)
	}

	resources := infrastructure.ProcessCollection(creds)
	infrastructure.FilterCollection(resources, config)

	visibleCount := resources.VisibleCount()
	fmt.Printf("Scan complete: Found %d resources in total in account %s: To be removed %d, Filtered %d\n", visibleCount, accountId, resources.NumOf(types.Ready), resources.NumOf(types.Filtered))
	utils.PrettyPrintStatus(resources)

	if !noDryRun {
		fmt.Println("Dry run complete.")

	} else {
		fmt.Println("Executing actual nuke operation... do you really want to continue (yes/no)?")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			err := utils.DeleteTemporaryR2Token(creds, resources)
			if err != nil {
				fmt.Println("Failed to delete temporary account token for R2/S3 operations:", err)
			}
			fmt.Println("Nuke operation aborted.")
			return
		}
		fmt.Println("Nuke operation confirmed.")

		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())

		// Start printer goroutine BEFORE removal to show progress during the operation
		wg.Add(1)
		go utils.PrintStatusWithContext(&wg, ctx, resources)

		if err := infrastructure.RemoveCollection(ctx, resources); err != nil {
			log.Printf("Error removing resources: %v", err)
		}

		// Cancel printer after removal completes, then wait for it to finish
		cancel()
		wg.Wait()

		err := utils.DeleteTemporaryR2Token(creds, resources)
		if err != nil {
			fmt.Println("Failed to delete temporary account token for R2/S3 operations:", err)
		}

		fmt.Println("Process finished.")
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
