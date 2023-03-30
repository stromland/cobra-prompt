package cmd

import (
	"github.com/avirtopeanu-ionos/cobra"
	"github.com/c-bata/go-prompt"
	cobraprompt "github.com/stromland/cobra-prompt"
)

var getFoodDynamicAnnotationValue = "GetFood"

var GetFoodDynamic = func(annotationValue string) []prompt.Suggest {
	if annotationValue != getFoodDynamicAnnotationValue {
		return nil
	}

	return []prompt.Suggest{
		{Text: "apple", Description: "Green apple"},
		{Text: "tomato", Description: "Red tomato"},
	}
}

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
		cobraprompt.DynamicSuggestionsAnnotation: getFoodDynamicAnnotationValue,
	},
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		name, _ := cmd.Flags().GetString("name")
		for _, v := range args {
			if verbose {
				cmd.Printf("Here you go, take this from %s: %s\n", name, v)
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
	getCmd.PersistentFlags().StringP("name", "n", "John", "name of the person")
	_ = getCmd.RegisterFlagCompletionFunc("name", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"John", "Mary", "Anne"}, cobra.ShellCompDirectiveNoFileComp
	})
}
