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

		if _, err := os.Stat("data/"); os.IsNotExist(err) {
			log.Fatalln("data/ not found\nThis directory is served when a new peer connection is made\nPlease make the data/ directory to connect")
		}

		go tui.AlreadyRegisteredFlow(metadata)

		// Block forever
		select {}
	}
}
