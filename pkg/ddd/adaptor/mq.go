package adaptor

type MessageProcessor interface {
	Command() error
}

type Consumer interface {
	Watch(newProcessor func ([]byte) MessageProcessor)error
	Stop() error
}
type Producer interface {
	Sync(tid int32, content []byte) error
	Stop() error
}