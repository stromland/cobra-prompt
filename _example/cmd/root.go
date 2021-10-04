package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	SilenceUsage:  true, // Only print usage when defined in command.
	SilenceErrors: true,
}
