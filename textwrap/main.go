package main

import (
	"fmt"
	"os"
	"strings"
	"bufio"
	"regexp"
	"flag"
)

const LENGTH = 40

func print_break(length int) {
	fmt.Printf("%s\n", strings.Repeat("-", length))
}

func print_wrapped(line string, length int, quote string) {
	len_quote := len(quote)
	buffer := quote
	for index, rune := range line[len_quote:] {
		buffer += string(rune)
		if (index + 1) % (length - len_quote) == 0 {
			fmt.Printf("%s\n", buffer)
			buffer = quote
		}
	}
	if buffer != "" {
		fmt.Printf("%s\n", buffer)
	}
}

func main() {
	// Open STDIN as scanner
	_, err := os.Stdin.Stat()
	if err != nil {
		fmt.Printf("%s\n", "cannot read input")
		os.Exit(1)
	}
	input := bufio.NewScanner(os.Stdin)

	// Look for arguments
	var width = flag.Int("width", 80, "target width for output")
	flag.Parse()

	// Compile regular expressions
	re_quote, err := regexp.Compile("^([> ]*)")
	if err != nil {
		fmt.Printf("internal error - %v\n", err)
		os.Exit(1)
	}
	re_break, err := regexp.Compile("^(?:-{5,}|={5,})$")
	if err != nil {
		fmt.Printf("internal error - %v\n", err)
		os.Exit(1)
	}

	// Scan line by line
	for input.Scan() {
		line := input.Text()
		line = strings.TrimSpace(line)

		if len(line) > *width {
			if re_break.MatchString(line) {
				print_break(*width)
			} else {
				quote := re_quote.FindString(line)
				print_wrapped(line, *width, quote)
			}
		} else {
			fmt.Printf("%s\n", line)
		}
	}

	// Check for scanner errors
	if err = input.Err(); err != nil {
		fmt.Printf("internal error - %v\n", err)
		os.Exit(1)
	}
}

