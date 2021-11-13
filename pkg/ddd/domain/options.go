package domain

import (
	"ddd/adaptor"
	"ddd/dto"
	"ddd/util"
)

type Codec interface {
	Marshal(val interface{}) ([]byte, error)
	Unmarshal(data []byte, val interface{}) error
}

type Input interface {
	Start(do func(msg *dto.Message) error) error
	Stop()error
}

type Options struct {
	Namespace      string
	Repo           adaptor.DomainRepo
	SessionFactory func(agg adaptor.DomainAggregate) Session
	CommandFactory util.SameNameStructFactory
	EventFactory   util.SameNameStructFactory
	OutBox         adaptor.DomainOutBox
	Codec Codec
	Input Input
}
