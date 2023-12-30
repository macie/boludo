// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/BurntSushi/toml"
	"github.com/macie/boludo/llama"
)

func TestParseArgs(t *testing.T) {
	testcases := []struct {
		args []string
		want ConfigArgs
	}{
		{[]string{"-h"}, ConfigArgs{ShowHelp: true}},
		{[]string{"-v"}, ConfigArgs{ShowVersion: true}},
		{[]string{"edit"}, ConfigArgs{ConfigId: "edit"}},
		{[]string{"edit", "-h"}, ConfigArgs{ConfigId: "edit", ShowHelp: true}},
		{[]string{"assistant", "--verbose"}, ConfigArgs{ConfigId: "assistant", ShowVerbose: true}},
		{[]string{"assistant", "--server", "./llm-server"}, ConfigArgs{ConfigId: "assistant", ServerPath: "./llm-server"}},
		{[]string{"chat", "How are you?"}, ConfigArgs{ConfigId: "chat", Prompt: "How are you?"}},
		{[]string{"chat", "-v", "How are you?"}, ConfigArgs{ConfigId: "chat", Prompt: "How are you?", ShowVersion: true}},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(strings.Join(tc.args, "_"), func(t *testing.T) {
			t.Parallel()
			got, err := ParseArgs(tc.args)
			if err != nil {
				t.Fatalf("ParseArgs(%v) returns error: %v", tc.args, err)
			}
			if got != tc.want {
				t.Fatalf("ParseArgs(%v) = %v, want %v", tc.args, got, tc.want)
			}
		})
	}
}

func TestParseArgs_Invalid(t *testing.T) {
	testcases := []struct {
		args []string
	}{
		{[]string{}},
		{[]string{"-x"}},
		{[]string{"edit", "--yyy"}},
		{[]string{"assistant", "--server"}},
		{[]string{"chat", "prompt", "prompt2"}},
	}
	want := ConfigArgs{}
	for _, tc := range testcases {
		tc := tc
		t.Run(strings.Join(tc.args, "_"), func(t *testing.T) {
			t.Parallel()
			got, err := ParseArgs(tc.args)
			if err == nil {
				t.Fatalf("ParseArgs(%v) does not return error", tc.args)
			}
			if got != want {
				t.Fatalf("ParseArgs(%v) = %v, want %v", tc.args, got, want)
			}
		})
	}
}

func TestConfigArgsOptions(t *testing.T) {
	testcases := []struct {
		args ConfigArgs
		want llama.Options
	}{
		{ConfigArgs{}, llama.DefaultOptions},
		{ConfigArgs{ModelPath: "model.gguf"}, llama.Options{ModelPath: "model.gguf", Format: "", Temp: 1, MinP: 0}},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.args.ConfigId, func(t *testing.T) {
			t.Parallel()
			got := tc.args.Options()
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf(".Options() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	testcases := []struct {
		content string
		want    ConfigFile
	}{
		{"", ConfigFile{}},
		{"[chat]\nmodel = \"model.gguf\"\ncreativity = 1.2\n cutoff = 0.5\n", ConfigFile{
			"chat": ModelSpec{
				Model:      "model.gguf",
				Format:     "",
				Creativity: 1.2,
				Cutoff:     0.5,
			},
		}},
		{"[edit]\nmodel = \"model.gguf\"\ncutoff = 0.5\n[unknown]", ConfigFile{
			"edit": ModelSpec{
				Model:      "model.gguf",
				Format:     "",
				Creativity: 1.0,
				Cutoff:     0.5,
			},
			"unknown": ModelSpec{
				Model:      "",
				Format:     "",
				Creativity: 1.0,
				Cutoff:     0.0,
			},
		}},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(strings.Split(tc.content, "\n")[0], func(t *testing.T) {
			t.Parallel()
			fs := fstest.MapFS{
				"boludo.toml": {Data: []byte(tc.content)},
			}

			got, err := ParseFile(fs, "boludo.toml")
			if err != nil {
				t.Fatalf("ParseFile(fs, \"boludo.toml\") returns error: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("ParseFile(fs, \"boludo.toml\") = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestParseFile_Invalid(t *testing.T) {
	filename := "invalid.toml"
	confDir := fstest.MapFS{
		filename: {Data: []byte("[invalid")},
	}
	var parseErr toml.ParseError
	want := ConfigFile{}

	got, err := ParseFile(confDir, filename)
	if !errors.As(err, &parseErr) {
		t.Fatalf("ParseFile(fs, \"%s\") want error %T; got: %v", filename, parseErr, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseFile(fs, \"%s\") = %v, want %v", filename, got, want)
	}
}

func TestParseFile_Missing(t *testing.T) {
	filename := "missing.toml"
	confDir := fstest.MapFS{}
	want := ConfigFile{}

	got, err := ParseFile(confDir, filename)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("ParseFile(fs, \"%s\") want error `%v`; got: `%v`", filename, os.ErrNotExist, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseFile(fs, \"%s\") = %v, want %v", filename, got, want)
	}

}

func TestConfigFileOptions(t *testing.T) {
	testcases := []struct {
		configId string
		file     ConfigFile
		want     llama.Options
	}{
		{"chat", ConfigFile{"edit": ModelSpec{Model: "editmodel.gguf"}, "chat": ModelSpec{Model: "chatmodel.gguf", Format: "", Creativity: 0.3, Cutoff: 2}}, llama.Options{ModelPath: "chatmodel.gguf", Format: "", Temp: 0.3, MinP: 2}},
		{"invalid", ConfigFile{}, llama.DefaultOptions},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.configId, func(t *testing.T) {
			t.Parallel()
			got := tc.file.Options(tc.configId)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf(".Options(\"%s\") = %v, want %v", tc.configId, got, tc.want)
			}
		})
	}
}
