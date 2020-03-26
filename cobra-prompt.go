package cobraprompt

import (
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CallbackAnnotation for dynamic suggestions.
const CallbackAnnotation = "cobra-prompt"

// PersistFlagValuesFlag the flag that will be avaiailable when PersistFlagValues is true
const PersistFlagValuesFlag = "persist-flag-values"

// CobraPrompt given a Cobra command it will make every flag and sub commands available as suggestions.
// Command.Short will be used as description for the suggestion.
type CobraPrompt struct {
	// RootCmd is the start point, all its sub commands and flags will be available as suggestions
	RootCmd *cobra.Command

	// GoPromptOptions is for customize go-prompt
	// see https://github.com/c-bata/go-prompt/blob/master/option.go
	GoPromptOptions []prompt.Option

	// DynamicSuggestionsFunc will be executed if an command has CallbackAnnotation as an annotation. If it's included
	// the value will be provided to the DynamicSuggestionsFunc function.
	DynamicSuggestionsFunc func(annotation string, document *prompt.Document) []prompt.Suggest

	// PersistFlagValues will persist flags. For example have verbose turned on every command.
	PersistFlagValues bool
}

// Run will automatically generate suggestions for all cobra commands and flags defined by RootCmd
// and execute the selected commands. Run will also reset all given flags by default, see PersistFlagValues
func (co CobraPrompt) Run() {
	co.prepare()
	p := prompt.New(
		func(in string) {
			promptArgs := strings.Fields(in)
			os.Args = append([]string{os.Args[0]}, promptArgs...)
			co.RootCmd.Execute()
		},
		func(d prompt.Document) []prompt.Suggest {
			return findSuggestions(&co, &d)
		},
		co.GoPromptOptions...,
	)
	p.Run()
}

func (co CobraPrompt) prepare() {
	if co.PersistFlagValues {
		co.RootCmd.PersistentFlags().BoolP(PersistFlagValuesFlag, "",
			false, "Persist last given value for flags")
	}
}

func findSuggestions(co *CobraPrompt, d *prompt.Document) []prompt.Suggest {
	command := co.RootCmd
	args := strings.Fields(d.CurrentLine())

	if found, _, err := command.Find(args); err == nil {
		command = found
	}

	var suggestions []prompt.Suggest
	persistFlagValues, _ := command.Flags().GetBool(PersistFlagValuesFlag)
	addFlags := func(flag *pflag.Flag) {
		if flag.Changed && !persistFlagValues {
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

	annotation := command.Annotations[CallbackAnnotation]
	if co.DynamicSuggestionsFunc != nil && annotation != "" {
		suggestions = append(suggestions, co.DynamicSuggestionsFunc(annotation, d)...)
	}
	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}
