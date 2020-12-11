package rtc

import (
	"log"

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
func (s *Session) Init(pubKey string, st state.State) error {
	s.RemotePeer.PubKey = pubKey
	s.State = st

	// TODO: Resolve the pubKey

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

	answer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		return err
	}

	peerConnection.SetLocalDescription(answer)

	answerNew := *peerConnection.LocalDescription()

	answerNew.SDP = utils.StripSDP(answerNew.SDP)
	resp, err := utils.Encode(answerNew)
	if err != nil {
		return err
	}

	log.Println(answerNew)
	log.Println(resp)

	return nil
}
