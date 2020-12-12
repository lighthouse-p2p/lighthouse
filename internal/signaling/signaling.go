package signaling

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/lighthouse-p2p/lighthouse/internal/models"
	"github.com/lighthouse-p2p/lighthouse/internal/utils"
	"golang.org/x/crypto/nacl/sign"
)

// Client connects to the server websockets and authenticates
type Client struct {
	Metadata   models.Metadata
	Addr       string
	Connection *websocket.Conn
	Chans      map[string]chan string
	SignalChan chan models.Signal
}

// Init will authenticate with the signalling server
func (c *Client) Init(metadata models.Metadata) error {
	addr := utils.TranslateURL(fmt.Sprintf("ws://%s/v1/ws/signaling?pub_key=%s", metadata.Host, metadata.PubKey))

	c.Metadata = metadata
	c.Chans = make(map[string]chan string)

	connection, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return err
	}

	c.Connection = connection

	privateKeyRaw, err := base64.StdEncoding.DecodeString(metadata.PrivKey)
	if err != nil {
		return err
	}

	mt, message, err := connection.ReadMessage()
	if err != nil {
		if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
			return errors.New("The public key is not registered with the server")
		}

		return err
	}

	if mt != 1 {
		return errors.New("Unexpected socket message")
	}

	var privateKey [64]byte
	copy(privateKey[:], privateKeyRaw)

	signature := sign.Sign(nil, message, &privateKey)
	connection.WriteMessage(2, signature)

	mt, message, err = connection.ReadMessage()
	if err != nil {
		if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
			return errors.New("The signature didn't match the public key")
		}

		return err
	}

	if mt != 1 {
		return errors.New("Unexpected socket message")
	}

	if string(message) != "OK" {
		return errors.New("The signature didn't match the public key")
	}

	c.SignalChan = make(chan models.Signal, 32)

	return nil
}

// Listen gets and processes all the messages
func (c *Client) Listen() {
	go func() {
		connection := c.Connection

		for {
			_, message, err := connection.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
					break
				}

				continue
			}

			var signal models.Signal
			err = json.Unmarshal(message, &signal)
			if err != nil || signal.From == "" {
				continue
			}

			if signal.Type == "a" {
				if _, ok := c.Chans[signal.From]; !ok {
					c.Chans[signal.From] = make(chan string, 5)
				}

				log.Println("Got answer")

				c.Chans[signal.From] <- signal.SDP
			} else if signal.Type == "o" {
				log.Println("Got offer")
				c.SignalChan <- signal
			} else {
				continue
			}
		}
	}()
}

// Push sends a message on the socket
func (c *Client) Push(msg string) error {
	log.Println("Pushed")
	return c.Connection.WriteMessage(1, []byte(msg))
}

// Close disconnects cleanly
func (c *Client) Close() {
	c.Connection.Close()
}
