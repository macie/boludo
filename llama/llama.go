// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package llama

// DefaultOptions represent neutral parameters for interacting with LLaMA model.
var DefaultOptions = Options{
	ModelPath: "",
	Seed:      0,
	Temp:      1,
	MinP:      0,
}

// Options represent parameters for interacting with LLaMA model.
type Options struct {
	ModelPath string
	Format    string
	Temp      float32
	MinP      float32
	Seed      uint
}
