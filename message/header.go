package message

import (
	"strings"
)

const (
	HeaderSubject         = "HeaderSubject"
	HeaderDate            = "HeaderDate"
	HeaderFrom            = "HeaderFrom"
	HeaderTo              = "HeaderTo"
	HeaderCc              = "HeaderCc"
	HeaderMessageID       = "HeaderMessageID"
	HeaderContentType     = "HeaderContentType"
	HeaderContentEncoding = "HeaderContentEncoding"
)

// Message headers
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

func NewHeader() *MessageHeader {
	return &MessageHeader{"", "", "", "", "", "", "", ""}
}

func (m *MessageHeader) SetHeader(s string) {
	if strings.HasPrefix(s, "\t") {
		m.append_last_set(s)
	} else if strings.HasPrefix(s, "Subject:") {
		m.Subject = strings.TrimSpace(s[8:])
		m.LastSet = HeaderSubject
	} else if strings.HasPrefix(s, "Date:") {
		m.Date = strings.TrimSpace(s[5:])
		m.LastSet = HeaderDate
	} else if strings.HasPrefix(s, "From:") {
		m.From = strings.TrimSpace(s[5:])
		m.LastSet = HeaderFrom
	} else if strings.HasPrefix(s, "To:") {
		m.To = strings.TrimSpace(s[3:])
		m.LastSet = HeaderTo
	} else if strings.HasPrefix(s, "Cc:") {
		m.Cc = strings.TrimSpace(s[3:])
		m.LastSet = HeaderCc
	} else if strings.HasPrefix(s, "Message-ID:") {
		m.MessageID = strings.TrimSpace(s[11:])
		m.LastSet = HeaderMessageID
	} else if strings.HasPrefix(s, "Content-Type:") {
		m.ContentType = strings.TrimSpace(s[13:])
		m.LastSet = HeaderContentType
	}
}

func (m *MessageHeader) append_last_set(s string) {
	s = strings.TrimSpace(s)
	switch m.LastSet {
	case HeaderSubject:
		m.Subject += " " + s
	case HeaderDate:
		m.Date += " " + s
	case HeaderFrom:
		m.From += " " + s
	case HeaderTo:
		m.To += " " + s
	case HeaderCc:
		m.Cc += " " + s
	case HeaderMessageID:
		m.MessageID += " " + s
	case HeaderContentType:
		m.ContentType += " " + s
	}
}

// Message part headers
type MessagePartHeader struct {
	ContentType     string
	ContentEncoding string
	LastSet         string
}

func NewPartHeader() *MessagePartHeader {
	return &MessagePartHeader{"", "", ""}
}

func (m *MessagePartHeader) SetHeader(s string) {
	if strings.HasPrefix(s, "\t") {
		m.append_last_set(s)
	} else if strings.HasPrefix(s, "Content-Type:") {
		m.ContentType = strings.TrimSpace(s[13:])
		m.LastSet = HeaderContentType
	} else if strings.HasPrefix(s, "Content-Transfer-Encoding:") {
		m.ContentEncoding = strings.TrimSpace(s[26:])
		m.LastSet = HeaderContentEncoding
	}
}

func (m *MessagePartHeader) append_last_set(s string) {
	s = strings.TrimSpace(s)
	switch m.LastSet {
	case HeaderContentType:
		m.ContentType += " " + s
	case HeaderContentEncoding:
		m.ContentEncoding += " " + s
	}
}

