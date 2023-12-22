// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package llama

import (
	"bufio"
	"bytes"
	"io"
	"log"
)

// CmdLogger is a log.Logger wrapper for exec.Cmd output.
type CmdLogger struct {
	logger *log.Logger
}

// NewCmdLogger creates new CmdLogger.
func NewCmdLogger(out io.Writer, prefix string, flag int) *CmdLogger {
	return &CmdLogger{logger: log.New(out, prefix, flag)}
}

// Write implements io.Writer.
func (e CmdLogger) Write(p []byte) (n int, err error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))
	for scanner.Scan() {
		if scanner.Text() == "" || scanner.Text() == "." {
			continue
		}
		e.logger.Println(scanner.Text())
	}
	return len(p), nil
}
