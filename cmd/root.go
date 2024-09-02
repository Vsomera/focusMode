package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vsomera/focusmode/hosts"
)

// TODO : implement cleaning domains, and selecting domains to delete

var store = hosts.NewHostsStore()

var rootCmd = &cobra.Command{
	Use:   "fmode",
	Short: `Cli tool to block distracting websites`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("\nFocusMode\n\n| Cli Tool to block distracting websites, focus on what actually matters.\n\n")
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: `Lists all domains that are currently being blocked`,
	Run: func(cmd *cobra.Command, args []string) {
		domains, err := store.GetDomainsFromHost()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Print("\nBlacklist:\n\n")
		for i, d := range domains {
			fmt.Printf("|  %v %s\n", i+1, d)
		}
		fmt.Println("")

	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: `Add domain(s) to block list`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Print("\nError:\n\n|  No domain arguments in command call\n\n")
			return
		}

		existingDomains, err := store.GetDomainsFromHost()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		domainsToAdd := append(existingDomains, args...)
		err = store.AddDomainsToHost(domainsToAdd)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Print("\nAdded domain(s) to Blacklist:\n\n")
		for i, d := range args {
			fmt.Printf("|  %v %s\n", i+1, d)
		}
		fmt.Println("")

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer store.Close()
}

func init() {
	rootCmd.AddCommand(listCmd, addCmd)
}
