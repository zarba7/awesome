package role

import (
	"awesome/internal/app/domain_role/role/aggregate"
	"awesome/internal/app/domain_role/role/commands"
	"awesome/internal/app/domain_role/role/events"
	"ddd/util"
)

func Events()  (factory util.SameNameStructFactory){
	table := map[uint32]aggregate.EvtApplier{
		2: &events.DomainEventSetName{},
	}
	for tid, h := range table {
		factory.RegisterStruct(tid, h)
	}
	return
}



func Commands()  (factory util.SameNameStructFactory){
	table := map[uint32]aggregate.CmdHandler{
		2: &commands.Query{},
	}
	for tid, h := range table {
		factory.RegisterStruct(tid, h)
	}
	return
}



