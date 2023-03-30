package cmd

import (
	"github.com/avirtopeanu-ionos/cobra"
)

var RootCmd = &cobra.Command{
	Use:           "cobra-prompt",
	SilenceUsage:  true, // Only print usage when defined in command.
	SilenceErrors: true,
}
