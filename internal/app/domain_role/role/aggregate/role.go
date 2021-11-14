package aggregate

import (
	"awesome/internal/app/domain_role/role/aggregate/modules/bag"
	"awesome/internal/app/domain_role/role/aggregate/modules/info"
	"ddd/adaptor"
	"ddd/domain"
	"google.golang.org/protobuf/proto"
)

func NewRole(agg adaptor.DomainAggregate) domain.Session{
	role := &Role{agg: agg}
	role.Info.IFace = role
	role.Bag.IFace = role
	return role
}

type Role struct {
	agg  adaptor.DomainAggregate
	Info info.Entity
	Bag  bag.Entity
}

type CmdHandler interface {
	Handle(role *Role)(result interface{}, err error)
}
type EvtApplier interface {
	Apply(role *Role)
}



func (the *Role) PublishDomainEvent(eTid uint32, eMsg proto.Message) {	the.agg.AppendEvent(eTid, eMsg) }
func (the *Role) HandleCommand(handler interface{}) (interface{}, error) {	return handler.(CmdHandler).Handle(the)}
func (the *Role) ApplyEvent(applier interface{})                     {	 applier.(EvtApplier).Apply(the)  }

func (the *Role) Load(snapshot []byte) error {

	return nil
}
func (the *Role) Save() (snapshot []byte, err error) {
	panic("implement me")
}


func (the *Role) Uid() uint64 {	return uint64(the.agg.Root()) }

