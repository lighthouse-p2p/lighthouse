package rtc

import (
	"encoding/json"
	"fmt"

	"github.com/lighthouse-p2p/lighthouse/internal/api"
	"github.com/lighthouse-p2p/lighthouse/internal/models"
	"github.com/lighthouse-p2p/lighthouse/internal/state"
	"github.com/lighthouse-p2p/lighthouse/internal/utils"
	"github.com/pion/webrtc/v2"
)

// Session holds all the session information for a WebRTC session
type Session struct {
	RemotePeer struct {
		PubKey   string
		NickName string
	}

	State state.State
}

// Init initializes a WebRTC session
func (s *Session) Init(nickname string, st state.State) error {
	pubKey, err := api.Resolve(fmt.Sprintf("http://%s/v1/resolve", st.Metadata.Host), nickname)
	if err != nil {
		return err
	}

	s.RemotePeer.PubKey = pubKey
	s.State = st

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return err
	}

	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		return err
	}

	peerConnection.SetLocalDescription(offer)

	localSDP := *peerConnection.LocalDescription()

	localSDP.SDP = utils.StripSDP(localSDP.SDP)
	encodedLocalSDP, err := utils.Encode(localSDP)
	if err != nil {
		return err
	}

	sig := models.Signal{
		To:   pubKey,
		From: st.Metadata.PubKey,
		SDP:  encodedLocalSDP,
	}
	jsonSignal, err := json.Marshal(sig)
	if err != nil {
		return err
	}

	st.SignalingClient.Push(string(jsonSignal))

	return nil
}
