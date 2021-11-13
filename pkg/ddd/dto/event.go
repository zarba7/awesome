package dto

type Event struct {
	Uid 		int64 //事件id
	TrxID   int64	//事务id
	CreatedAt   int64 //创建时间
	Tid         uint32 //类型
	Tag         string //标签
	Content     interface{}
	AggID int64
	AggTag string
}

