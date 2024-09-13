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
	Use:   "list",
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

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Removes all domains in blacklist",
	Run: func(cmd *cobra.Command, args []string) {
		selectedDomain, _ := cmd.Flags().GetString("d")
		action := Color("cleared all domains", Green)
		confirmMsg := Color("Clear all domains?", Yellow)

		if selectedDomain != "" {
			action = Color(fmt.Sprintf("removed %s from blacklist", selectedDomain), Green)
			confirmMsg = Color(fmt.Sprintf("Remove %s from blacklist?", selectedDomain), Yellow)
		}

		var confirm string

		fmt.Printf("\n%s\n\n|  Type 'y' to confirm [y/n] ", confirmMsg)
		fmt.Scan(&confirm)

		if confirm == "y" || confirm == "Y" {
			var err error
			if selectedDomain != "" {
				err = store.DeleteDomainFromHost(selectedDomain)
			} else {
				err = store.CleanDomains()
			}

			if err != nil {
				errMsg := fmt.Sprintf("\nError:\n\n|  %s\n\n", err)
				fmt.Print(Color(errMsg, Red))
				return
			}
			fmt.Print(Color(fmt.Sprintf("\n|  %s\n", action), Green))
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
	rootCmd.AddCommand(listCmd, addCmd, cleanCmd)
	cleanCmd.Flags().String("d", "", "deletes a selected domain in the blacklist matching the given string")
}
