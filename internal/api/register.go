package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gojektech/heimdall/v6/httpclient"
	"github.com/lighthouse-p2p/lighthouse/internal/models"
)

// Register registers a public key and nickname that on the network
func Register(url, pubKey, nickName string) error {
	client := httpclient.NewClient()

	body, err := json.Marshal(&models.RegisterRequest{
		PubKey:   pubKey,
		NickName: nickName,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 201 {
		body, _ := ioutil.ReadAll(res.Body)

		return errors.New(string(body))
	}

	return nil
}
