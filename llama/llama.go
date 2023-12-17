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
