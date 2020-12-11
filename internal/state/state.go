package state

import (
	"github.com/lighthouse-p2p/lighthouse/internal/models"
	"github.com/lighthouse-p2p/lighthouse/internal/signaling"
)

// State holds the global application state
type State struct {
	Metadata        models.Metadata
	SignalingClient *signaling.Client
}
