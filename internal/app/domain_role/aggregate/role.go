package aggregate

import (
	"ddd/adaptor"
	"domain-role/aggregate/modules/bag"
	"domain-role/aggregate/modules/info"
	"google.golang.org/protobuf/proto"
)


type Role struct {
	agg adaptor.DomainAggregate
	Info          info.Entity
	Bag           bag.Entity
}



func (the *Role) PublishDomainEvent(eTid uint32, eMsg proto.Message) {	the.agg.AppendEvent(eTid, eMsg) }
func (the *Role) HandleCommand(handler interface{}) (interface{}, error) {	return handler.(CmdHandler).Handle(the)}
func (the *Role) ApplyEvent(applier interface{}) {	 applier.(EvtApplier).Apply(the)  }

func (the *Role) Load(snapshot []byte) error {

	return nil
}
func (the *Role) Save() (snapshot []byte, err error) {
	panic("implement me")
}


func (the *Role) Uid() uint64 {	return uint64(the.agg.Root()) }

