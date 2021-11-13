package session

import (
	"ddd/dto"
	"base/json"
)

type ContextQ struct {

}

func (the *ContextQ) Query(pack *dto.AppPacket) *dto.AppResponse  {
	Q := the.s.Dispatcher.FindQueryHandler(pack.Tid, the.from)
	if Q != nil{
		the.Result.Packet = pack
		json.UnmarshalPanic(pack.Content, Q)
		result := Q.Query()
		pack.Content,_ = result.MarshalJSON()
		return &the.Result
	}
	return nil
}
