package domain

import (
	"ddd/po"
)

type aggregate struct {
	po  po.DomainAggregate
	srv *service
	//cmd dto.Response
	ss      Session
	events  []*event
	expired int64
}

type event struct {
	tid  uint32
	args interface{}
	//content []byte
}

func (agg *aggregate) AppendEvent(eTid uint32, eMsg interface{}) {
	agg.events = append(agg.events, &event{
		tid:  eTid,
		args: eMsg,
	})
}

func (agg *aggregate) reset() {
	agg.clear()
	agg.po.CurrVersion = 0
	agg.po.SnapshotVersion = 0
}

func (agg *aggregate) clear() {
	agg.events = agg.events[:]
	agg.po.Events = agg.po.Events[:]
}

func (agg *aggregate) Tag() string    { return agg.po.Tag }
func (agg *aggregate) Root() int64    { return agg.po.ID }
func (agg *aggregate) Version() int64 { return agg.po.CurrVersion }
