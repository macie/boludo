// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package llama

import (
	"strings"
	"testing"
)

func TestServerSystemInfo(t *testing.T) {
	server := NewServer()
	defer server.Close()

	got := server.SystemInfo()
	if !strings.HasPrefix(got, "AVX = ") {
		t.Fatalf("unexpeted result of SystemInfo() = \"%v\"", got)
	}
}
