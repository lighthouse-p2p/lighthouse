package tui

import (
	"github.com/AlecAivazis/survey/v2"
)

// NewUserFlow is the TUI flow used when lighthouse metadata is missing
var NewUserFlow = []*survey.Question{
	{
		Name: "register",
		Prompt: &survey.Confirm{
			Message: "metadata.json was not found in the current directory, do you want to register and create one now?",
		},
	},
}

// NewUserAnswers contains the answers from NewUserFlow
type NewUserAnswers struct {
	Register bool
}
