package message

import (
	"fmt"
	"strings"

	textwrap "git.dominic-ricottone.com/textwrap/common"

	"git.dominic-ricottone.com/digestion/decoder"
)

func (m *Message) determine_best_part() int {
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
		value := m.Parts[i].evaluate_type()
		value += m.Parts[i].evaluate_encoding()
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

func (m *Message) format_header() []string {
	buffer := []string{}
	if m.Header.Subject != "" {
		buffer = append(buffer, fmt.Sprintf("Subject: %s", m.Header.Subject))
	}
	if m.Header.Date != "" {
		buffer = append(buffer, fmt.Sprintf("Date: %s", m.Header.Date))
	}
	if m.Header.From != "" {
		buffer = append(buffer, fmt.Sprintf("From: %s", m.Header.From))
	}
	if m.Header.To != "" {
		buffer = append(buffer, fmt.Sprintf("To: %s", m.Header.To))
	}
	if m.Header.Cc != "" {
		buffer = append(buffer, fmt.Sprintf("Cc: %s", m.Header.Cc))
	}
	//if m.Header.MessageID != "" {
	//	buffer = append(buffer, fmt.Sprintf("MessageID: %s", m.Header.MessageID))
	//}
	//if m.Header.ContentType != "" {
	//	buffer = append(buffer, fmt.Sprintf("ContentType: %s", m.Header.ContentType))
	//}
	return buffer
}

func (m *Message) format_content(length int) ([]string, error) {
	best_part := m.determine_best_part()
	buffer := []string{}

	// Handle messages with no content
	if best_part == -1 {
		return buffer, nil
	}

	// Decode best part's content
	decoded, err := decoder.DecodeArray(m.Parts[best_part].Content, m.Parts[best_part].Header.ContentEncoding)
	if err != nil {
		return buffer, err
	}

	// Wrap text content
	wrapped, err := textwrap.WrapArray(decoded, length)
	if err != nil {
		return decoded, err
	}

	return wrapped, nil
}

func (m *Message) Dump() {
	header := m.format_header()
	content, err := m.format_content(80)
	if err != nil {
		fmt.Printf("error: %s", err)
	}

	for _, line := range append(header, content...) {
		fmt.Printf("%s\n", line)
	}
	fmt.Printf("\n%s\n\n", strings.Repeat("-", 80))
}

