package fsm

type IState interface {
	String()
	OnEnter()
	Do()
	OnLeave()
}
type ICondition interface {
	String() bool
	OK() bool
}
