package cmd

import (
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:     "get",
	Short:   "Get something",
	Aliases: []string{"eat"},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var getFoodCmd = &cobra.Command{
	Use:   "food",
	Short: "Get some food",
	Annotations: map[string]string{
		CallbackAnnotation: "GetFood",
	},
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		for _, v := range args {
			if verbose {
				cmd.Println("Here you go, take this:", v)
			} else {
				cmd.Println(v)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getFoodCmd)
	getCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose log")
}
