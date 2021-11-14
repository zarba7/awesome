package events

import (
	"awesome/internal/app/domain_role/role/aggregate"
	"awesome/internal/proto/pbRole"
)


type DomainEventSetName struct {
	pbRole.DomainEventSetName
}

func (ev *DomainEventSetName) Apply(role *aggregate.Role)  {

}

