package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/macie/boludo"
	"github.com/macie/boludo/llama"
)

func main() {
	defaultLogHandler := boludo.UnstructuredHandler{Prefix: "[boludo]", Level: slog.LevelError}
	slog.SetDefault(slog.New(defaultLogHandler))

	config, err := NewAppConfig(os.Args[1:])
	if err != nil {
		slog.Error(fmt.Sprint(err))
		os.Exit(1)
	}

	if config.ExitMessage != "" {
		fmt.Fprintln(os.Stdin, config.ExitMessage)
		os.Exit(0)
	}
	if config.Verbose {
		defaultLogHandler.Level = slog.LevelInfo
		slog.SetDefault(slog.New(defaultLogHandler))
	}

	ctx, cancel := NewAppContext(config)
	defer cancel()

	server := llama.Server{
		Path:   config.ServerPath,
		Logger: slog.New(boludo.UnstructuredHandler{Prefix: "[llm-server]", Level: defaultLogHandler.Level}),
	}
	client := llama.Client{
		Options: &config.Options,
		Logger:  slog.New(boludo.UnstructuredHandler{Prefix: "[llm-client]", Level: defaultLogHandler.Level}),
	}
	llama.SetDefault(server, client)

	if err := llama.Serve(ctx, config.Options.ModelPath); err != nil {
		slog.Error(fmt.Sprint(err))
		os.Exit(1)
	}
	defer llama.Close()

	userPrompt := strings.Builder{}
	userPrompt.WriteString(config.UserPrompt)

	input, err := os.Stdin.Stat()
	if err == nil && (input.Mode()&os.ModeCharDevice) == 0 {
		// something is redirected to stdin
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			userPrompt.Write(scanner.Bytes())
			userPrompt.WriteRune('\n')
		}
	}

	config.Prompt.Add(strings.Trim(userPrompt.String(), "\n"))

	output, err := llama.Complete(ctx, config.Prompt)
	if err != nil {
		slog.Error(fmt.Sprint(err))
		os.Exit(1)
	}

	for token := range output {
		fmt.Fprint(os.Stdout, token)
	}

	switch ctx.Err() {
	case nil:
		// no error
	case context.Canceled:
		slog.Info("completion cancelled by user")
	case context.DeadlineExceeded:
		slog.Info("completion needs more time than expected")
	default:
		slog.Info("completion was interrupted")
	}

	os.Exit(0)
}
