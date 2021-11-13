package application

import (
	"ddd/application"
	"src/services/application/command"
	"src/services/application/query"
	"src/services/application/session"
)



type xxxxx int

func commandTable1() map[xxxxx]session.CommandHandler {
	return map[xxxxx]session.CommandHandler{
		1: &command.CreateRole{},
		2: &command.SetRoleName{},
	}
}

func queryTable1() map[xxxxx]session.QueryHandler {
	return map[xxxxx]session.QueryHandler{
		0: &query.RoleInfo{},
	}
}

func protocolTable()(dsp application.Dispatcher) {
	for tid, h := range commandTable1() {
		dsp.RegisterCommandHandler(int32(tid), h, 1)
	}
	for tid, h := range queryTable1() {
		dsp.RegisterQueryHandler(int32(tid), h, 1)
	}
	return
}
