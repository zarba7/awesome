package aggregate

import (
	"ddd/adaptor"
	"ddd/domain"
)

func RoleFactory() func(agg adaptor.DomainAggregate) domain.Session {
	return func(agg adaptor.DomainAggregate) domain.Session {
		role := &Role{agg: agg}
		role.Info.IRole = role
		role.Bag.IRole = role
		return role
	}
}

type CmdHandler interface {
	Handle(role *Role)(result interface{}, err error)
}
type EvtApplier interface {
	Apply(role *Role)
}
