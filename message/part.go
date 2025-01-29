package message

import (
	"git.sr.ht/~dricottone/digestion/decoder"
	"git.sr.ht/~dricottone/digestion/renderer"
)

type MessagePart struct {
	Header             *MessagePartHeader
	Content            []string
}

func NewPart() *MessagePart {
	return &MessagePart{NewPartHeader(), []string{""}}
}

func (m *MessagePart) evaluate_type() int {
	switch renderer.DetermineFormatting(m.Header.ContentType) {
	case renderer.FormattedTextPlain:
		return 0
	case renderer.FormattedTextHTML:
		return 20
	default:
		return 30
	}
}

func (m *MessagePart) evaluate_encoding() int {
	switch decoder.DetermineEncoding(m.Header.ContentEncoding) {
	case decoder.EncodedUTF8:
		return 0
	case decoder.EncodedBase64:
		return 1
	case decoder.EncodedQuotedPrintable:
		return 1
	default:
		return 10
	}
}

