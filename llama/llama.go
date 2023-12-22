// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package llama

// Options represent parameters for interacting with LLAMA model.
type Options struct {
	ModelPath string
	Seed      uint
	Temp      float32
	Min_P     float32
}
