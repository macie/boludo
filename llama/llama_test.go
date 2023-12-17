// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package llama

import (
	"fmt"
	"strings"
	"testing"
)

const testingModelPath = "../external/TinyLLama-v0.Q8_0.gguf"

func TestServerSystemInfo(t *testing.T) {
	server := NewServer()
	defer server.Close()

	got := server.SystemInfo()
	if !strings.HasPrefix(got, "AVX = ") {
		t.Fatalf("unexpeted result of SystemInfo() = \"%v\"", got)
	}
}

func TestModelString(t *testing.T) {
	server := NewServer()
	defer server.Close()
	model, err := NewModel(Options{
		ModelPath: testingModelPath,
	})
	if err != nil {
		t.Fatalf("NewModel(options) returns error: %v", err)
	}
	defer model.Close()

	want := fmt.Sprintf("Model (path: %s; context size (trainied): 2048 (2048); vocabulary size: 32000; embeddings size: 64)", testingModelPath)
	got := model.String()
	if got != want {
		t.Fatalf("String() = %v; want %v", got, want)
	}
}

func TestModelMaxContextSize(t *testing.T) {
	server := NewServer()
	defer server.Close()
	model, err := NewModel(Options{
		ModelPath: testingModelPath,
	})
	if err != nil {
		t.Fatalf("NewModel(options) returns error: %v", err)
	}
	defer model.Close()

	want := 2048
	got := model.MaxContextSize()
	if got != want {
		t.Fatalf("MaxContextSize() = %v; want %v", got, want)
	}
}
