package main

import (
	"fmt"
	"os"
	"strings"
	"io"
	"bufio"
	"regexp"

	"git.dominic-ricottone.com/textwrap/common"
)

// An enumeration of header parts
const(
	HeaderSubject         = "HeaderSubject"
	HeaderDate            = "HeaderDate"
	HeaderFrom            = "HeaderFrom"
	HeaderTo              = "HeaderTo"
	HeaderCc              = "HeaderCc"
	HeaderMessageID       = "HeaderMessageID"
	HeaderContentType     = "HeaderContentType"
	HeaderContentEncoding = "HeaderContentEncoding"
)

// A message header container, used within message containers
type MessageHeader struct {
	Subject     string
	Date        string
	From        string
	To          string
	Cc          string
	MessageID   string
	ContentType string
	LastSet     string
}

// Builder for a message header
func NewHeader() *MessageHeader {
	return &MessageHeader{"", "", "", "", "", "", "", ""}
}

// A message part header container, used within message part containers
type MessagePartHeader struct {
	ContentType     string
	ContentEncoding string
	LastSet         string
}

// Builder for a message part header
func NewPartHeader() *MessagePartHeader {
	return &MessagePartHeader{"", "", ""}
}

// A message part container, used within message containers
type MessagePart struct {
	Header      *MessagePartHeader
	Content     []string
}

// Builder for a message part
func NewPart() *MessagePart {
	return &MessagePart{NewPartHeader(), []string{""}}
}

// A message container
type Message struct {
	Header       *MessageHeader
	Parts        []*MessagePart
	PartBoundary *regexp.Regexp
}

// Builder for a message
func NewMessage() *Message {
	return &Message{NewHeader(), []*MessagePart{NewPart()}, nil}
}

// Message setters
func (m *Message) SetHeader(line string) {
	if strings.HasPrefix(line, "Subject:") {
		m.Header.Subject = line[8:]
		m.Header.LastSet = HeaderSubject
	} else if strings.HasPrefix(line, "Date:") {
		m.Header.Date = line[5:]
		m.Header.LastSet = HeaderDate
	} else if strings.HasPrefix(line, "From:") {
		m.Header.From = line[5:]
		m.Header.LastSet = HeaderFrom
	} else if strings.HasPrefix(line, "To:") {
		m.Header.To = line[3:]
		m.Header.LastSet = HeaderTo
	} else if strings.HasPrefix(line, "Cc:") {
		m.Header.Cc = line[3:]
		m.Header.LastSet = HeaderCc
	} else if strings.HasPrefix(line, "Message-ID:") {
		m.Header.MessageID = line[11:]
		m.Header.LastSet = HeaderMessageID
	} else if strings.HasPrefix(line, "Content-Type:") {
		m.Header.ContentType = line[13:]
		m.Header.LastSet = HeaderContentType
	}
}

func (m *Message) AppendLastHeader(s string) {
	switch m.Header.LastSet {
	case HeaderSubject:
		m.Header.Subject += " " + s
	case HeaderDate:
		m.Header.Date += " " + s
	case HeaderFrom:
		m.Header.From += " " + s
	case HeaderTo:
		m.Header.To += " " + s
	case HeaderCc:
		m.Header.Cc += " " + s
	case HeaderMessageID:
		m.Header.MessageID += " " + s
	case HeaderContentType:
		m.Header.ContentType += " " + s
	}
}

func (m *Message) SetPartHeader(line string) {
	if strings.HasPrefix(line, "Content-Type:") {
		m.Parts[len(m.Parts)-1].Header.ContentType = line[13:]
		m.Parts[len(m.Parts)-1].Header.LastSet = HeaderContentType
	} else if strings.HasPrefix(line, "Content-Transfer-Encoding:") {
		m.Parts[len(m.Parts)-1].Header.ContentEncoding = line[26:]
		m.Parts[len(m.Parts)-1].Header.LastSet = HeaderContentEncoding
	}
}

func (m *Message) AppendLastPartHeader(s string) {
	switch m.Parts[len(m.Parts)-1].Header.LastSet {
	case HeaderContentType:
		m.Parts[len(m.Parts)-1].Header.ContentType += " " + s
	case HeaderContentEncoding:
		m.Parts[len(m.Parts)-1].Header.ContentEncoding += " " + s
	}
}

func (m *Message) AppendPart() {
	m.Parts = append(m.Parts, NewPart())
}

func (m *Message) AppendContent(s string) {
	i := len(m.Parts)-1
	m.Parts[i].Content = append(m.Parts[i].Content, s)
}

func (m *Message) FindBoundary(re *regexp.Regexp) {
	match := re.FindStringSubmatch(m.Header.ContentType)
	if match != nil {
		boundary := strings.Replace(match[1], " ", "", -1)
		m.PartBoundary, _ = regexp.Compile(".*" + boundary + ".*")
	}
}

// Message logic
func (m *Message) MatchBoundary(line string) bool {
	if m.PartBoundary != nil {
		return m.PartBoundary.MatchString(line)
	} else {
		return false
	}
}

// A message printer
func (m *Message) Dump() {
	fmt.Printf("Subject: %s\n", m.Header.Subject)
	fmt.Printf("Date: %s\n", m.Header.Date)
	fmt.Printf("From: %s\n", m.Header.From)
	fmt.Printf("To: %s\n", m.Header.To)
	fmt.Printf("Cc: %s\n", m.Header.Cc)
	fmt.Printf("MessageID: %s\n", m.Header.MessageID)
	fmt.Printf("ContentType: %s\n", m.Header.ContentType)
	for i := 0; i < len(m.Parts); i++ {
		fmt.Printf("Part %d:\n", i)
		fmt.Printf("ContentType: %s\n", m.Parts[i].Header.ContentType)
		fmt.Printf("ContentEncoding: %s\n", m.Parts[i].Header.ContentEncoding)

		wrapped, _ := common.WrapArray(m.Parts[i].Content, 80)
		for j := 0; j < len(wrapped); j++ {
			fmt.Printf("%s\n", wrapped[j])
		}
		fmt.Println("EOF")
	}
}

// Parser statuses
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
	message := NewMessage()

	for input.Scan() {
		line := input.Text()
		tline := strings.TrimSpace(line)

		if parsing == ParsingPreHeader {
			if re_header.MatchString(tline) {
				parsing = ParsingHeader
				message.SetHeader(tline)
			}
		} else if parsing == ParsingHeader {
			if tline == "" {
				parsing = ParsingContent
				message.FindBoundary(re_multipart)
			} else if strings.HasPrefix(line, "\t") {
				message.AppendLastHeader(tline)
			} else {
				message.SetHeader(tline)
			}
		} else if parsing == ParsingPartHeader {
			if tline == "" {
				parsing = ParsingContent
			} else if strings.HasPrefix(line, "\t") {
				message.AppendLastPartHeader(tline)
			} else {
				message.SetPartHeader(tline)
			}
		} else if parsing == ParsingContent {
			if re_message_break.MatchString(tline) {
				parsing = ParsingPreHeader
				message.Dump()
				message = NewMessage()
			} else if message.MatchBoundary(tline) {
				parsing = ParsingPartHeader
				message.AppendPart()
			} else {
				message.AppendContent(tline)
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

