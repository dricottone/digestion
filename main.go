package main

import (
	"fmt"
	"os"
	"strings"
	"io"
	"bufio"
	"regexp"

	"git.dominic-ricottone.com/digestion/message"
)

const (
	ParsingPreHeader  = "ParsingPreHeader"
	ParsingHeader     = "ParsingHeader"
	ParsingPartHeader = "ParsingPartHeader"
	ParsingContent    = "ParsingContent"
)

func parse_stream(reader io.Reader) {
	// Create scanner from reader
	input := bufio.NewScanner(reader)

	// Compile regular expressions
	re_message_break, err := regexp.Compile("^-+$")
	if err != nil {
		fmt.Printf("internal error - %v\n", err)
		os.Exit(1)
	}
	re_header, err := regexp.Compile(
		"^(?:Date|From|Subject|To|Cc|Message-ID|" +
		"Content-(?:Type|Transfer-Encoding)):",
	)
	if err != nil {
		fmt.Printf("internal error - %v\n", err)
		os.Exit(1)
	}
	re_multipart, err := regexp.Compile(".*boundary=\"(.*)\".*")
	if err != nil {
		fmt.Printf("internal error - %v\n", err)
		os.Exit(1)
	}

	parsing := ParsingPreHeader
	current_message := message.NewMessage()

	for input.Scan() {
		line := input.Text()
		tline := strings.TrimSpace(line)

		if parsing == ParsingPreHeader {
			if re_header.MatchString(tline) {
				parsing = ParsingHeader
				current_message.SetHeader(tline)
			}
		} else if parsing == ParsingHeader {
			if tline == "" {
				parsing = ParsingContent
				current_message.FindBoundary(re_multipart)
			} else {
				current_message.SetHeader(line)
			}
		} else if parsing == ParsingPartHeader {
			if tline == "" {
				parsing = ParsingContent
			} else {
				current_message.SetPartHeader(line)
			}
		} else if parsing == ParsingContent {
			if re_message_break.MatchString(tline) {
				parsing = ParsingPreHeader
				current_message.Dump()
				current_message = message.NewMessage()
			} else if current_message.MatchBoundary(tline) {
				parsing = ParsingPartHeader
				current_message.AppendPart()
			} else {
				current_message.AppendContent(tline)
			}
		}
	}

	// Check for scanner errors
	if err = input.Err(); err != nil {
		fmt.Printf("internal error - %v\n", err)
		os.Exit(1)
	}
}

func parse_file(filename string) {
	// Check file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("cannot read file '%s'\n", filename)
		os.Exit(1)
	}
	defer file.Close()

	// Parse
	parse_stream(file)
}

func main() {
	// Check STDIN
	_, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println("cannot read input")
		os.Exit(1)
	}

	// Parse
	parse_stream(os.Stdin)
}

