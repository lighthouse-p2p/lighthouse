package main

import (
	"os"

	"github.com/lighthouse-p2p/lighthouse/internal/tui"
)

func main() {
	tui.GenerateASCIIArt()

	if _, err := os.Stat("metadata.json"); os.IsNotExist(err) {
		tui.StartNewUserFlow()
	}
}
