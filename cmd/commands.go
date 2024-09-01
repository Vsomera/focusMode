package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vsomera/focusmode/hosts"
)

var store = hosts.NewHostsStore()

func init() {
	rootCmd.AddCommand(listCmd)
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
		fmt.Printf(`List of blocked domains: %v`, domains)
	},
}
