package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var exitCmd = &cobra.Command{
	Use:   "exit",
	Short: "Exit prompt",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}

func init() {
	RootCmd.AddCommand(exitCmd)
}
