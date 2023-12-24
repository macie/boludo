// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package llama

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"
)

func setupTestServer(t *testing.T) Server {
	const testingModel = "./external/TinyLLama-v0.Q8_0.gguf"

	_, filename, _, _ := runtime.Caller(0)
	if err := os.Chdir(path.Join(path.Dir(filename), "../")); err != nil {
		t.Fatalf(fmt.Sprintf("cannot setup test server: %v", err))
	}

	return Server{
		ModelPath: testingModel,
		Addr:      "localhost:24114",
	}
}

func TestServerStart(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	if err := server.Start(context.TODO()); err != nil {
		t.Fatalf("Start(ctx) returns error: %v", err)
	}
}
