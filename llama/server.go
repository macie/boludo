// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package llama

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"
)

// Server represents LLM server.
type Server struct {
	cmd  *exec.Cmd
	log  *CmdLogger
	host string
	port int
}

// NewServer configures a new LLM server. It uses `llm-server` command which is
// `server` command from `llama.cpp` library.
//
// The server is not started by default.
func NewServer(ctx context.Context, serverPath string, options Options) Server {
	s := Server{
		host: "127.0.0.1",
		port: 24114,
		log:  NewCmdLogger(os.Stderr, "[llm-server] ", 0),
	}

	cmd := exec.CommandContext(ctx,
		serverPath,
		"--host", s.host,
		"--port", fmt.Sprint(s.port),
		"--model", options.ModelPath,
		"--threads", fmt.Sprint(runtime.NumCPU()),
		"--ctx-size", fmt.Sprint(2048),
	)
	cmd.Stdout = s.log
	cmd.Stderr = s.log
	s.cmd = cmd

	return s
}

// Start starts LLM server.
//
// It is the caller's responsibility to close Server.
func (s *Server) Start() error {
	err := s.cmd.Start()
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("cannot start a LLM server: command '%s' not found in directory '%s': %w", path.Base(s.cmd.Path), path.Dir(s.cmd.Path), os.ErrNotExist)
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
	addr := net.JoinHostPort(s.host, fmt.Sprint(s.port))
	conn, err := net.DialTimeout("tcp", addr, timeoutDNS)
	if err != nil || conn == nil {
		return false
	}

	conn.Close()
	return true
}

// Close frees all resources associated with server.
func (s *Server) Close() error {
	// FIXME: to release resources should be called: s.Wait() ???
	if s.cmd.Process == nil {
		// server is not running
		return nil
	}

	// FIXME(windows): doesn't work on Windows. See: https://pkg.go.dev/os#Process.Signal
	return s.cmd.Process.Signal(os.Interrupt)
}
