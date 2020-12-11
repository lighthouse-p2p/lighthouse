package models

// Signal is used to define the signal msg transmitted over the network
type Signal struct {
	Type string `json:"type"`
	To   string `json:"to"`
	From string `json:"from"`
	SDP  string `json:"sdp"`
}
