package fsm

//finite-state machine
type Machine struct {
	currS IState
}

//fsm de
func (M *Machine)Update(){
	
}

func (M *Machine)End (cond ICondition, stat IState) *Machine {
	return M
}
type state struct {
	*Machine
	stat IState
}
func (S *state)NextState (cond ICondition, stat IState) *state{
	return S
}
func Begin (stat IState) *state{
	M := &Machine{}
	S := &state{
		Machine: M,
		stat:    stat,
	}
	return S
}

