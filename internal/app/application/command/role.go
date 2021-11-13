package command

import (
	"src/services/application/session"
)

type CreateRole struct {

}

func (C *CreateRole) Command(ctx *session.Context) int32 {
	panic("implement me")
}



type SetRoleName struct {

}

func (C *SetRoleName) Command(ctx *session.Context) int32 {
	panic("implement me")
}

