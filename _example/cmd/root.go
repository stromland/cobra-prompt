package cmd

import (
	"github.com/spf13/cobra"
	cobraprompt "github.com/stromland/cobra-prompt"
)

const CallbackAnnotation = cobraprompt.CallbackAnnotation

var RootCmd = &cobra.Command{}
