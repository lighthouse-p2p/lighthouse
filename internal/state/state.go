package state

import (
	"github.com/lighthouse-p2p/lighthouse/internal/models"
)

// State holds the global application state
type State struct {
	Metadata models.Metadata
}
