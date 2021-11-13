package events

import (
	"internal/app/domain_role/aggregate"
	"pbRole"
)


type DomainEventSetName struct {
	pbRole.DomainEventSetName
}

func (ev *DomainEventSetName) Apply(role *aggregate.Role)  {

}

