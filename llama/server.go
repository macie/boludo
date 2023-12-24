// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package llama

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"
)

// Server represents LLM server.
type Server struct {
	// ModelPath specifies the path to the model file in GGUF format.
	ModelPath string

	// Addr optionally specifies the TCP address for the server to listen on,
	// in the form "host:port".
	// If empty, "localhost:24114"is used.
	Addr string

	// Cmd specifies a command for underlying LLM server.
	// If nil, the default command is used: `./llm-server --ctx-size 2048`.
	Cmd *exec.Cmd

	// Logger specifies an optional logger for underlying server errors and
	// debug messages.
	// If nil, logging is done to stderr.
	Logger *log.Logger
}

// Start starts LLM server.
//
// It is the caller's responsibility to close Server.
func (s *Server) Start(ctx context.Context) error {
	f, err := os.Stat(s.ModelPath)
	if err != nil {
		return fmt.Errorf("cannot start a LLM server: model '%s' not found: %w", s.ModelPath, err)
	}
	if f.IsDir() {
		return fmt.Errorf("cannot start a LLM server: model path '%s' is a directory", s.ModelPath)
	}

	if s.Logger == nil {
		s.Logger = log.New(os.Stderr, "[llm-server] ", 0)
	}

	if s.Addr == "" {
		s.Addr = "localhost:24114"
	}
	host, port, err := net.SplitHostPort(s.Addr)
	if err != nil {
		return fmt.Errorf("cannot start a LLM server: invalid address %s: %w", s.Addr, err)
	}

	if s.Cmd == nil {
		cmdLogger := CmdLogger{
			ErrorLog: s.Logger,
		}
		s.Cmd = exec.CommandContext(ctx,
			"./llm-server",
			"--host", host,
			"--port", port,
			"--model", s.ModelPath,
			"--threads", fmt.Sprint(runtime.NumCPU()),
			"--ctx-size", fmt.Sprint(2048),
		)
		s.Cmd.Stdout = cmdLogger
		s.Cmd.Stderr = cmdLogger
	}

	if err := s.Cmd.Start(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("cannot start a LLM server: command '%s' not found in directory '%s': %w", path.Base(s.Cmd.Path), path.Dir(s.Cmd.Path), os.ErrNotExist)
		}

		return fmt.Errorf("cannot start a LLM server: %w", err)
	}

	// wait for server to start. Check frequency is limited
	i := 1
	for {
		if ok := s.Ping(); ok {
			break
		}
		time.Sleep(time.Duration(i*25) * time.Millisecond)
		if i < 8 { // maximum wait time is 0.2 seconds
			i = i * 2
		}
	}

	return err
}

// Ping checks if server is running.
func (s *Server) Ping() bool {
	// 130 ms is an average DNS lookup time observed by Googlebot
	// See: https://developers.google.com/speed/public-dns/docs/performance#cache_misses
	timeoutDNS := 130 * time.Millisecond
	conn, err := net.DialTimeout("tcp", s.Addr, timeoutDNS)
	if err != nil || conn == nil {
		return false
	}

	conn.Close()
	return true
}

// Close frees all resources associated with server.
func (s *Server) Close() error {
	// FIXME: to release resources should be called: s.Wait() ???
	if s.Cmd == nil || s.Cmd.Process == nil {
		// server is not running
		return nil
	}

	// FIXME(windows): doesn't work on Windows. See: https://pkg.go.dev/os#Process.Signal
	return s.Cmd.Process.Signal(os.Interrupt)
}