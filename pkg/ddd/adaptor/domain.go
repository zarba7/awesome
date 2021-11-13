package adaptor

import (
	"ddd/dto"
	"ddd/po"
)

type DomainOutBox interface {
	Publish(channel string, events []*dto.Event) error
}

type DomainAggregate interface {
	Tag() string
	Root() int64
	Version() int64
	AppendEvent(eTid uint32, eMsg interface{})
}


type DomainRepo interface {
	LoadAll(agg *po.DomainAggregate) error
	Load(agg *po.DomainAggregate) error
	Save(agg *po.DomainAggregate) error
	UpdateSnapshot(agg *po.DomainAggregate) error
}

