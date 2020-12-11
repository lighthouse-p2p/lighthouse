package tui

import (
	"fmt"
	"sync"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/logrusorgru/aurora"
	"github.com/tj/go-spin"
)

// NewUserFlow is the TUI flow used when lighthouse metadata is missing
var NewUserFlow = []*survey.Question{
	{
		Name: "register",
		Prompt: &survey.Confirm{
			Message: "metadata.json was not found in the current directory, do you want to register and create one now?",
		},
		Validate: survey.Required,
	},
}

// NewUserAnswers contains the answers from NewUserFlow
type NewUserAnswers struct {
	Register bool
}

// StartNewUserFlow starts the TUI flow and registers the user
func StartNewUserFlow() {
	answers := &NewUserAnswers{}
	survey.Ask(NewUserFlow, answers)

	if answers.Register {
		s := spin.New()
		s.Set(spin.Default)

		var wg sync.WaitGroup
		doneCh := make(chan struct{})

		wg.Add(1)
		go func() {
			time.Sleep(1 * time.Second)
			wg.Done()
		}()

		go func() {
			wg.Wait()
			close(doneCh)
		}()

	loop:
		for {
			select {
			case <-doneCh:
				fmt.Printf("\r  %s        ", aurora.Bold(aurora.Green("Registered ✓")))
				break loop
			default:
				fmt.Printf("\r  %s %s ", aurora.Bold(aurora.Cyan("Registering")), s.Next())
				time.Sleep(64 * time.Millisecond)
			}
		}
	} else {
		return
	}
}
