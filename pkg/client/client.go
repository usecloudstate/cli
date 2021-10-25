package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	hc    *http.Client
	token string
}

var (
	rootApp = "https://api.usecloudstate.io"
)

func Init() *Client {
	return &Client{
		hc: &http.Client{},
	}
}

func (client *Client) Request(method string, path string, body interface{}) (*http.Response, error) {
	postBody, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	requestBody := bytes.NewBuffer(postBody)

	url := fmt.Sprintf("%s/%s", rootApp, path)

	req, err := http.NewRequest(method, url, requestBody)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.token))
	req.Header.Set("Content-Type", "application/json")

	return client.hc.Do(req)
}

func (client *Client) SetAuthToken(token string) {
	client.token = token
}
