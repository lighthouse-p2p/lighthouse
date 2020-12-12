package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gojektech/heimdall/v6/httpclient"
)

// Coins fetches the total coins
func Coins(url, pubKey string) (string, error) {
	client := httpclient.NewClient()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s?pub_key=%s", url, pubKey), nil)
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
