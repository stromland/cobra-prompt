package cobraprompt

import (
	"testing"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestFindSuggestions(t *testing.T) {
	rootCmd := newTestCommand("root", "The root cmd")
	getCmd := newTestCommand("get", "Get something")
	getObjectCmd := newTestCommand("object", "Get the object")
	getThingCmd := newTestCommand("thing", "The thing")
	getFoodCmd := newTestCommand("food", "Get some food")
	getFoodCmd.PersistentFlags().StringP("name", "n", "John", "name of the person to get some food from")
	_ = getFoodCmd.RegisterFlagCompletionFunc("name", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"John", "Mary\tMarianne - John's Mother", "Anne"}, cobra.ShellCompDirectiveNoFileComp
	})

	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getObjectCmd, getThingCmd, getFoodCmd)
	getObjectCmd.Flags().BoolP("verbose", "v", false, "Verbose log")

	cp := &CobraPrompt{
		RootCmd: rootCmd,
	}

	tests := []struct {
		name            string
		input           string
		expectedResults []string
	}{
		{
			name:            "Root suggestions",
			input:           "",
			expectedResults: []string{"get"},
		},
		{
			name:            "Get command suggestions",
			input:           "get ",
			expectedResults: []string{"object", "food", "thing"},
		},
		{
			name:            "Verbose flag suggestions",
			input:           "get object -",
			expectedResults: []string{"-v"},
		},
		{
			name:            "Verbose long flag suggestions",
			input:           "get object --",
			expectedResults: []string{"--verbose"},
		},
		{
			name:            "Name flag suggestions after flag",
			input:           "get food --name ",
			expectedResults: []string{"John", "Mary", "Anne"},
		},
		{
			name:            "Name flag suggestions with partial value",
			input:           "get food --name J",
			expectedResults: []string{"John"},
		},
		{
			name:            "Shorthand name flag suggestions",
			input:           "get food -n ",
			expectedResults: []string{"John", "Mary", "Anne"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := prompt.NewBuffer()
			buf.InsertText(test.input, false, true)
			suggestions := cp.findSuggestions(*buf.Document())

			assert.Len(t, suggestions, len(test.expectedResults), "Incorrect number of suggestions")

			actualSuggestionsMap := make(map[string]struct{})
			for _, suggestion := range suggestions {
				actualSuggestionsMap[suggestion.Text] = struct{}{}
			}

			// Check each expected result is present in actual suggestions
			for _, expected := range test.expectedResults {
				_, exists := actualSuggestionsMap[expected]
				assert.True(t, exists, "Expected suggestion not found: "+expected)
			}
		})
	}
}

func newTestCommand(use string, short string) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Run:   func(cmd *cobra.Command, args []string) {},
	}
}
