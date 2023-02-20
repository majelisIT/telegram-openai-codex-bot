package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gojek/heimdall/v7/httpclient"
)

type DalleResponse struct {
	Created int `json:"created"`
	Data    []struct {
		URL string `json:"url"`
	} `json:"data"`
}

type DalleRequest struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type OpenApiDalle interface {
	GetDalleImage(prompt string) (res DalleResponse, err error)
}

type DalleApi struct {
	apiKey string
}

func NewDalleApi(apiKey string) OpenApiDalle {
	return &DalleApi{
		apiKey: apiKey,
	}
}

func (c *DalleApi) GetDalleImage(prompt string) (art DalleResponse, err error) {
	timeout := 100 * time.Second
	client := httpclient.NewClient(httpclient.WithHTTPTimeout(timeout))

	payloadBody := DalleRequest{
		Prompt: prompt,
		N:      1,
		Size:   "1024x1024",
	}

	payloadByte, err := json.Marshal(payloadBody)
	if err != nil {
		panic(err)
	}

	reqHeader := http.Header{}
	reqHeader.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	reqHeader.Set("Content-Type", "application/json")

	res, err := client.Post("https://api.openai.com/v1/images/generations", bytes.NewBuffer(payloadByte), reqHeader)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, &art)

	return
}
