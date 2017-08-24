package coprompt

import (
	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

// CALLBACK_ANNOTATION
const CALLBACK_ANNOTATION = "coprompt"

// CoPrompt struct
// DynamicSuggestionsFunc will be executed if an command has CALLBACK_ANNOTATION as an annotation. If it's included the
// value will be provided to the DynamicSuggestionsFunc function.
type CoPrompt struct {
	RootCmd                *cobra.Command
	GoPromptOptions        []prompt.Option
	DynamicSuggestionsFunc func(annotation string, document prompt.Document) []prompt.Suggest
}

// Run will automatically generate suggestions for all your cobra commands and flags and execute the selected commands
func (coprompt CoPrompt) Run() {
	p := prompt.New(
		func(in string) {
			promptArgs := strings.Fields(in)
			os.Args = append([]string{os.Args[0]}, promptArgs...)
			coprompt.RootCmd.Execute()
		},
		func(d prompt.Document) []prompt.Suggest {
			return findSuggestions(coprompt, d)
		},
		coprompt.GoPromptOptions...,
	)
	p.Run()
}

func findSuggestions(coprompt CoPrompt, d prompt.Document) []prompt.Suggest {
	command := coprompt.RootCmd
	args := strings.Fields(d.CurrentLine())

	if found, _, err := command.Find(args); err == nil {
		command = found
	}

	var suggestions []prompt.Suggest
	addFlags := func(flag *pflag.Flag) {
		if flag.Changed {
			flag.Value.Set(flag.DefValue)
		}
		if flag.Hidden {
			return
		}
		if strings.HasPrefix(d.GetWordBeforeCursor(), "--") {
			suggestions = append(suggestions, prompt.Suggest{Text: "--" + flag.Name, Description: flag.Usage})
		} else if strings.HasPrefix(d.GetWordBeforeCursor(), "-") && flag.Shorthand != "" {
			suggestions = append(suggestions, prompt.Suggest{Text: "-" + flag.Shorthand, Description: flag.Usage})
		}
	}

	command.LocalFlags().VisitAll(addFlags)
	command.InheritedFlags().VisitAll(addFlags)

	if command.HasAvailableSubCommands() {
		for _, c := range command.Commands() {
			if !c.Hidden {
				suggestions = append(suggestions, prompt.Suggest{Text: c.Name(), Description: c.Short})
			}
		}
	}

	if coprompt.DynamicSuggestionsFunc != nil && command.Annotations[CALLBACK_ANNOTATION] != "" {
		annotation := command.Annotations[CALLBACK_ANNOTATION]
		suggestions = append(suggestions, coprompt.DynamicSuggestionsFunc(annotation, d)...)
	}
	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}
