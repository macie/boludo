// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package llama

/*
#cgo CFLAGS: -I${SRCDIR}/../external/llama.cpp
#cgo LDFLAGS: -L${SRCDIR}/../external/llama.cpp -Wl,-R${SRCDIR}/../external/llama.cpp -lllama
#include "llama.h"
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"unsafe"
)

// Server represents computation server.
type Server struct{}

// NewServer creates new server.
//
// It is the caller's responsibility to close Server.
func NewServer() Server {
	C.llama_backend_init(false)
	return Server{}
}

// SystemInfo returns system information.
func (Server) SystemInfo() string {
	return C.GoString(C.llama_print_system_info())
}

// Close frees all resources associated with server.
func (Server) Close() {
	C.llama_backend_free()
}

// Options represent parameters for interacting with LLAMA model.
type Options struct {
	ModelPath string
	Seed      uint
	Temp      float32
	Min_P     float32
}

// Model represents LLAMA model.
type Model struct {
	context *C.struct_llama_context
	path    string
}

// NewModel creates new model from given parameters.
//
// It is the caller's responsibility to close Model.
func NewModel(options Options) (Model, error) {
	modelPath := C.CString(options.ModelPath)
	defer C.free(unsafe.Pointer(modelPath))
	modelParams := C.llama_model_default_params()
	model := C.llama_load_model_from_file(modelPath, modelParams)
	if model == nil {
		return Model{}, errors.New("couldn't load model")
	}

	ctxParams := C.llama_context_default_params()
	ctxParams.n_ctx = C.uint(0) // context from model
	ctxParams.n_threads = C.uint(runtime.NumCPU())
	ctxParams.n_threads_batch = C.uint(runtime.NumCPU())
	if options.Seed > 0 {
		ctxParams.seed = C.uint(options.Seed)
	}
	context := C.llama_new_context_with_model(model, ctxParams)

	return Model{context: context, path: options.ModelPath}, nil
}

// String returns string representation of model.
func (m Model) String() string {
	model := C.llama_get_model(m.context)

	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("Model (path: %s; ", m.path))
	s.WriteString(fmt.Sprintf("context size (trainied): %d (%d); ", C.llama_n_ctx(m.context), C.llama_n_ctx_train(model)))
	s.WriteString(fmt.Sprintf("vocabulary size: %d; ", C.llama_n_vocab(model)))
	s.WriteString(fmt.Sprintf("embeddings size: %d)", C.llama_n_embd(model)))

	return s.String()
}

// MaxContextSize returns maximum context size of model.
func (m Model) MaxContextSize() int {
	return int(C.llama_n_ctx_train(C.llama_get_model(m.context)))
}

// Close frees all resources associated with model.
func (m Model) Close() {
	C.llama_free(m.context)
	C.llama_free_model(C.llama_get_model(m.context))
}
