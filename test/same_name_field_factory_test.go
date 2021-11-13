package test

import (
	"ddd/util"
	"encoding/json"
	"awesome/internal/proto/pbRole"
	"reflect"
	"testing"
)

type DomainEventSetName struct {
	*pbRole.DomainEventSetName
}

type DomainEventAddItem struct {
	ii int32
	*pbRole.DomainEventAddItem
}

func (t *DomainEventSetName) Set() {

}

type Set interface {
	Set()
}

func Import(dispatcher *util.SameNameStructFactory) {
	table := map[int32]interface{}{
		2: DomainEventSetName{},
		3: &DomainEventAddItem{},
	}
	for tid, h := range table {
		dispatcher.RegisterStruct(tid, h)
	}
}

var dispatcher = util.SameNameStructFactory{}

func init() {
	Import(&dispatcher)
}

func TestSameNameFieldHandler_Init3(t *testing.T) {
	var x interface{}
	x = DomainEventSetName{DomainEventSetName: &pbRole.DomainEventSetName{}}
	fun := func() interface{} { return x }
	yy := reflect.ValueOf(fun()).Interface()
	yt := reflect.ValueOf(yy).Addr()
	y := yt.Addr().Interface()
	y.(*DomainEventSetName).DomainEventSetName = nil
	t.Log(y)
}

func TestSameNameFieldHandler_Init(t *testing.T) {

	h := dispatcher.FindStructByFieldVal(2, &pbRole.DomainEventSetName{Name: "xxxx"})
	data, _ := json.Marshal(&pbRole.DomainEventSetName{Name: "dddd"})
	h = dispatcher.FindStructByUnmarshal(2, func(args interface{}) {
		json.Unmarshal(data, args)
	})
	h.(Set).Set()
}

func BenchmarkSameNameFieldDispatcher_FindHandlerByFieldVal(b *testing.B) {
	xx := &pbRole.DomainEventAddItem{
		Tid: 3,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dispatcher.FindStructByFieldVal(3, xx)
	}
}
