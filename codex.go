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

type CodexSuggestion struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string      `json:"text"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type CodexRequest struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	Temperature      float64 `json:"temperature"`
	MaxTokens        int     `json:"max_tokens"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  int     `json:"presence_penalty"`
}

type OpenApiCodex interface {
	GetCodexSuggestion(prompt string) (suggestion CodexSuggestion, err error)
}

type CodexApi struct {
	apiKey string
}

func NewCodexApi(apiKey string) OpenApiCodex {
	return &CodexApi{
		apiKey: apiKey,
	}
}

func (c *CodexApi) GetCodexSuggestion(prompt string) (suggestion CodexSuggestion, err error) {
	timeout := 100 * time.Second
	client := httpclient.NewClient(httpclient.WithHTTPTimeout(timeout))

	payloadBody := CodexRequest{
		Model:            "text-davinci-003",
		Prompt:           prompt,
		Temperature:      0.7,
		MaxTokens:        500,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}

	payloadByte, err := json.Marshal(payloadBody)
	if err != nil {
		panic(err)
	}

	reqHeader := http.Header{}
	reqHeader.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	reqHeader.Set("Content-Type", "application/json")

	res, err := client.Post("https://api.openai.com/v1/completions", bytes.NewBuffer(payloadByte), reqHeader)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, &suggestion)

	return
}
