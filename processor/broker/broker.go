package broker

type Broker interface {
	FetchMessages() error
	Close() error
}
