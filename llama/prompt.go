// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package llama

import (
	"fmt"
	"strings"
)

// supported prompt formats
var promptFormats = map[string]func(Prompt) string{
	"": func(p Prompt) string {
		return fmt.Sprintf("%s\n%s", p.System, strings.Join(p.userPrompt, "\n"))
	},
	"chatml": func(p Prompt) string {
		systemPrompt := fmt.Sprintf("<|im_start|>system\n%s<|im_end|>\n", p.System)
		userPrompt := ""
		for i := range p.userPrompt {
			userPrompt += fmt.Sprintf("<|im_start|>user\n%s<|im_end|>\n", p.userPrompt[i])
		}
		if userPrompt == "" {
			userPrompt = "<|im_start|>user\n<|im_end|>\n"
		}

		return fmt.Sprintf("%s%s<|im_start|>assistant", systemPrompt, userPrompt)
	},
}

// Prompt represents prompt for the LLM.
type Prompt struct {
	Format     string
	System     string
	userPrompt []string
}

// String returns prompt string in format specified by Format.
// If Format is not specified, returns prompt in default format.
func (p *Prompt) String() string {
	formatFunc, ok := promptFormats[strings.ToLower(p.Format)]
	if !ok {
		formatFunc = promptFormats[""]
	}

	return formatFunc(*p)
}

// Add adds user prompt to the prompt.
func (p *Prompt) Add(userPrompt string) {
	p.userPrompt = append(p.userPrompt, userPrompt)
}
