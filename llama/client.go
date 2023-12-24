// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package llama

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// completionRequest represents completion request to LLM server.
//
// See: https://github.com/ggerganov/llama.cpp/blob/master/examples/server/README.md#api-endpoints
type completionRequest struct {
	Prompt          string  `json:"prompt"`
	Temp            float32 `json:"temperature"`
	TopK            int     `json:"top_k"`
	MinP            float32 `json:"min_p"`
	TopP            float32 `json:"top_p"`
	Seed            int     `json:"seed"`
	WithoutNewlines bool    `json:"penalize_nl"`
	PredictNum      int     `json:"n_predict"`
	Streaming       bool    `json:"stream"`
}

// completionResponse represents completion response from LLM server.
//
// See: https://github.com/ggerganov/llama.cpp/blob/master/examples/server/README.md#api-endpoints
type completionResponse struct {
	Content string `json:"content"`
	Stop    bool   `json:"stop"`
}

// Client represents client for LLM server.
type Client struct {
	// URL specifies the address of the underlying LLM server.
	// If empty, the default URL is used: `http://localhost:24114`.
	URL string
}

// Complete returns a channel with completion results for given string.
func (c *Client) Complete(ctx context.Context, s string) (chan string, error) {
	req := completionRequest{
		Prompt:          s,
		Temp:            1.2,
		TopK:            40,
		MinP:            0.05,
		TopP:            0.0,
		Seed:            -1,
		PredictNum:      150,
		Streaming:       true,
		WithoutNewlines: false,
	}
	ch, err := c.infer(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("could not complete: %w", err)
	}
	return ch, nil
}

// infer is a low-level function for sending completion requests to the LLM server.
func (c *Client) infer(ctx context.Context, req completionRequest) (chan string, error) {
	url := c.URL
	if c.URL == "" {
		url = "http://localhost:24114"
	}
	url = fmt.Sprintf("%s/completion", url)

	// FIXME: use context
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("completion request cannot be serialized: %w", err)
	}
	resp, err := http.Post(url, "text/event-stream", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("completion request cannot be sent: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("LLM server returned error: %s", resp.Status)
	}

	ch := make(chan string)
	go func(respBody io.ReadCloser) {
		defer close(ch)
		defer respBody.Close()

		scanner := bufio.NewScanner(respBody)
		for scanner.Scan() {
			line := scanner.Bytes()

			var response completionResponse
			if bytes.HasPrefix(line, []byte("data: ")) {
				json.Unmarshal(line[6:], &response)
			} else {
				json.Unmarshal(line, &response)
			}

			if response.Stop {
				return
			}

			if response.Content == "" {
				continue
			}

			ch <- response.Content
		}
	}(resp.Body)

	return ch, nil
}
