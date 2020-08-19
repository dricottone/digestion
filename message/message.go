package message

type Message struct {
	Header       *MessageHeader
	Parts        []*MessagePart
}

func NewMessage() *Message {
	return &Message{NewHeader(), []*MessagePart{NewPart()}}
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

