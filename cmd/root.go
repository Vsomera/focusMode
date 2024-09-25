package cmd

import (
	"fmt"
	"os"

	"github.com/Vsomera/focusmode/hosts"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
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

// TODO : implement '--a' flag
var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes selected or all domains in blacklist",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Print(Color("\nError:\n\n|  No domain arguments or flag in command call\n\n", Red))
			return
		}

		var confirm string

		confirmMsg := fmt.Sprintf("\nDelete %v domain(s)? \n\n|  Type 'y' to confirm [y/n] ", len(args))
		fmt.Print(Color(confirmMsg, Yellow))
		fmt.Scan(&confirm)

		if confirm == "y" || confirm == "Y" {
			var err error
			// delete selected domains in args
			for _, d := range args {
				err = store.DeleteDomainFromHost(d)
			}

			if err != nil {
				errMsg := fmt.Sprintf("\nError:\n\n|  %s\n\n", err)
				fmt.Print(Color(errMsg, Red))
				os.Exit(1)
			}

			deleteMsg := fmt.Sprintf("\nDeleted %d domain(s) from Blacklist:\n\n", len(args))
			fmt.Print(Color(deleteMsg, Green))
			for i, d := range args {
				fmt.Printf("|  %v %s\n", i+1, d)
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
	removeCmd.Flags().String("a", "", "deletes all domains")
}
