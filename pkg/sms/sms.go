package sms

type Message struct {
	Content string
}

type SMS interface {
	Send(phone string, message Message) error
}
