package llama

import (
	"strings"
	"testing"
)

func TestPromptString(t *testing.T) {
	testcases := []struct {
		prompt Prompt
		want   string
	}{
		{Prompt{Format: "alpaca"}, "### Instruction:\n\n### Response:\n"},
		{Prompt{Format: "Alpaca", System: "You are a helpful assistant."}, "You are a helpful assistant.\n\n### Instruction:\n\n### Response:\n"},
		{Prompt{Format: "chatml"}, "<|im_start|>system\n<|im_end|>\n<|im_start|>user\n<|im_end|>\n<|im_start|>assistant\n"},
		{Prompt{Format: "ChatML", System: "You are a helpful assistant."}, "<|im_start|>system\nYou are a helpful assistant.<|im_end|>\n<|im_start|>user\n<|im_end|>\n<|im_start|>assistant\n"},
		{Prompt{Format: "openchat"}, "GPT4 Correct User: <|end_of_turn|>GPT4 Correct Assistant: "},
		{Prompt{Format: "openchat", System: "You are a helpful assistant."}, "You are a helpful assistant.<|end_of_turn|>GPT4 Correct User: <|end_of_turn|>GPT4 Correct Assistant: "},
		{Prompt{Format: "Zephyr"}, "<|system|>\n</s>\n<|user|>\n</s>\n<|assistant|>\n"},
		{Prompt{Format: "Zephyr", System: "You are a helpful assistant."}, "<|system|>\nYou are a helpful assistant.</s>\n<|user|>\n</s>\n<|assistant|>\n"},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.prompt.String(), func(t *testing.T) {
			t.Parallel()
			got := tc.prompt.String()
			if got != tc.want {
				t.Fatalf("Prompt.String() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestPromptAdd(t *testing.T) {
	testcases := []struct {
		format      string
		userPrompts []string
		want        string
	}{
		{"Alpaca", []string{}, "### Instruction:\n\n### Response:\n"},
		{"alpaca", []string{"How are you?"}, "### Instruction:\nHow are you?\n\n### Response:\n"},
		{"ChatML", []string{}, "<|im_start|>system\n<|im_end|>\n<|im_start|>user\n<|im_end|>\n<|im_start|>assistant\n"},
		{"chatml", []string{"How are you?"}, "<|im_start|>system\n<|im_end|>\n<|im_start|>user\nHow are you?<|im_end|>\n<|im_start|>assistant\n"},
		{"OpenChat", []string{}, "GPT4 Correct User: <|end_of_turn|>GPT4 Correct Assistant: "},
		{"openchat", []string{"How are you?"}, "GPT4 Correct User: How are you?<|end_of_turn|>GPT4 Correct Assistant: "},
		{"Zephyr", []string{}, "<|system|>\n</s>\n<|user|>\n</s>\n<|assistant|>\n"},
		{"zephyr", []string{"How are you?"}, "<|system|>\n</s>\n<|user|>\nHow are you?</s>\n<|assistant|>\n"},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(strings.Join(tc.userPrompts, "; "), func(t *testing.T) {
			t.Parallel()
			prompt := Prompt{Format: tc.format}
			for i := range tc.userPrompts {
				prompt.Add(tc.userPrompts[i])
			}
			got := prompt.String()
			if got != tc.want {
				t.Fatalf("Prompt{Format: %v}.String() = %v, want %v", tc.format, got, tc.want)
			}
		})
	}
}
