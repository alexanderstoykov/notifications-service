package service

type Sender interface {
	Send(message *Message) error
}

type Message struct {
	Sender   string
	Message  string
	Receiver string
}
