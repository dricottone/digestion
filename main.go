package main

import (
	"fmt"
	"os"
	"strings"
	"io"
	"bufio"
	"regexp"
	"flag"

	"git.dominic-ricottone.com/~dricottone/digestion/message"
	parcels "git.dominic-ricottone.com/~dricottone/parcels/common"
)

const LINE_LENGTH = 80

const (
	ParsingPreHeader  = "ParsingPreHeader"
	ParsingHeader     = "ParsingHeader"
	ParsingPartHeader = "ParsingPartHeader"
	ParsingContent    = "ParsingContent"
)

func contains_nonempty(s1, s2 string) bool {
	return s2 != "" && strings.Contains(s1, s2)
}

func first_submatch(r regexp.Regexp, s string) string {
	matches := r.FindStringSubmatch(s)
	if matches != nil {
		return strings.Replace(matches[1], " ", "", -1)
	} else {
		return ""
	}
}

func parse_stream(reader io.Reader, length int) {
	fmt.Println("ping")
	content, urls, err := parcels.ParseFromReader(reader, 0)
	if err != nil {
		fmt.Printf("internal error - %v\n", err)
		os.Exit(1)
	}

	// create scanner
	input := bufio.NewScanner(strings.NewReader(content))

	// Compile regular expressions
	re_message_break, err := regexp.Compile("^-{5,}$")
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

	parsing := ParsingContent
	current_message := message.NewMessage()
	current_boundary := ""

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
				current_boundary = first_submatch(*re_multipart, current_message.Header.ContentType)
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
				current_message.Dump(length)
				current_message = message.NewMessage()
				current_boundary = ""
			} else if contains_nonempty(tline, current_boundary) {
				parsing = ParsingPartHeader
				current_message.AppendPart()
			} else {
				current_message.AppendContent(tline)
			}
		}
	}

	fmt.Printf("%s", urls)

	// Check for scanner errors
	if err = input.Err(); err != nil {
		fmt.Printf("internal error - %v\n", err)
		os.Exit(1)
	}
}

func parse_file(filename string, length int) {
	// Check file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("cannot read file '%s'\n", filename)
		os.Exit(1)
	}
	defer file.Close()

	// Parse
	parse_stream(file, length)
}

func main() {
	// Check STDIN
	_, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println("cannot read input")
		os.Exit(1)
	}

	// Look for arguments
	var length = flag.Int("length", LINE_LENGTH, "maximum length of lines")
	flag.Parse()

	// Parse
	parse_stream(os.Stdin, *length)
}

