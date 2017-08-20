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

var lRootCmd *cobra.Command
var lHandleDynamicSuggestions HandleDynamicSuggestions

/*
Run
CoPrompt will automatically generate suggestions for all your cobra commands and flags
CoPrompt will check if your command includes the annotation "coprompt". If it's included the value will be
provided to the HandleDynamicSuggestions function.
*/
func Run(rootCmd *cobra.Command, fn HandleDynamicSuggestions, options []prompt.Option) {
	lRootCmd = rootCmd
	lHandleDynamicSuggestions = fn

	p := prompt.New(
		executor,
		completer,
		options...,
	)
	p.Run()
}

// executor executes command and print the output.
func executor(in string) {
	promptArgs := strings.Split(in, " ")
	os.Args = append([]string{os.Args[0]}, promptArgs...)
	lRootCmd.Execute()
}

// completer returns the completion items from user input.
func completer(d prompt.Document) []prompt.Suggest {
	currentCommand := lRootCmd
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

	var suggestions []prompt.Suggest
	if strings.HasPrefix(d.GetWordBeforeCursor(), "-") {
		flags := []*pflag.FlagSet{currentCommand.LocalNonPersistentFlags(), lRootCmd.PersistentFlags()}
		for _, fs := range flags {
			fs.VisitAll(func(flag *pflag.Flag) {
				suggestions = append(suggestions, prompt.Suggest{Text: "--" + flag.Name, Description: flag.Usage})
			})
		}
	} else if currentCommand.HasAvailableSubCommands() {
		for _, c := range currentCommand.Commands() {
			suggestions = append(suggestions, prompt.Suggest{Text: c.Name(), Description: c.Short})
		}
	} else if currentCommand.Annotations["coprompt"] != "" {
		copromptAnnotation := currentCommand.Annotations["coprompt"]
		suggestions = lHandleDynamicSuggestions(copromptAnnotation, d)
	}

	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}
