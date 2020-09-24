package server

// Message structure expected from publishers
type Message struct {
	Payload string
	Topic   string
}
