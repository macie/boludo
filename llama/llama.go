// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

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

// SetDefault sets default Server.
func SetDefault(server Server) {
	defaultServer = server
	defaultClient = Client{Addr: defaultServer.Addr}
}

// Serve starts LLM server and returns error if it fails. It is the caller's
// responsibility to close Server.
func Serve(ctx context.Context, modelPath string) error {
	return defaultServer.Start(ctx, modelPath)
}

// Complete returns a channel with completion results for given string.
func Complete(ctx context.Context, s string) (chan string, error) {
	return defaultClient.Complete(ctx, s)
}

// Close releases all resources used by LLM server.
func Close() error {
	return defaultServer.Close()
}

// Options represent parameters for interacting with LLaMA model.
type Options struct {
	ModelPath string
	Format    string
	Temp      float32
	MinP      float32
	Seed      uint
}
