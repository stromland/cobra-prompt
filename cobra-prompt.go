package cobraprompt

import (
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// DynamicSuggestionsAnnotation for dynamic suggestions.
const DynamicSuggestionsAnnotation = "cobra-prompt-dynamic-suggestions"

// PersistFlagValuesFlag the flag that will be avaiailable when PersistFlagValues is true
const PersistFlagValuesFlag = "persist-flag-values"

type CobraPromptOptions struct {
	// GoPromptOptions is for customize go-prompt
	// see https://github.com/c-bata/go-prompt/blob/master/option.go
	GoPromptOptions []prompt.Option

	// FindSuggestionsOptions contains options available when traversing command and flag tree
	FindSuggestionsOptions FindSuggestionsOptions

	// PersistFlagValues will persist flags. For example have verbose turned on every command.
	PersistFlagValues bool

	// DisableCompletionCommand will disable the default completion command for cobra
	DisableCompletionCommand bool

	// AddDefaultExitCommand adds a command for exiting prompt loop
	AddDefaultExitCommand bool

	// OnErrorFunc handle error for command.Execute, if not set print error and exit
	OnErrorFunc func(err error)
}

type FindSuggestionsOptions struct {
	// DynamicSuggestionsFunc will be executed if an command has CallbackAnnotation as an annotation. If it's included
	// the value will be provided to the DynamicSuggestionsFunc function.
	DynamicSuggestionsFunc func(annotationValue string, document *prompt.Document) []prompt.Suggest

	// ShowHiddenCommands makes hidden commands available
	ShowHiddenCommands bool

	// ShowHiddenFlags makes hidden flags available
	ShowHiddenFlags bool

	// ShowHelpCommandAndFlags will make help command and flag for every command available.
	ShowHelpCommandAndFlags bool
}

// CobraPrompt given a Cobra command it will make every flag and sub commands available as suggestions.
// Command.Short will be used as description for the suggestion.
type CobraPrompt struct {
	// RootCmd is the start point, all its sub commands and flags will be available as suggestions
	rootCmd *cobra.Command

	options *CobraPromptOptions

	prompt *prompt.Prompt
}

func New(rootCmd cobra.Command, options CobraPromptOptions) *CobraPrompt {
	prompt := prompt.New(
		func(in string) {
			promptArgs := strings.Fields(in)
			os.Args = append([]string{os.Args[0]}, promptArgs...)
			if err := rootCmd.Execute(); err != nil {
				if options.OnErrorFunc != nil {
					options.OnErrorFunc(err)
				} else {
					rootCmd.PrintErrln(err)
					os.Exit(1)
				}
			}
		},
		func(d prompt.Document) []prompt.Suggest {
			return FindSuggestions(&rootCmd, &d, options.FindSuggestionsOptions)
		},
		options.GoPromptOptions...,
	)

	co := &CobraPrompt{
		rootCmd: &rootCmd,
		options: &options,
		prompt:  prompt,
	}

	co.prepare()

	return co
}

// Run will automatically generate suggestions for all cobra commands and flags defined by RootCmd
// and execute the selected commands. Run will also reset all given flags by default, see PersistFlagValues
func (co CobraPrompt) Run() {
	co.prompt.Run()
}

func (co CobraPrompt) prepare() {
	if co.options.FindSuggestionsOptions.ShowHelpCommandAndFlags {
		// TODO: Find help commands
		co.rootCmd.InitDefaultHelpCmd()
	}

	if co.options.DisableCompletionCommand {
		co.rootCmd.CompletionOptions.DisableDefaultCmd = true
	}

	if co.options.AddDefaultExitCommand {
		co.rootCmd.AddCommand(&cobra.Command{
			Use:   "exit",
			Short: "Exit prompt",
			Run: func(cmd *cobra.Command, args []string) {
				os.Exit(0)
			},
		})
	}

	if co.options.PersistFlagValues {
		co.rootCmd.PersistentFlags().BoolP(PersistFlagValuesFlag, "",
			false, "Persist last given value for flags")
	}
}

func FindSuggestions(rootCmd *cobra.Command, d *prompt.Document, options FindSuggestionsOptions) []prompt.Suggest {
	args := strings.Fields(d.CurrentLine())

	cmd := rootCmd
	if found, _, err := rootCmd.Find(args); err == nil {
		cmd = found
	}

	var suggestions []prompt.Suggest
	persistFlagValues, _ := cmd.Flags().GetBool(PersistFlagValuesFlag)
	addFlags := func(flag *pflag.Flag) {
		if flag.Changed && !persistFlagValues {
			flag.Value.Set(flag.DefValue)
		}
		if flag.Hidden && !options.ShowHiddenFlags {
			return
		}
		if strings.HasPrefix(d.GetWordBeforeCursor(), "--") {
			suggestions = append(suggestions, prompt.Suggest{Text: "--" + flag.Name, Description: flag.Usage})
		} else if strings.HasPrefix(d.GetWordBeforeCursor(), "-") && flag.Shorthand != "" {
			suggestions = append(suggestions, prompt.Suggest{Text: "-" + flag.Shorthand, Description: flag.Usage})
		}
	}

	cmd.LocalFlags().VisitAll(addFlags)
	cmd.InheritedFlags().VisitAll(addFlags)

	if cmd.HasAvailableSubCommands() {
		for _, c := range cmd.Commands() {
			if !c.Hidden && !options.ShowHiddenCommands {
				suggestions = append(suggestions, prompt.Suggest{Text: c.Name(), Description: c.Short})
			}
			if options.ShowHelpCommandAndFlags {
				c.InitDefaultHelpFlag()
			}
		}
	}

	annotation := cmd.Annotations[DynamicSuggestionsAnnotation]
	if options.DynamicSuggestionsFunc != nil && annotation != "" {
		suggestions = append(suggestions, options.DynamicSuggestionsFunc(annotation, d)...)
	}
	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}
