package models

// Message holds the text written and the writer of it.
type Message struct {
	User string
	Text string
}

// NewMessage creates a new message.
func NewMessage(user string, text string) Message {
	return Message{User: user, Text: text}
}
