package llama

import (
	"context"
	"testing"
)

func TestServerStart(t *testing.T) {
	const testingModel = "../external/TinyLLama-v0.Q8_0.gguf"
	server := Server{Path: "../llm-server"}
	defer server.Close()

	if err := server.Start(context.TODO(), testingModel); err != nil {
		t.Fatalf("Start(ctx) returns error: %v", err)
	}
}
