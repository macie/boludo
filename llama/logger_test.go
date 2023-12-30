// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package llama

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"testing"
)

type TestHandler struct{ Output io.Writer }

func (TestHandler) Enabled(_ context.Context, _ slog.Level) bool { return true }
func (TestHandler) WithGroup(string) slog.Handler                { return nil }
func (TestHandler) WithAttrs([]slog.Attr) slog.Handler           { return nil }
func (h TestHandler) Handle(_ context.Context, r slog.Record) error {
	_, err := h.Output.Write([]byte(r.Level.String() + " " + r.Message + "\n"))
	return err
}

func TestCmdLogger(t *testing.T) {
	testcases := []struct {
		pattern string
		want    string
	}{
		{"foo", "INFO foo\n"},
		{"foo\n{\"level\":\"WARNING\", \"message\":\"bar\"}\n\n", "INFO foo\nWARN {\"level\":\"WARNING\", \"message\":\"bar\"}\n"},
		{"foo\n{\"level\":\"ERROR\", \"message\":\"bar\"}\nbaz\n", "INFO foo\nERROR {\"level\":\"ERROR\", \"message\":\"bar\"}\nINFO baz\n"},
		{"", ""},
		{".", ""},
		{"foo\n...\n", "INFO foo\n"},
		{"foo\n..bar", "INFO foo\nINFO ..bar\n"},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.pattern, func(t *testing.T) {
			t.Parallel()
			output := new(bytes.Buffer)
			c := CmdLogger{Log: slog.New(TestHandler{output})}
			if _, err := c.Write([]byte(tc.pattern)); err != nil {
				t.Fatalf("Write([]byte(%q)) returns error: %v", tc.pattern, err)
			}
			got := output.String()
			if got != tc.want {
				t.Fatalf("Write([]byte(%q)) = %q, want %q", tc.pattern, got, tc.want)
			}
		})
	}
}
