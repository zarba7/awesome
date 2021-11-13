package domain_role

import (
	"awesome/internal/app/domain_role/aggregate"
	"awesome/internal/app/domain_role/events"
	"ddd/util"
)

func CommandFactory()  (factory util.SameNameStructFactory){
	table := map[uint32]aggregate.CmdHandler{
		2: &createRole{},
	}
	for tid, h := range table {
		factory.RegisterStruct(tid, h)
	}
	return
}


func EventsFactory()  (factory util.SameNameStructFactory){
	table := map[uint32]aggregate.EvtApplier{
		2: &events.DomainEventSetName{},
	}
	for tid, h := range table {
		factory.RegisterStruct(tid, h)
	}
	return
}


