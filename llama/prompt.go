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
	"alpaca": func(p Prompt) string {
		systemPrompt := ""
		if p.System != "" {
			systemPrompt = fmt.Sprintf("%s\n\n", p.System)
		}

		userPrompt := ""
		for i := range p.userPrompt {
			userPrompt += fmt.Sprintf("### Instruction:\n%s\n\n", p.userPrompt[i])
		}
		if userPrompt == "" {
			userPrompt = "### Instruction:\n\n"
		}

		return fmt.Sprintf("%s%s### Response:\n", systemPrompt, userPrompt)
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

		return fmt.Sprintf("%s%s<|im_start|>assistant\n", systemPrompt, userPrompt)
	},
	"openchat": func(p Prompt) string {
		systemPrompt := p.System
		if systemPrompt != "" {
			systemPrompt += "<|end_of_turn|>"
		}
		userPrompt := ""
		for i := range p.userPrompt {
			userPrompt += fmt.Sprintf("GPT4 Correct User: %s<|end_of_turn|>", p.userPrompt[i])
		}
		if userPrompt == "" {
			userPrompt = "GPT4 Correct User: <|end_of_turn|>"
		}

		return fmt.Sprintf("%s%sGPT4 Correct Assistant: ", systemPrompt, userPrompt)
	},
	"zephyr": func(p Prompt) string {
		systemPrompt := fmt.Sprintf("<|system|>\n%s</s>\n", p.System)
		userPrompt := ""
		for i := range p.userPrompt {
			userPrompt += fmt.Sprintf("<|user|>\n%s</s>\n", p.userPrompt[i])
		}
		if userPrompt == "" {
			userPrompt = "<|user|>\n</s>\n"
		}

		return fmt.Sprintf("%s%s<|assistant|>\n", systemPrompt, userPrompt)
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
