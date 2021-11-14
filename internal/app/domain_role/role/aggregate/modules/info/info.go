package info

import (
	"awesome/internal/app/domain_role/role/aggregate/modules"
)

type Entity struct {
	modules.IFace
	name string
}

func (the *Entity) SetName(name string )  {
	//the.PublishDomainEvent(1, &pbRole.DomainEventSetName{} )
}
