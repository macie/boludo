// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package boludo

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

// UnstructuredHandler is a slog.Handler that writes all log records to stderr
// without structure.
type UnstructuredHandler struct {
	// Prefix specifies an optional prefix for each log record.
	Prefix string

	// Level specifies a minimum level for log records.
	Level slog.Level
}

// WithAttrs implements slog.Handler.
func (UnstructuredHandler) WithAttrs([]slog.Attr) slog.Handler {
	return nil
}

// WithGroup implements slog.Handler.
func (UnstructuredHandler) WithGroup(string) slog.Handler { return nil }

// Enabled filters log records based on their level.
func (u UnstructuredHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= u.Level
}

// Handle writes log records to stderr.
func (u UnstructuredHandler) Handle(_ context.Context, r slog.Record) error {
	line := strings.Builder{}

	if u.Prefix != "" {
		line.WriteString(u.Prefix)
		line.WriteString(" ")
	}
	line.WriteString(r.Level.String())
	line.WriteString(" ")
	line.WriteString(r.Message)

	r.Attrs(func(a slog.Attr) bool {
		line.WriteString(" ")
		line.WriteString(a.String())
		return true
	})

	fmt.Println(line.String())

	return nil
}
