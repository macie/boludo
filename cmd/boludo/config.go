// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/macie/boludo/llama"
)

const helpMsg = "boludo - AI personal assistant\n" +
	"\n" +
	"Usage:\n" +
	"   boludo <CONFIG_ID> [PROMPT]\n" +
	"   boludo [-h] [-v]\n" +
	"\n" +
	"Options:\n" +
	"   -h            show this help message and exit\n" +
	"   -v            show version information and exit\n" +
	"\n" +
	"boludo reads prompt from PROMPT, and then from standard input"

// ConfigArgs contains configuration options for the program provided by the user.
type ConfigArgs struct {
	ConfigId    string
	Prompt      string
	ModelPath   string
	ShowHelp    bool
	ShowVersion bool
}

// ParseArgs creates a new ConfigArgs from the given command line arguments.
func ParseArgs(cliArgs []string) (ConfigArgs, error) {
	conf := ConfigArgs{}

	if len(cliArgs) == 0 {
		return ConfigArgs{}, fmt.Errorf(helpMsg)
	}

	// first argument is a config name
	if !strings.HasPrefix(cliArgs[0], "-") {
		conf.ConfigId = cliArgs[0]
		cliArgs = cliArgs[1:]
	}

	// global options
	f := flag.NewFlagSet("boludo", flag.ContinueOnError)
	f.SetOutput(io.Discard)
	f.BoolVar(&conf.ShowHelp, "h", false, "")
	f.BoolVar(&conf.ShowVersion, "v", false, "")
	if err := f.Parse(cliArgs); err != nil {
		return ConfigArgs{}, fmt.Errorf("%w. See 'boludo -h' for help", err)
	}

	switch f.NArg() {
	case 0:
		break
	case 1:
		conf.Prompt = f.Arg(0)
	default:
		return ConfigArgs{}, fmt.Errorf("too much arguments: '%s'. See 'boludo -h' for help", strings.Join(cliArgs, "', '"))
	}

	return conf, nil
}

// Options returns the llama.Options based on the ConfigArgs.
//
// It uses default values from llama.DefaultOptions for options not specified by
// the user.
func (a ConfigArgs) Options() llama.Options {
	options := llama.DefaultOptions
	if a.ModelPath != "" {
		options.ModelPath = a.ModelPath
	}

	return options
}
