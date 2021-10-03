package main

import (
	"github.com/c-bata/go-prompt"
	cobraprompt "github.com/stromland/cobra-prompt"
	"github.com/stromland/cobra-prompt/_example/cmd"
)

func main() {
	shell := &cobraprompt.CobraPrompt{
		RootCmd:                  cmd.RootCmd,
		DynamicSuggestionsFunc:   handleDynamicSuggestions,
		PersistFlagValues:        true,
		ShowHelpCommandAndFlags:  true,
		DisableCompletionCommand: true,
		AddDefaultExitCommand:    true,
		GoPromptOptions: []prompt.Option{
			prompt.OptionTitle("cobra-prompt"),
			prompt.OptionPrefix(">(^'^)> "),
			prompt.OptionMaxSuggestion(10),
		},
	}
	shell.Run()
}

func handleDynamicSuggestions(annotationValue string, _ *prompt.Document) []prompt.Suggest {
	switch annotationValue {
	case cmd.GET_FOOD_DYNAMIC:
		return GetFood()
	default:
		return []prompt.Suggest{}
	}
}

func GetFood() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "apple", Description: "Green apple"},
		{Text: "tomato", Description: "Red tomato"},
	}
}
