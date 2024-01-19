package llama

import (
	"context"
	"strings"
	"testing"
)

func TestComplete(t *testing.T) {
	const testingModel = "../external/TinyLLama-v0.Q8_0.gguf"
	server := Server{Path: "../llm-server"}
	if err := server.Start(context.TODO(), testingModel); err != nil {
		t.Fatalf("cannot start a LLM server: %v", err)
	}
	defer server.Close()

	prompt := Prompt{}
	prompt.Add("Once upon a time")

	client := Client{}
	c, err := client.Complete(context.TODO(), prompt)
	if err != nil {
		t.Fatalf("client.Complete() returns error: %v", err)
	}

	result := strings.Builder{}
	for s := range c {
		result.WriteString(s)
	}
	got := result.String()

	if got == "" {
		t.Fatalf(`client.Complete(ctx, s) = "%s", want %s`, got, "")
	}

}
