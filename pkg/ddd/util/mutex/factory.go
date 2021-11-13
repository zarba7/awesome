package mutex

import "sync"

type Factory interface {
	LocalLocker(key string) sync.Locker
}

func NewFactory() Factory {
	the := &factory{list: make(map[string]*mutex)}
	the.pool.New = func() interface{} { return &mutex{} }
	the.lockerPool.New = func() interface{} { return &locker{fac: the} }
	return the
}

type mutex struct {
	sync.Mutex
	n int
	//remoteUnlock func()
}

type factory struct {
	mtx0       sync.Mutex
	list       map[string]*mutex
	pool       sync.Pool
	lockerPool sync.Pool
}

func (the *factory) lock(k string) (mtx *mutex, first bool) {
	the.mtx0.Lock()
	mtx = the.list[k]
	if mtx == nil {
		mtx = the.pool.Get().(*mutex)
		the.list[k] = mtx
		first = true
	}
	mtx.n++
	the.mtx0.Unlock()
	mtx.Lock()
	return
}

func (the *factory) unlock(k string, mtx *mutex) (last bool) {
	mtx.Unlock()
	the.mtx0.Lock()
	mtx.n--
	if mtx.n == 0 {
		last = true
		delete(the.list, k)
	}
	the.mtx0.Unlock()
	return
}

func (the *factory) LocalLocker(key string) sync.Locker {
	var lk = the.lockerPool.Get().(*locker)
	lk.key = key
	return lk
}
