// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package llama

import (
	"context"
	"testing"
)

const testingModel = "../external/TinyLLama-v0.Q8_0.gguf"
const serverPath = "../llm-server"

func TestServerStart(t *testing.T) {
	server := NewServer(context.TODO(),
		serverPath, Options{ModelPath: testingModel},
	)
	defer server.Close()

	if err := server.Start(); err != nil {
		t.Fatalf("Start(ctx, options) returns error: %v", err)
	}
}
