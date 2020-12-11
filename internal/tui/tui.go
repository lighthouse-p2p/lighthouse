package tui

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lighthouse-p2p/lighthouse/internal/api"
	"github.com/lighthouse-p2p/lighthouse/internal/models"
	"github.com/logrusorgru/aurora"
	"github.com/tj/go-spin"
	"golang.org/x/crypto/nacl/sign"
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
		validationRegex := regexp.MustCompile("^[a-z]+$")

		nickname := ""
		survey.AskOne(
			&survey.Input{
				Message: "Nickname:",
			},
			&nickname,
			survey.WithValidator(survey.Required),
			survey.WithValidator(func(val interface{}) error {
				if !validationRegex.Match([]byte(val.(string))) {
					return errors.New("Nickname must be a-z in lower case")
				}

				return nil
			}),
		)

		done := make(chan bool)
		go Spinner(done, "Generating a keypair", "Done")

		publicKey, privateKey, err := sign.GenerateKey(rand.Reader)
		if err != nil {
			done <- true

			fmt.Printf("%s\n", aurora.Bold(aurora.Red("Cannot generate keypair ✕")))
			fmt.Printf("Error: %s\n", err)

			os.Exit(1)
		}

		publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKey[:])
		privateKeyBase64 := base64.StdEncoding.EncodeToString(privateKey[:])

		time.Sleep(1 * time.Second)
		done <- true

		done = make(chan bool)
		go Spinner(done, "Registering", "Registered")
		err = api.Register("http://localhost:3000/v1/register", publicKeyBase64, nickname)

		if err != nil {
			done <- true
			// time.Sleep(64 * time.Millisecond)

			fmt.Printf("\r  %s\n", aurora.Bold(aurora.Red("Registration failed ✕")))
			fmt.Printf("  %s %s\n", aurora.Bold(aurora.Red("Error:")), err)

			os.Exit(1)
		}

		time.Sleep(1 * time.Second)
		done <- true

		time.Sleep(64 * time.Millisecond)

		metadata, err := json.Marshal(models.Metadata{
			PubKey:   publicKeyBase64,
			PrivKey:  privateKeyBase64,
			NickName: nickname,
		})
		if err != nil {
			log.Fatalf("%s\n", err)
		}

		f, err := os.Create("metadata.json")
		if err != nil {
			log.Fatalf("%s\n", err)
		}

		_, err = f.WriteString(string(metadata))
		if err != nil {
			log.Fatalf("%s\n", err)
		}
	} else {
		return
	}
}

// Spinner creates a terminal loading prompt
func Spinner(done chan bool, loading, loaded string) {
	go func() {
		s := spin.New()
		s.Set(spin.Default)

	loop:
		for {
			select {
			case <-done:
				fmt.Printf("\r  %s %s", aurora.Bold(aurora.Green(fmt.Sprintf("%s ✓", loaded))), strings.Repeat(" ", Abs(len(loading)-len(loaded))))
				break loop
			default:
				fmt.Printf("\r  %s %s ", aurora.Bold(aurora.Cyan(loading)), s.Next())
				time.Sleep(64 * time.Millisecond)
			}
		}

		time.Sleep(64 * time.Millisecond)
	}()
}

// Abs gives the absolute value of the int
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
