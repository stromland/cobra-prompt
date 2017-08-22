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

// Run will automatically generate suggestions for all your cobra commands and flags and execute the commands
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

	for _, arg := range args {
		if command.HasAvailableSubCommands() {
			for _, cmd := range command.Commands() {
				if cmd.Name() == arg || isAlias(arg, cmd.Aliases) {
					command = cmd
					break
				}
			}
		} else {
			break
		}
	}

	var suggestions []prompt.Suggest
	command.Flags().VisitAll(func(flag *pflag.Flag) {
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
	})

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

func isAlias(name string, aliases []string) bool {
	for _, alias := range aliases {
		if name == alias {
			return true
		}
	}
	return false
}
