package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/lcrownover/banner-tools/internal/banner"
	"github.com/lcrownover/banner-tools/internal/keys"
	"github.com/spf13/cobra"
)

func main() {
	apiKey, exists := os.LookupEnv("BANNER_TOOLS_API_KEY")
	if !exists {
		fmt.Println("BANNER_TOOLS_API_KEY environment variable not set")
		os.Exit(1)
	}
	apiURL, exists := os.LookupEnv("BANNER_TOOLS_API_URL")
	if !exists {
		fmt.Println("BANNER_TOOLS_API_URL environment variable not set ('api.example.edu')")
		os.Exit(1)
	}
	ctx := context.WithValue(context.Background(), keys.ApiKeyKey, apiKey)
	ctx = context.WithValue(context.Background(), keys.ApiURL, apiURL)

	// Root command
	var rootCmd = &cobra.Command{
		Use:   "banner-tool",
		Short: "A Banner CLI tool to get data from the Banner API",
	}

	// "duckid" subcommand
	var duckidCmd = &cobra.Command{
		Use:   "duckid <bannerID>,<bannerID>",
		Short: "Fetches Duck IDs from any number of given Banner IDs",
		Args:  cobra.ExactArgs(1), // Expects exactly 1 positional argument
		Run: func(cmd *cobra.Command, args []string) {
			bannerId := args[0]
			bannerIds := strings.Split(bannerId, ",")
			fmt.Printf("Banner IDs: %v\n", bannerIds)
			var wg sync.WaitGroup
			wg.Add(len(bannerIds))
			for _, bid := range bannerIds {
				go func() {
					defer wg.Done()
					duckid, err := banner.BannerIDToDuckID(ctx, bid)
					if err != nil {
						fmt.Printf("failed to get duckid from bannerid: %v\n", err)
					}
					fmt.Printf("%s\n", duckid)
				}()
			}
			wg.Wait()
		},
	}

	// "roster" subcommand
	var rosterCmd = &cobra.Command{
		Use:   "roster [termCode] [crn]",
		Short: "Gets a list of Banner IDs from a given class registration number and term code",
		Args:  cobra.ExactArgs(2), // Expects exactly 2 positional arguments
		Run: func(cmd *cobra.Command, args []string) {
			termCode := args[0]
			crn := args[1]
			fmt.Printf("Roster - Term Code: %s, CRN: %s\n", termCode, crn)
		},
	}

	// "term" subcommand
	var nextTerm bool
	var termCmd = &cobra.Command{
		Use:   "term",
		Short: "Shows the current or next term",
		Run: func(cmd *cobra.Command, args []string) {
			if nextTerm {
				fmt.Println("Term: next")
			} else {
				fmt.Println("Term: current")
			}
		},
	}

	// Set up the "next" flag for the "term" command
	termCmd.Flags().BoolVar(&nextTerm, "next", false, "Get the next term")

	// Add subcommands to the root command
	rootCmd.AddCommand(duckidCmd)
	rootCmd.AddCommand(rosterCmd)
	rootCmd.AddCommand(termCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
