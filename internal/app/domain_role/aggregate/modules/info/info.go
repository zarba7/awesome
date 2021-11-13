package info

import (
	"domain-role/aggregate/modules"
)


type Entity struct {
	modules.IRole
	Name string
}

func (the *Entity) SetName(name string )  {
	//the.PublishDomainEvent(1, &pbRole.DomainEventSetName{} )
}
