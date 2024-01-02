// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
		{Prompt{Format: "chatml"}, "<|im_start|>system\n<|im_end|>\n<|im_start|>user\n<|im_end|>\n<|im_start|>assistant\n"},
		{Prompt{Format: "ChatML", System: "You are a helpful assistant."}, "<|im_start|>system\nYou are a helpful assistant.<|im_end|>\n<|im_start|>user\n<|im_end|>\n<|im_start|>assistant\n"},
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
		userPrompts []string
		want        string
	}{
		{[]string{}, "<|im_start|>system\n<|im_end|>\n<|im_start|>user\n<|im_end|>\n<|im_start|>assistant\n"},
		{[]string{"How are you?"}, "<|im_start|>system\n<|im_end|>\n<|im_start|>user\nHow are you?<|im_end|>\n<|im_start|>assistant\n"},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(strings.Join(tc.userPrompts, "; "), func(t *testing.T) {
			t.Parallel()
			prompt := Prompt{
				Format: "ChatML",
			}
			for i := range tc.userPrompts {
				prompt.Add(tc.userPrompts[i])
			}
			got := prompt.String()
			if got != tc.want {
				t.Fatalf("Prompt.String() = %v, want %v", got, tc.want)
			}
		})
	}
}
