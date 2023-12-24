// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package llama

import (
	"bufio"
	"bytes"
	"log"
)

// CmdLogger is a log.Logger wrapper for exec.Cmd output.
type CmdLogger struct {
	// Log specifies an optional logger for exec.Cmd output
	// If nil, logging is done via the log package's standard logger.
	ErrorLog *log.Logger
}

// Write implements io.Writer.
func (c CmdLogger) Write(p []byte) (n int, err error) {
	if c.ErrorLog == nil {
		c.ErrorLog = log.Default()
	}

	scanner := bufio.NewScanner(bytes.NewReader(p))
	for scanner.Scan() {
		if scanner.Text() == "" || scanner.Text() == "." {
			continue
		}
		c.ErrorLog.Println(scanner.Text())
	}
	return len(p), nil
}
