package coprompt

import (
	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

/*
HandleDynamicSuggestions
CoPrompt will check if your command includes the annotation "coprompt". If it's included the value will be
provided to the HandleDynamicSuggestions function.
*/
type HandleDynamicSuggestions func(annotation string, document prompt.Document) []prompt.Suggest

type CoPrompt struct {
	RootCmd                  *cobra.Command
	GoPromptOptions          []prompt.Option
	HandleDynamicSuggestions HandleDynamicSuggestions
}

const CallbackAnnotation = "coprompt"

/*
Run
CoPrompt will automatically generate suggestions for all your cobra commands and flags
CoPrompt will check if your command includes the annotation "coprompt". If it's included the value will be
provided to the HandleDynamicSuggestions function.
*/
func (coprompt CoPrompt) Run() {
	p := prompt.New(
		coprompt.copromtExecutor(),
		coprompt.copromptCompleter(),
		coprompt.GoPromptOptions...,
	)
	p.Run()
}

// executor executes command and print the output.
func (coprompt CoPrompt) copromtExecutor() func(string) {
	return func(in string) {
		promptArgs := strings.Split(in, " ")
		os.Args = append([]string{os.Args[0]}, promptArgs...)
		coprompt.RootCmd.Execute()
	}
}

// completer returns the completion items from user input.
func (coprompt CoPrompt) copromptCompleter() func(prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		currentCommand := coprompt.RootCmd
		cText := d.CurrentLine()
		args := strings.Fields(cText)

		for _, v := range args {
			if currentCommand.HasAvailableSubCommands() {
				for _, c := range currentCommand.Commands() {
					if c.Name() == v {
						currentCommand = c
						break
					}
				}
			} else {
				break
			}
		}

		suggestions := coprompt.collectSuggestions(currentCommand, d)

		return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
	}
}

func (coprompt CoPrompt) collectSuggestions(command *cobra.Command, d prompt.Document) []prompt.Suggest {
	var suggestions []prompt.Suggest

	persistentFlags := parentPersistentFlags(command)
	flags := append([]*pflag.FlagSet{command.LocalNonPersistentFlags()}, persistentFlags...)

	loopFlags := func(fn func(flag *pflag.Flag)) {
		for _, fs := range flags {
			fs.VisitAll(fn)
		}
	}

	if strings.HasPrefix(d.GetWordBeforeCursor(), "--") {
		loopFlags(func(flag *pflag.Flag) {
			if flag.Changed {
				flag.Value.Set(flag.DefValue)
			}
			suggestions = append(suggestions, prompt.Suggest{Text: "--" + flag.Name, Description: flag.Usage})
		})
	} else if strings.HasPrefix(d.GetWordBeforeCursor(), "-") {
		loopFlags(func(flag *pflag.Flag) {
			if flag.Changed {
				flag.Value.Set(flag.DefValue)
			}
			if flag.Shorthand != "" {
				suggestions = append(suggestions, prompt.Suggest{Text: "-" + flag.Shorthand, Description: flag.Usage})
			}
		})
	} else if command.HasAvailableSubCommands() {
		for _, c := range command.Commands() {
			if !c.Hidden {
				suggestions = append(suggestions, prompt.Suggest{Text: c.Name(), Description: c.Short})
			}
		}
	} else if coprompt.HandleDynamicSuggestions != nil && command.Annotations[CallbackAnnotation] != "" {
		copromptAnnotation := command.Annotations[CallbackAnnotation]
		suggestions = coprompt.HandleDynamicSuggestions(copromptAnnotation, d)
	}

	return suggestions
}

func parentPersistentFlags(cc *cobra.Command) []*pflag.FlagSet {
	fs := []*pflag.FlagSet{cc.PersistentFlags()}
	if cc.HasParent() {
		fs = append(fs, parentPersistentFlags(cc.Parent())...)
	}
	return fs
}
