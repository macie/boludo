// Copyright (C) 2023 Maciej Żok
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	fmt.Fprintln(os.Stdout, "Hello! I'm boludo. You can ask me anything or enter empty line to exit.")
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "How can I help you?")

	for {
		fmt.Fprint(os.Stdout, "> ")
		output, err := readline(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}

		if fmt.Sprint(output) == "" {
			fmt.Fprintln(os.Stdout, "< Goodbye!")
			break
		}

		fmt.Fprintf(os.Stdout, "< I don't understand: %s", output)
		fmt.Fprintln(os.Stdout, "")
		fmt.Fprintln(os.Stdout, "")
	}

	os.Exit(0)
}

func readline(r io.Reader) (io.Writer, error) {
	s := bufio.NewScanner(r)
	output := new(strings.Builder)

	s.Scan()
	output.WriteString(strings.TrimRight(s.Text(), " "))
	if err := s.Err(); err != nil {
		return nil, err
	}

	return output, nil
}
