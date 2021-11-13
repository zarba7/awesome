package query

import (
	"encoding/json"
	"src/services/application/session"
)

type RoleInfo struct {

}

func (req *RoleInfo) Query(ctx *session.ContextQ) json.Marshaler {
	panic("implement me")
}


