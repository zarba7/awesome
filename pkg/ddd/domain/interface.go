package domain

type Session interface {
	HandleCommand(handler interface{}) (result interface{}, err error)
	ApplyEvent(applier interface{})
	Load(snapshot []byte) error
	Save() (snapshot []byte, err error)
}

type EventSubmitter interface {
	Apply()
	Submit()
}
