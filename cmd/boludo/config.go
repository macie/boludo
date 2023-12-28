// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/BurntSushi/toml"
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

// ConfigFile represents a configuration file.
type ConfigFile map[string]ModelSpec

// ModelSpec represents a model specification in the configuration file.
type ModelSpec struct {
	Model      string
	Format     string
	Creativity float32
	Cutoff     float32
}

// ParseFile reads the TOML configuration file and returns a ConfigFile.
//
// Config file should exist in default location (see `os.UserConfigDir()`):
//   - on Linux in `$XDG_CONFIG_HOME/boludo/boludo.toml` or `$HOME/.config/boludo/boludo.toml`
//   - on macOS in `$HOME/Library/Application Support/boludo/boludo.toml`
//   - on Windows in `%APPDATA%\boludo\boludo.toml`.
//
// If the file doesn't exist, an empty ConfigFile is returned.
func ParseFile(configDir fs.FS, filename string) (ConfigFile, error) {
	file, err := configDir.Open(filename)
	if err != nil {
		return ConfigFile{}, fmt.Errorf("could not open '%s': %w", filename, err)
	}
	defer file.Close()

	config := ConfigFile{}
	_, decodeErr := toml.NewDecoder(file).Decode(&config)
	if decodeErr != nil {
		return ConfigFile{}, fmt.Errorf("could not read config file: %w", decodeErr)
	}

	return config, nil
}

// UnmarshalTOML implements toml.Unmarshaler interface.
//
// Values not defined in the config file will be set to the default values
func (c *ConfigFile) UnmarshalTOML(data interface{}) error {
	definedConfigs, _ := data.(map[string]interface{})
	for configId := range definedConfigs {
		defaultSpec := ModelSpec{
			Model:      "",
			Format:     llama.DefaultOptions.Format,
			Creativity: llama.DefaultOptions.Temp,
			Cutoff:     llama.DefaultOptions.MinP,
		}
		for k, v := range definedConfigs[configId].(map[string]interface{}) {
			switch k {
			case "model":
				defaultSpec.Model = v.(string)
			case "creativity":
				defaultSpec.Creativity = (float32)(v.(float64))
			case "cutoff":
				defaultSpec.Cutoff = (float32)(v.(float64))
			}
		}
		(*c)[configId] = defaultSpec
	}
	return nil
}

// Options returns the llama.Options based on the ConfigFile.
//
// It uses default values from llama.DefaultOptions for options not specified in
// config file.
func (c *ConfigFile) Options(configId string) llama.Options {
	if spec, ok := (*c)[configId]; ok {
		return llama.Options{
			ModelPath: spec.Model,
			Format:    spec.Format,
			Temp:      spec.Creativity,
			MinP:      spec.Cutoff,
		}
	}

	return llama.DefaultOptions
}
