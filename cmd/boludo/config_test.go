// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"reflect"
	"strings"
	"testing"

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
