package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/lighthouse-p2p/lighthouse/internal/models"
	"github.com/lighthouse-p2p/lighthouse/internal/tui"
)

func main() {
	tui.GenerateASCIIArt()

	if _, err := os.Stat("metadata.json"); os.IsNotExist(err) {
		tui.StartNewUserFlow()
	} else {
		data, err := ioutil.ReadFile("metadata.json")
		if err != nil {
			log.Fatalf("%s\n", err)
		}

		var metadata models.Metadata
		json.Unmarshal(data, &metadata)

		if metadata.Host == "" ||
			metadata.NickName == "" ||
			metadata.PrivKey == "" ||
			metadata.PubKey == "" {
			log.Fatalln("Invalid metadata")
		}

		tui.AlreadyRegisteredFlow(metadata)
	}
}
