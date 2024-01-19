package llama

import (
	"context"
)

var (
	defaultServer = Server{Addr: "localhost:24114"}
	defaultClient = Client{Addr: "localhost:24114"}
	// DefaultOptions represent neutral parameters for interacting with LLaMA model.
	DefaultOptions = Options{
		ModelPath: "",
		Seed:      0,
		Temp:      1,
		MinP:      0,
	}
)

// SetDefault sets default Client and Server.
func SetDefault(server Server, client Client) {
	defaultServer = server
	defaultClient = client
}

// Serve starts LLM server and returns error if it fails. It is the caller's
// responsibility to close Server.
func Serve(ctx context.Context, modelPath string) error {
	return defaultServer.Start(ctx, modelPath)
}

// Complete returns a channel with completion results for given string.
func Complete(ctx context.Context, p Prompt) (chan string, error) {
	return defaultClient.Complete(ctx, p)
}

// Close releases all resources used by LLM server.
func Close() error {
	return defaultServer.Close()
}

// Options represent parameters for interacting with LLaMA model.
type Options struct {
	ModelPath string
	Temp      float32
	MinP      float32
	Seed      uint
}

// Update updates the Options with the non-default values from other Options.
func (o *Options) Update(other Options) {
	if other.ModelPath != DefaultOptions.ModelPath {
		o.ModelPath = other.ModelPath
	}
	if other.Temp != DefaultOptions.Temp {
		o.Temp = other.Temp
	}
	if other.MinP != DefaultOptions.MinP {
		o.MinP = other.MinP
	}
	if other.Seed != DefaultOptions.Seed {
		o.Seed = other.Seed
	}
}
