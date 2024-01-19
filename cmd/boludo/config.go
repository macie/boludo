package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/macie/boludo/llama"
)

const helpMsg = "boludo - AI personal assistant\n" +
	"\n" +
	"Usage:\n" +
	"   boludo <CONFIG_ID> [--server PATH] [-t <timeout>] [PROMPT]\n" +
	"   boludo [-h] [-v]\n" +
	"\n" +
	"Options:\n" +
	"   -t <timeout>    timeout after which the program exits (default: 0).\n" +
	"                   Valid time units: ns, us, ms, s, m, h\n" +
	"   --server PATH   path to LLM server executable\n" +
	"   --verbose       show more verbose debug output\n" +
	"   -h              show this help message and exit\n" +
	"   -v              show version information and exit\n" +
	"\n" +
	"boludo reads prompt from PROMPT, and then from standard input"

var AppVersion = "local-dev"

// Version returns string with full version description.
func Version() string {
	return fmt.Sprintf("boludo %s", AppVersion)
}

// AppConfig contains configuration options for the program.
type AppConfig struct {
	Options     llama.Options
	ServerPath  string
	Prompt      llama.Prompt
	UserPrompt  string
	Timeout     time.Duration
	Verbose     bool
	ExitMessage string
}

// NewAppConfig creates a new AppConfig from:
//   - command line arguments
//   - config file
//   - default values
//
// in that order.
func NewAppConfig(cliArgs []string) (AppConfig, error) {
	configArgs, err := ParseArgs(cliArgs)
	if err != nil {
		return AppConfig{}, fmt.Errorf("could not read CLI arguments: %w", err)
	}
	if configArgs.ShowHelp {
		return AppConfig{ExitMessage: helpMsg}, nil
	}
	if configArgs.ShowVersion {
		return AppConfig{ExitMessage: Version()}, nil
	}

	configRoot, err := os.UserConfigDir()
	if err != nil {
		return AppConfig{}, fmt.Errorf("could not locate config directory: %w", err)
	}
	configDir := filepath.Join(configRoot, "boludo")
	configFile, err := ParseFile(os.DirFS(configDir), "boludo.toml")
	switch {
	case err == nil:
		// no error
	case errors.Is(err, os.ErrNotExist):
		// ignore missing config file
	default:
		return AppConfig{}, fmt.Errorf("could not read `%s`: %w", filepath.Join(configDir, "boludo.toml"), err)
	}

	options := llama.DefaultOptions
	options.Update(configFile.Options(configArgs.ConfigId))
	options.Update(configArgs.Options())

	prompt := configFile.Prompt(configArgs.ConfigId)
	initialPrompt := configFile.InitialPrompt(configArgs.ConfigId)
	userPrompt := configArgs.Prompt
	if initialPrompt != "" {
		userPrompt = fmt.Sprintf("%s %s", initialPrompt, userPrompt)
	}

	return AppConfig{
		Prompt:     prompt,
		UserPrompt: userPrompt,
		Options:    options,
		ServerPath: configArgs.ServerPath,
		Timeout:    configArgs.Timeout,
		Verbose:    configArgs.ShowVerbose,
	}, nil
}

// NewAppContext returns a cancellable context which is:
// - cancelled when the interrupt signal is received
// - cancelled after the timeout (if any).
func NewAppContext(config AppConfig) (context.Context, context.CancelFunc) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	if config.Timeout != 0 {
		return context.WithTimeout(ctx, config.Timeout)
	}

	return ctx, cancel
}

// ConfigArgs contains configuration options for the program provided by the user.
type ConfigArgs struct {
	ConfigId    string
	Prompt      string
	Timeout     time.Duration
	ModelPath   string
	ServerPath  string
	ShowHelp    bool
	ShowVersion bool
	ShowVerbose bool
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
	f.DurationVar(&conf.Timeout, "t", 0, "")
	f.BoolVar(&conf.ShowVerbose, "verbose", false, "")
	f.StringVar(&conf.ServerPath, "server", "", "")
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
	Model         string
	SystemPrompt  string
	InitialPrompt string
	Format        string
	Creativity    float32
	Cutoff        float32
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
			Model:         "",
			SystemPrompt:  "",
			InitialPrompt: "",
			Format:        "",
			Creativity:    llama.DefaultOptions.Temp,
			Cutoff:        llama.DefaultOptions.MinP,
		}
		for k, v := range definedConfigs[configId].(map[string]interface{}) {
			switch k {
			case "model":
				defaultSpec.Model = os.ExpandEnv(v.(string))
			case "creativity":
				defaultSpec.Creativity = (float32)(v.(float64))
			case "cutoff":
				defaultSpec.Cutoff = (float32)(v.(float64))
			case "format":
				defaultSpec.Format = v.(string)
			case "system-prompt":
				defaultSpec.SystemPrompt = v.(string)
			case "initial-prompt":
				defaultSpec.InitialPrompt = v.(string)
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
			Temp:      spec.Creativity,
			MinP:      spec.Cutoff,
		}
	}

	return llama.DefaultOptions
}

// Prompt returns the llama.Prompt based on the ConfigFile.
func (c *ConfigFile) Prompt(configId string) llama.Prompt {
	if spec, ok := (*c)[configId]; ok {
		return llama.Prompt{
			Format: spec.Format,
			System: spec.SystemPrompt,
		}
	}
	return llama.Prompt{}
}

// InitialPrompt returns the initial prompt specified in the ConfigFile.
func (c *ConfigFile) InitialPrompt(configId string) string {
	if spec, ok := (*c)[configId]; ok {
		return spec.InitialPrompt
	}
	return ""
}
