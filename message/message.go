package message

import (
	"fmt"
	"strings"
	"regexp"

	textwrap "git.dominic-ricottone.com/textwrap/common"

	"git.dominic-ricottone.com/digestion/decoder"
)

type Message struct {
	Header       *MessageHeader
	Parts        []*MessagePart
	PartBoundary *regexp.Regexp
}

func NewMessage() *Message {
	return &Message{NewHeader(), []*MessagePart{NewPart()}, nil}
}

func (m *Message) SetHeader(s string) {
	m.Header.SetHeader(s)
}

func (m *Message) SetPartHeader(s string) {
	m.Parts[len(m.Parts)-1].Header.SetHeader(s)
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

func (m *Message) MatchBoundary(line string) bool {
	if m.PartBoundary != nil {
		return m.PartBoundary.MatchString(line)
	} else {
		return false
	}
}

func (m *Message) DetermineBestPart() int {
	// Handle cases with obvious best part
	number_parts := len(m.Parts)
	if number_parts == 0 {
		return -1
	} else if number_parts == 1 {
		return 0
	}

	// Evaluate each part--lower is better
	evaluations := []int{}
	for i := 0; i < number_parts; i++ {
		value := m.Parts[i].evaluateContentType()
		value += m.Parts[i].evaluateContentEncoding()
		evaluations = append(evaluations, value)
	}

	// Find minimum value and return that part index
	best_part_index := 0
	for i := 1; i < number_parts; i++ {
		if evaluations[i] < evaluations[best_part_index] {
			best_part_index = i
		}
	}
	return best_part_index
}

func (m *Message) Dump() {
	fmt.Printf("Subject: %s\n", m.Header.Subject)
	fmt.Printf("Date: %s\n", m.Header.Date)
	fmt.Printf("From: %s\n", m.Header.From)
	fmt.Printf("To: %s\n", m.Header.To)
	fmt.Printf("Cc: %s\n", m.Header.Cc)
	fmt.Printf("MessageID: %s\n", m.Header.MessageID)
	fmt.Printf("ContentType: %s\n", m.Header.ContentType)

	if index := m.DetermineBestPart(); index != -1 {
		fmt.Printf("ContentType: %s\n", m.Parts[index].Header.ContentType)
		fmt.Printf("ContentEncoding: %s\n", m.Parts[index].Header.ContentEncoding)

		decoded, _ := decoder.DecodeArray(m.Parts[index].Content, m.Parts[index].Header.ContentEncoding)

		wrapped, _ := textwrap.WrapArray(decoded, 80)

		for _, line := range wrapped {
			fmt.Printf("%s\n", line)
		}
	}
}

