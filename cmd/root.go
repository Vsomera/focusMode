package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "fmode",
	Short: `Cli tool to block distracting websites`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(`FocusMode Cli Tool`)
	},
}
