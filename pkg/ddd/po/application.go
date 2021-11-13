package po

type Trx struct {
	AggregateID int64
	Payload     interface{}
}

type Saga struct {
	//事务id
	TransactionID int64
	TrxList       []*Trx
}
