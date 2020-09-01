package domain

type Message struct {
	Text string
}

type Client interface {
	SendMessage(m Message) error
	Close()
}
