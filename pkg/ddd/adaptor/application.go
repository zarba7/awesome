package adaptor


type UidProducer interface {
	Incr() int64
}

type SagaRepository interface {
	//跟据类型产生Uid
	UidProducer
	Save() error
	Load() error
}