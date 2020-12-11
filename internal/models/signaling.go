package models

// Signal is used to define the signal msg transmitted over the network
type Signal struct {
	To  string `json:"to"`
	SDP string `json:"sdp"`
}
