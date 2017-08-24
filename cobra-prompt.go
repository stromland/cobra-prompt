package cobraprompt

import (
	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

// CALLBACK_ANNOTATION
const CALLBACK_ANNOTATION = "cobra-prompt"

// CobraPrompt struct
// DynamicSuggestionsFunc will be executed if an command has CALLBACK_ANNOTATION as an annotation. If it's included the
// value will be provided to the DynamicSuggestionsFunc function.
type CobraPrompt struct {
	RootCmd                *cobra.Command
	GoPromptOptions        []prompt.Option
	DynamicSuggestionsFunc func(annotation string, document prompt.Document) []prompt.Suggest
	ResetFlagsFlag         bool
}

// Run will automatically generate suggestions for all your cobra commands and flags and execute the selected commands
func (co CobraPrompt) Run() {
	co.prepare()
	p := prompt.New(
		func(in string) {
			promptArgs := strings.Fields(in)
			os.Args = append([]string{os.Args[0]}, promptArgs...)
			co.RootCmd.Execute()
		},
		func(d prompt.Document) []prompt.Suggest {
			return findSuggestions(co, d)
		},
		co.GoPromptOptions...,
	)
	p.Run()
}

func (co CobraPrompt) prepare() {
	if co.ResetFlagsFlag {
		co.RootCmd.PersistentFlags().BoolP("flags-no-reset", "",
			false, "Turn off flags explicit")
	}
}

func findSuggestions(co CobraPrompt, d prompt.Document) []prompt.Suggest {
	command := co.RootCmd
	args := strings.Fields(d.CurrentLine())

	if found, _, err := command.Find(args); err == nil {
		command = found
	}

	var suggestions []prompt.Suggest
	resetFlags, _ := command.Flags().GetBool("flags-no-reset")
	addFlags := func(flag *pflag.Flag) {
		if flag.Changed && !resetFlags {
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

	annotation := command.Annotations[CALLBACK_ANNOTATION]
	if co.DynamicSuggestionsFunc != nil && annotation != "" {
		suggestions = append(suggestions, co.DynamicSuggestionsFunc(annotation, d)...)
	}
	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}
