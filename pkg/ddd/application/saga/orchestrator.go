package saga

import "base/fsm"


type ISubTransaction interface {
	fsm.IState
}

type ICondition interface {
	fsm.ICondition
}


type Orchestrator struct {
	fsm.Machine
}

