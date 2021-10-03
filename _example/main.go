package main

import (
	"github.com/c-bata/go-prompt"
	cobraprompt "github.com/stromland/cobra-prompt"
	"github.com/stromland/cobra-prompt/_example/cmd"
)

func main() {
	shell := &cobraprompt.CobraPrompt{
		RootCmd:                  cmd.RootCmd,
		PersistFlagValues:        true,
		ShowHelpCommandAndFlags:  true,
		DisableCompletionCommand: true,
		AddDefaultExitCommand:    true,
		GoPromptOptions: []prompt.Option{
			prompt.OptionTitle("cobra-prompt"),
			prompt.OptionPrefix(">(^'^)> "),
			prompt.OptionMaxSuggestion(10),
		},
		DynamicSuggestionsFunc: func(annotationValue string, document *prompt.Document) []prompt.Suggest {
			if suggestions := cmd.GetFoodDynamic(annotationValue); suggestions != nil {
				return suggestions
			}

			return []prompt.Suggest{}
		},
	}
	shell.Run()
}
