package application

import (
	"reflect"
	"sync"
)

type Dispatcher struct {
	commandTable protocol
	queryTable   protocol
}

func (the *Dispatcher) RegisterCommandHandler(tid int32, h interface{}, from uint8) {	the.commandTable.register(tid, h, from)}
func (the *Dispatcher) RegisterQueryHandler(tid int32, h interface{}, from uint8) {	the.queryTable.register(tid, h, from)}
func (the *Dispatcher) FindQueryHandler(tid int32, from uint8)interface{}{
	return  the.queryTable.find(tid, from)
}
func (the *Dispatcher) FindCommandHandler(tid int32, from uint8)interface{}{
	return  the.commandTable.find(tid, from)
}


type handler struct {
	pool sync.Pool
	from uint8
}

type protocol map[int32]handler

func (the protocol) register(tid int32, h interface{}, from uint8) {
	var e = reflect.TypeOf(h).Elem()
	var pool = sync.Pool{New: func() interface{} { return reflect.New(e).Interface()}}
	the[tid] = handler{pool, from}
}

func (the protocol) find(tid int32, from uint8) interface{} {
	p, ok := the[tid]
	if ok && (from == 0 || p.from == from){
		return p.pool.Get()
	}
	return nil
}


