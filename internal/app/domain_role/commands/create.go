package commands

import (
	"internal/app/domain_role/aggregate"
)

type CreateRole struct {
	pbRole.CreateRole
}

func (args *CreateRole) Handle(role *aggregate.Role) (result interface{}, err error) {
	panic("implement me")
}
