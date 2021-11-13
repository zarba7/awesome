package saga

import (
	"base/fsm"
)


type OrchestratorFactory struct {
	lSubTransaction []ISubTransaction
	hConditions []ICondition
}
func (orh *OrchestratorFactory) RegisterSubTransaction(trx ISubTransaction) {

}
func (orh *OrchestratorFactory) RegisterCondition(trx fsm.ICondition) {

}
func (orh *OrchestratorFactory) ParseFrom(jsonStr string) *Orchestrator  {

}
