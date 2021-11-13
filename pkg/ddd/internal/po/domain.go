package po

import "encoding/json"

type DomainAggregate struct {
	ID              int64 //`聚合id`
	Tag             string
	CurrVersion     int64  //`当前版本`
	SnapshotVersion int64  //`快照版本`
	Snapshot        []byte //`快照`
	Events          []*DomainEvent
}

type DomainEvent struct {
	Uid       int64  //事件id
	TrxID     int64  //事务id
	CreatedAt int64  //创建时间
	Tid       uint32 //类型
	Tag       string //标签
	Content   []byte
}

func (e *DomainEvent) UnmarshalContent(args interface{}) error {
	return json.Unmarshal(e.Content, args)
}
func (e *DomainEvent) MarshalContent(args interface{}) (err error) {
	e.Content, err = json.Marshal(args)
	return
}
