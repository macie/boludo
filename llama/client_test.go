// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package llama

import (
	"context"
	"strings"
	"testing"
)

func TestComplete(t *testing.T) {
	server := NewServer(context.TODO(), serverPath, Options{ModelPath: testingModel})
	defer server.Close()

	if err := server.Start(); err != nil {
		t.Fatalf("server.Start() returns error: %v", err)
	}

	client := NewClient("http://localhost:24114")

	c, err := client.Complete(context.TODO(), "Once upon a time")
	if err != nil {
		t.Fatalf("client.Complete() returns error: %v", err)
	}

	result := strings.Builder{}
	for s := range c {
		result.WriteString(s)
	}
	got := result.String()

	if got == "" {
		t.Fatalf(`client.Complete(ctx, s) = "%s", want %s`, got, "")
	}

}
