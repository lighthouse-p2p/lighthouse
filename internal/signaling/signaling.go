package signaling

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/lighthouse-p2p/lighthouse/internal/models"
	"golang.org/x/crypto/nacl/sign"
)

// Client connects to the server websockets and authenticates
type Client struct {
	Metadata   models.Metadata
	Addr       string
	Connection *websocket.Conn
}

// Init will authenticate with the signalling server
func (c *Client) Init(metadata models.Metadata) error {
	addr := fmt.Sprintf("ws://%s/v1/ws/signaling?pub_key=%s", metadata.Host, metadata.PubKey)

	c.Metadata = metadata

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

	return nil
}
