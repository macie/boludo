// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package boludo

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// UnstructuredHandler is a slog.Handler that writes all log records to stderr
// without structure.
type UnstructuredHandler struct {
	// Prefix specifies an optional prefix for each log record.
	Prefix string

	// Level specifies a minimum level for log records.
	Level slog.Level

	// Output specifies a destination for log records.
	// If nil, os.Stderr is used.
	Output io.Writer
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
	output := u.Output
	if output == nil {
		output = os.Stderr
	}

	if u.Prefix != "" {
		output.Write([]byte(u.Prefix))
		output.Write([]byte{' '})
	}
	output.Write([]byte(r.Level.String()))
	output.Write([]byte{' '})
	output.Write([]byte(r.Message))
	r.Attrs(func(a slog.Attr) bool {
		output.Write([]byte{' '})
		output.Write([]byte(a.String()))
		return true
	})
	output.Write([]byte{'\n'})

	return nil
}
