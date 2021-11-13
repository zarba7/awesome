package session

import (
	"ddd/dto"
	"base/json"
)

type Context struct {
	*ContextQ
	From    uint8
	Request *dto.AppPacket
}

func (the *Context) Command() error {
	var pack = the.Result.Packet
	C := the.s.dispatcher.FindCommandHandler(pack.Tid, the.From)
	if C != nil{
		json.UnmarshalPanic(pack.Content, C)
		C.Command(&the.Result)
		the.Result.Packet = nil
	}
	return nil
}



