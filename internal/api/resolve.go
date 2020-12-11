package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gojektech/heimdall/v6/httpclient"
)

// Resolve fetches the public key for a nickname
func Resolve(url, nickname string) (string, error) {
	client := httpclient.NewClient()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", url, nickname), nil)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != 200 {
		return "", errors.New(string(body))
	}

	return string(body), nil
}
