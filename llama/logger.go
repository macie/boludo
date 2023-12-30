// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package llama

import (
	"bufio"
	"bytes"
	"log/slog"
	"regexp"
	"strings"
)

// CmdLogger is a slog.Logger wrapper for exec.Cmd output.
type CmdLogger struct {
	// Log specifies an optional logger for exec.Cmd output
	// If nil, logging is done via the slog package's standard logger.
	Log *slog.Logger
}

// Write implements io.Writer.
func (c *CmdLogger) Write(p []byte) (n int, err error) {
	if c.Log == nil {
		c.Log = slog.Default()
	}

	progressbarPattern := regexp.MustCompile(`^\.+$`)
	scanner := bufio.NewScanner(bytes.NewReader(p))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, ".") && progressbarPattern.MatchString(line) {
			continue
		}

		if strings.Contains(line, `"level":"ERROR"`) {
			c.Log.Error(line)
		} else if strings.Contains(line, `"level":"WARNING"`) {
			c.Log.Warn(line)
		} else {
			// "level":"INFO" and unstructured lines (see: server.cpp)
			c.Log.Info(line)
		}
	}
	return len(p), nil
}
