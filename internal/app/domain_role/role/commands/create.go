package commands

import (
	"awesome/internal/app/domain_role/role/aggregate"
	"awesome/internal/proto/pbRole"
)

type Query struct {
	pbRole.Query
}

func (args *Query) Handle(role *aggregate.Role) (result interface{}, err error) {
	panic("implement me")
}
