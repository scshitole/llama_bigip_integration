package llama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	apiURL string
}

type CompletionRequest struct {
	Prompt string `json:"prompt"`
}

type CompletionResponse struct {
	Text string `json:"text"`
}

func NewClient(apiURL string) (*Client, error) {
	return &Client{apiURL: apiURL}, nil
}

func (c *Client) GetCompletion(prompt string) (string, error) {
	reqBody, err := json.Marshal(CompletionRequest{Prompt: prompt})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	resp, err := http.Post(c.apiURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to send request to LLaMA API: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLaMA API returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}

	var completionResp CompletionResponse
	err = json.Unmarshal(body, &completionResp)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return completionResp.Text, nil
}
