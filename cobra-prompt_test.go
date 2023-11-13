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

	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getObjectCmd, getThingCmd)
	getObjectCmd.Flags().BoolP("verbose", "v", false, "Verbose log")

	cp := &CobraPrompt{
		RootCmd: rootCmd,
	}

	tests := []struct {
		name            string
		input           string
		expectedCount   int
		expectedResults []string
	}{
		{
			name:            "Root suggestions",
			input:           "",
			expectedCount:   1,
			expectedResults: []string{"get"},
		},
		{
			name:            "Get command suggestions",
			input:           "get ",
			expectedCount:   2,
			expectedResults: []string{"object", "thing"},
		},
		{
			name:            "Verbose flag suggestions",
			input:           "get object -",
			expectedCount:   1,
			expectedResults: []string{"-v"},
		},
		{
			name:            "Verbose long flag suggestions",
			input:           "get object --",
			expectedCount:   1,
			expectedResults: []string{"--verbose"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := prompt.NewBuffer()
			buf.InsertText(test.input, false, true)
			suggestions := cp.findSuggestions(*buf.Document())

			assert.Len(t, suggestions, test.expectedCount, "Incorrect number of suggestions")
			for i, expected := range test.expectedResults {
				assert.Equal(t, expected, suggestions[i].Text, "Incorrect suggestion")
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
