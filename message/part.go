package message

import (
	"git.dominic-ricottone.com/digestion/decoder"
)

type MessagePart struct {
	Header      *MessagePartHeader
	Content     []string
}

func NewPart() *MessagePart {
	return &MessagePart{NewPartHeader(), []string{""}}
}

func (m *MessagePart) evaluateContentType() int {
	return 0
}

func (m *MessagePart) evaluateContentEncoding() int {
	return decoder.EvaluateEncoding(m.Header.ContentEncoding)
}


