package llama

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/macie/boludo"
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
	RepeatPenalty   float32 `json:"repeat_penalty"`
	RepeatLastN     int     `json:"repeat_last_n"`
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
	// Addr specifies the address of the LLM server.
	// If empty, "localhost:24114" is used.
	Addr string

	// Options specifies options for underlying LLM server.
	// If nil, DefaultOptions are used.
	Options *Options

	// Logger specifies logger for the client.
	Logger *slog.Logger
}

// Complete returns a channel with completion results for given string.
func (c *Client) Complete(ctx context.Context, p Prompt) (chan string, error) {
	if c.Options == nil {
		c.Options = &DefaultOptions
	}
	req := completionRequest{
		Prompt:          p.String(),
		Temp:            c.Options.Temp, // 1.0 means disable
		TopK:            0,              // 0 means disable
		MinP:            c.Options.MinP, // 0.0 means disable
		TopP:            1.0,            // 1.0 means disable
		Seed:            5489,           // -1 means random seed
		PredictNum:      -1,             // -1 means infinite
		RepeatPenalty:   1.0,            // 1.0 means disable
		RepeatLastN:     0,              // 0 means disable
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
	if c.Logger == nil {
		c.Logger = slog.New(boludo.UnstructuredHandler{Prefix: "[llm-client]", Level: slog.LevelInfo})
	}
	if c.Addr == "" {
		c.Addr = "localhost:24114"
	}
	url := fmt.Sprintf("http://%s/completion", c.Addr)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("completion request cannot be serialized: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("completion request cannot be created: %w", err)
	}
	c.Logger.Info("completion request", slog.String("url", url), slog.String("body", string(reqBody)))
	resp, err := http.DefaultClient.Do(httpReq)
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
