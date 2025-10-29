package eventbus

// Publisher is an interface for publishing events.
type Publisher interface {
	Publish(topic string, message []byte) error
}

// Subscriber is an interface for subscribing to events.
type Subscriber interface {
	Subscribe(topic string, handler func(message []byte)) error
}
