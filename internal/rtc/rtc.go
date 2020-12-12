package rtc

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/lighthouse-p2p/lighthouse/internal/api"
	"github.com/lighthouse-p2p/lighthouse/internal/models"
	"github.com/lighthouse-p2p/lighthouse/internal/rtc/wrapper"
	"github.com/lighthouse-p2p/lighthouse/internal/state"
	"github.com/lighthouse-p2p/lighthouse/internal/utils"
	"github.com/pion/webrtc/v2"
	"github.com/xtaci/smux"
)

// Sessions holds a map for open RTC sessions and the assigned ports
type Sessions struct {
	RTCSessions map[string]*Session
	PortMap     map[string]int
}

// Init initializes the session maps with empty values
func (s *Sessions) Init() {
	s.PortMap = make(map[string]int)
	s.RTCSessions = make(map[string]*Session)
}

// Session holds all the session information for a WebRTC session
type Session struct {
	RemotePeer struct {
		PubKey   string
		NickName string
	}

	Listener *net.Listener
	*webrtc.PeerConnection
}

// Init initializes a WebRTC session as an offer
func (s *Session) Init(nickname string, st state.State, port int) error {
	pubKey, err := api.Resolve(fmt.Sprintf("http://%s/v1/resolve", st.Metadata.Host), nickname)
	if err != nil {
		return err
	}

	s.RemotePeer.PubKey = pubKey
	s.RemotePeer.NickName = nickname

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
		Type: "o",
	}
	jsonSignal, err := json.Marshal(sig)
	if err != nil {
		return err
	}

	st.SignalingClient.Push(string(jsonSignal))

	for {
		if _, ok := st.SignalingClient.Chans[pubKey]; ok {
			remoteSDP := <-st.SignalingClient.Chans[pubKey]

			var remoteDesc webrtc.SessionDescription
			err = utils.Decode(remoteSDP, &remoteDesc)
			if err != nil {
				continue
			}

			err = peerConnection.SetRemoteDescription(remoteDesc)
			if err != nil {
				return err
			}

			dataChannel, err := peerConnection.CreateDataChannel("data", nil)
			if err != nil {
				return err
			}

			dataChannel.OnOpen(func() {
				go func() {
					log.Println("Connection seems to be up, preparing a multiplexed session")

					proxySrv, err := net.Listen("tcp4", fmt.Sprintf("127.0.0.1:%d", port))
					if err != nil {
						panic(err)
					}

					s.Listener = &proxySrv

					conn, err := wrapper.WrapConn(dataChannel, &wrapper.NilAddr{}, &wrapper.NilAddr{})
					if err != nil {
						dataChannel.Close()
						peerConnection.Close()
					}

					session, err := smux.Client(conn, nil)
					if err != nil {
						panic(err)
					}

					for {
						l, err := proxySrv.Accept()
						if err != nil {
							panic(err)
						}

						go func(c net.Conn) {
							stream, err := session.OpenStream()
							if err != nil {
								panic(err)
							}

							log.Printf("New stream, %d total streams!\n", session.NumStreams())

							go wrapper.JoinStreams(stream, c)
						}(l)
					}
				}()
			})

			break
		} else {
			time.Sleep(250 * time.Millisecond)
			continue
		}
	}

	s.PeerConnection = peerConnection

	return nil
}

// InitAnswer initializes a WebRTC session as an answer
func (s *Session) InitAnswer(signal models.Signal, push func(string)) error {
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

	var remoteDesc webrtc.SessionDescription
	err = utils.Decode(signal.SDP, &remoteDesc)
	if err != nil {
		return err
	}

	err = peerConnection.SetRemoteDescription(remoteDesc)
	if err != nil {
		return err
	}

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		return err
	}

	peerConnection.SetLocalDescription(answer)

	localSDP := *peerConnection.LocalDescription()

	localSDP.SDP = utils.StripSDP(localSDP.SDP)
	encodedLocalSDP, err := utils.Encode(localSDP)
	if err != nil {
		return err
	}

	sig := models.Signal{
		To:   signal.From,
		From: signal.To,
		SDP:  encodedLocalSDP,
		Type: "a",
	}
	jsonSignal, err := json.Marshal(sig)
	if err != nil {
		return err
	}

	push(string(jsonSignal))

	s.PeerConnection = peerConnection

	// setup data channel
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		// Register channel opening handling
		d.OnOpen(func() {
			log.Println("Connection seems to be up, preparing a multiplexed session")

			conn, err := wrapper.WrapConn(d, &wrapper.NilAddr{}, &wrapper.NilAddr{})
			if err != nil {
				d.Close()
				peerConnection.Close()
				panic(err)
			}

			session, err := smux.Server(conn, nil)
			if err != nil {
				d.Close()
				peerConnection.Close()
				panic(err)
			}

			for {
				stream, err := session.AcceptStream()
				if err != nil {
					session.Close()
					d.Close()
					peerConnection.Close()
					panic(err)
				}

				go func(stream *smux.Stream) {
					proxyConn, err := net.Dial("tcp4", "localhost:42011")
					if err != nil {
						panic(err)
					}

					wrapper.JoinStreams(stream, proxyConn)
				}(stream)
			}
		})
	})

	return nil
}
