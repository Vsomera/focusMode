package cmd

import (
	"fmt"
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	"github.com/vsomera/focusmode/hosts"
)

var store = hosts.NewHostsStore()
var logo = Color(figure.NewFigure("Focus Mode", "smisome1", true).String(), Blue)

var rootCmd = &cobra.Command{
	Use:     "focusmode",
	Version: "1",
	Short:   Color(`FocusMode - Cli tool to block distracting websites`, Blue),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(logo)
		fmt.Print("\n| Cli Tool to block distracting websites, run \"focusmode help\" for command info.\n\n")
	},
}

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: `Lists all domains that are currently being blocked`,
	Run: func(cmd *cobra.Command, args []string) {
		domains, err := store.GetDomainsFromHost()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(domains) == 0 {
			fmt.Print(Color("\nBlacklist:\n\n", Blue))
			fmt.Print("|  No domains added\n\n")
			return
		}

		fmt.Print(Color("\nBlacklist:\n\n", Blue))
		for i, d := range domains {
			fmt.Printf("|  %v %s\n", i+1, d)
		}
		fmt.Println("")

	},
}

var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes selected domain in blacklist, can only remove 1 domain at a time",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Print(Color("\nError:\n\n|  No domain arguments in command call\n\n", Red))
			return
		}

		var confirm string
		domainToDelete := args[0]

		fmt.Print(Color(fmt.Sprintf("\nRemove %s from blacklist? \n\n|  Type 'y' to confirm [y/n] ", domainToDelete), Yellow))
		fmt.Scan(&confirm)

		if confirm == "y" || confirm == "Y" {

			err := store.DeleteDomainFromHost(domainToDelete)
			if err != nil {
				fmt.Print(Color(fmt.Sprintf("\nError:\n\n|  %s\n\n", err), Red))
				os.Exit(1)
			}

		}
		fmt.Println("")
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: `Add domain(s) to blacklist`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Print(Color("\nError:\n\n|  No domain arguments in command call\n\n", Red))
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

		fmt.Print(Color("\nAdded domain(s) to Blacklist:\n\n", Green))
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
	rootCmd.AddCommand(listCmd, addCmd, removeCmd)
}
