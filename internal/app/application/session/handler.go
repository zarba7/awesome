package session

import (
	"encoding/json"
)

type CommandHandler interface {
	Command(ctx *Context)int32
}
type QueryHandler interface {
	Query(ctx *ContextQ)json.Marshaler
}
